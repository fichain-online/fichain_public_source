// package minter
package minter

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"FichainCore/client"
	"FichainCore/common"
	"FichainCore/params"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"FichainBridge/models" // Your models package
)

const (
	FichainGasLimit = uint64(200000)
)

var (
	ErrUnsupportedToken    = errors.New("unsupported token")
	ErrInsufficientBalance = errors.New("insufficient available balance to mint")
)

type Minter struct {
	DB       *gorm.DB
	txClient *client.Client
	tokenMap map[string]common.Address
}

func NewMinter(db *gorm.DB, txClient *client.Client, tokenMap map[string]common.Address) *Minter {
	return &Minter{DB: db, txClient: txClient, tokenMap: tokenMap}
}

// MintToken finds the available balance for a user's wallet and attempts to mint it.
// It creates a full audit trail using the DepositLog model.
func (m *Minter) MintToken(
	tokenName string,
	fichainUserAddress common.Address,
	sourceTxHash common.Hash,
	amount *big.Int,
) error {
	var logEntry *models.DepositLog

	// Phase 1: Calculate available funds and lock them in an atomic DB transaction.
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		var wallet models.DepositWallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("fichain_address = ? AND token_name = ?", fichainUserAddress.Bytes(), tokenName).
			First(&wallet).Error; err != nil {
			return err // Wallet not found or DB error.
		}

		lockedBalance := new(big.Int).SetBytes(wallet.LockedBalance)

		// Create the audit log for this specific operation *inside* the transaction.
		logEntry = &models.DepositLog{
			FichainAddress:    fichainUserAddress.Bytes(),
			SourceChainTxHash: sourceTxHash.Bytes(),
			TokenName:         tokenName,
			Amount:            amount.Bytes(),
			Status:            models.LogStatusProcessing,
		}
		if err := tx.Create(logEntry).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		// Lock the available funds by adding them to the LockedBalance.
		wallet.LockedBalance = new(big.Int).Add(lockedBalance, amount).Bytes()

		return tx.Save(&wallet).Error
	})
	// Handle outcomes of the locking transaction.
	if err != nil {
		if errors.Is(err, ErrInsufficientBalance) {
			// This is an expected case, not an error to be propagated.
			return nil
		}
		logger.Error(
			"Failed to lock funds for minting",
			"user",
			fichainUserAddress.Hex(),
			"token",
			tokenName,
			"error",
			err,
		)
		return err
	}

	logger.Info(
		"Funds locked for minting",
		"logID",
		logEntry.ID,
		"user",
		fichainUserAddress.Hex(),
		"amount",
		amount.String(),
	)

	// Phase 2: Attempt the on-chain transaction.
	// Pass the logEntry so it can be updated based on the outcome.
	return m.attemptOnChainMint(tokenName, fichainUserAddress, amount, logEntry)
}

// RetryMint attempts to mint a wallet's *entire locked balance*.
// This is for a worker to call on wallets with stuck funds.
func (m *Minter) RetryMint(wallet *models.DepositWallet) error {
	if wallet == nil {
		return errors.New("cannot retry mint on a nil wallet")
	}

	lockedBalance := new(big.Int).SetBytes(wallet.LockedBalance)
	if lockedBalance.Cmp(big.NewInt(0)) <= 0 {
		return nil // Nothing to retry.
	}

	// Create a new audit log for this specific retry attempt.
	logEntry := &models.DepositLog{
		FichainAddress: wallet.FichainAddress,
		TokenName:      wallet.TokenName,
		Amount:         wallet.LockedBalance,
		Status:         models.LogStatusProcessing,
		RetryCount:     1, // Indicate this is a retry.
	}
	if err := m.DB.Create(logEntry).Error; err != nil {
		logger.Error(
			"Failed to create audit log for retry",
			"user",
			common.BytesToAddress(wallet.FichainAddress).Hex(),
			"error",
			err,
		)
		return err
	}
	logger.Info(
		"Retrying to mint locked balance",
		"logID",
		logEntry.ID,
		"user",
		common.BytesToAddress(wallet.FichainAddress).Hex(),
		"amount",
		lockedBalance.String(),
	)

	// The balance is already locked, so we go straight to the on-chain attempt.
	return m.attemptOnChainMint(
		wallet.TokenName,
		common.BytesToAddress(wallet.FichainAddress),
		lockedBalance,
		logEntry,
	)
}

// attemptOnChainMint contains the logic for sending the Fichain transaction and handling the outcome.
func (m *Minter) attemptOnChainMint(
	tokenName string,
	fichainUserAddress common.Address,
	amount *big.Int,
	logEntry *models.DepositLog,
) error {
	tokenContractAddress, ok := m.tokenMap[tokenName]
	if !ok {
		m.updateLogState(logEntry, models.LogStatusFailed, nil, ErrUnsupportedToken.Error())
		return ErrUnsupportedToken
	}

	transferData, err := createTransferDataPayload(fichainUserAddress, amount)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create transfer data payload: %v", err)
		m.updateLogState(logEntry, models.LogStatusFailed, nil, errMsg)
		return errors.New(errMsg)
	}

	txMessage := fmt.Sprintf("Bridge Mint for log %d", logEntry.ID)
	receipt, err := m.txClient.SendTransaction(
		tokenContractAddress,
		big.NewInt(0),
		transferData,
		FichainGasLimit,
		params.TempGasPrice,
		txMessage,
	)

	// Handle transaction failure. The funds remain locked. We only update the log.
	if err != nil || (receipt != nil && receipt.Status != 1) {
		errMsg := "unknown error"
		if err != nil {
			errMsg = err.Error()
		} else if receipt != nil {
			errMsg = fmt.Sprintf("transaction reverted on-chain (status %d)", receipt.Status)
		}

		logger.Error(
			"Fichain transaction failed. Funds remain locked for retry.",
			"logID",
			logEntry.ID,
			"user",
			fichainUserAddress.Hex(),
			"error",
			errMsg,
		)
		m.updateLogState(logEntry, models.LogStatusFailed, nil, errMsg)
		return errors.New(errMsg)
	}

	// Transaction was successful. Finalize the state.
	logger.Info(
		"Fichain transaction successful",
		"logID",
		logEntry.ID,
		"hash",
		receipt.TxHash.Hex(),
	)
	// set tx hash
	logEntry.DestChainTxHash = receipt.TxHash.Bytes()
	m.finalizeMint(tokenName, fichainUserAddress, amount, logEntry)

	return nil
}

// finalizeMint confirms the mint by moving the amount from LockedBalance to TotalWithdrawn.
func (m *Minter) finalizeMint(
	tokenName string,
	fichainUserAddress common.Address,
	amount *big.Int,
	logEntry *models.DepositLog,
) {
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		var wallet models.DepositWallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("fichain_address = ? AND token_name = ?", fichainUserAddress.Bytes(), tokenName).
			First(&wallet).Error; err != nil {
			return err
		}

		lockedBalance := new(big.Int).SetBytes(wallet.LockedBalance)
		totalWithdrawn := new(big.Int).SetBytes(wallet.TotalWithdrawn)

		// This check is a safeguard.
		if lockedBalance.Cmp(amount) < 0 {
			return fmt.Errorf(
				"cannot finalize mint, locked balance %s is less than mint amount %s",
				lockedBalance,
				amount,
			)
		}

		// Atomically move funds from locked to withdrawn.
		wallet.LockedBalance = new(big.Int).Sub(lockedBalance, amount).Bytes()
		wallet.TotalWithdrawn = new(big.Int).Add(totalWithdrawn, amount).Bytes()

		return tx.Save(&wallet).Error
	})
	if err != nil {
		// This is a critical error. The user has been credited on-chain, but our DB state is inconsistent.
		logger.Error(
			"CRITICAL: FAILED TO FINALIZE MINT IN DB AFTER SUCCESSFUL TRANSACTION",
			"logID",
			logEntry.ID,
			"user",
			fichainUserAddress.Hex(),
			"error",
			err,
		)
		// We still try to update the log, but mark the error.
		m.updateLogState(
			logEntry,
			models.LogStatusFailed,
			logEntry.DestChainTxHash,
			fmt.Sprintf("DB finalization failed: %v", err),
		)
		return
	}

	// DB update was successful, now update the log to success.
	m.updateLogState(logEntry, models.LogStatusSuccess, logEntry.DestChainTxHash, "")
}

// createTransferDataPayload encodes the arguments for an ERC20 'transfer' call.
func createTransferDataPayload(recipient common.Address, amount *big.Int) ([]byte, error) {
	parsedABI, err := abi.JSON(
		strings.NewReader(
			`[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	// Pack the arguments with the method name.
	return parsedABI.Pack("transfer", recipient, amount)
}

// updateLogState is a helper to centralize log updates.
func (m *Minter) updateLogState(
	log *models.DepositLog,
	status models.LogStatus,
	destTxHash []byte,
	errMsg string,
) {
	updates := map[string]interface{}{
		"status":             status,
		"dest_chain_tx_hash": destTxHash,
		"error_message":      errMsg,
	}
	if err := m.DB.Model(log).Updates(updates).Error; err != nil {
		logger.Error("CRITICAL: Failed to update audit log state", "logID", log.ID, "error", err)
	}
}
