package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"sync"
	"time"

	"FichainCore/common"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"FichainBridge/minter"
	"FichainBridge/models" // Your models package
)

// ... (ERC20TransferTopic and struct definition remain the same) ...
var ERC20TransferTopic = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

type BlockScannerService struct {
	DB                      *gorm.DB
	externalClient          *ethclient.Client
	lastBlockNumberSavePath string
	depositAddresses        map[common.Address]*models.DepositWallet
	tokenContracts          map[common.Address]string
	minter                  *minter.Minter
	mu                      sync.RWMutex
}

func NewBlockScannerService(
	db *gorm.DB,
	tokenContracts map[string]string, // Map of token NAME to token ADDRESS
	network string,
	minter *minter.Minter,
	lastBlockNumberSavePath string,
) (*BlockScannerService, error) {
	// ... (constructor logic is the same) ...
	cli, err := ethclient.Dial(network)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to connect to Ethereum client: %v", err))
		return nil, err
	}

	addrToName := make(map[common.Address]string)
	for name, addr := range tokenContracts {
		addrToName[common.HexToAddress(addr)] = name
	}

	return &BlockScannerService{
		DB:                      db,
		externalClient:          cli,
		tokenContracts:          addrToName,
		depositAddresses:        make(map[common.Address]*models.DepositWallet),
		lastBlockNumberSavePath: lastBlockNumberSavePath,
		minter:                  minter,
	}, nil
}

// LoadInitialAddresses loads all active deposit wallets from the DB into the
// in-memory map at startup.
func (s *BlockScannerService) LoadInitialAddresses() {
	slog.Info("Loading initial deposit addresses into cache...")
	var wallets []*models.DepositWallet
	if err := s.DB.Where("status = ?", models.StatusActive).Find(&wallets).Error; err != nil {
		slog.Error("Failed to load initial wallets from DB", "error", err)
		return
	}

	// Lock once to perform a bulk update
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, wallet := range wallets {
		addr := common.BytesToAddress(wallet.Address)
		s.depositAddresses[addr] = wallet
	}

	slog.Info("Initial deposit address cache loaded", "count", len(s.depositAddresses))
}

// AddWalletToCache is a thread-safe method to add a single newly created wallet
// to the scanner's in-memory cache.
func (s *BlockScannerService) AddWalletToCache(wallet *models.DepositWallet) {
	if wallet == nil {
		return
	}
	addr := common.BytesToAddress(wallet.Address)

	s.mu.Lock()
	s.depositAddresses[addr] = wallet
	s.mu.Unlock()

	slog.Info("Added new wallet to scanner cache", "address", addr.Hex())
}

func (s *BlockScannerService) SubscribeAndProcess() {
	slog.Info("Starting BlockScannerService...")

	// Initial load of addresses, replacing the periodic refresh.
	s.LoadInitialAddresses()

	lastBlockNumber := s.loadLastBlock()
	for {
		// ... (main loop logic for getting blocks is the same) ...
		currentBlock, err := s.externalClient.BlockByNumber(context.Background(), nil)
		if err != nil {
			slog.Error("Error getting latest block number:", "error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		if lastBlockNumber == 0 {
			lastBlockNumber = currentBlock.NumberU64() - 1
		}

		if currentBlock.NumberU64() <= lastBlockNumber {
			time.Sleep(3 * time.Second)
			continue
		}

		for blockNum := lastBlockNumber + 1; blockNum <= currentBlock.NumberU64(); blockNum++ {
			s.processBlock(blockNum)
			s.saveLastBlock(blockNum)
			lastBlockNumber = blockNum
		}
	}
}

// processBlock is the core logic for scanning a single block.
func (s *BlockScannerService) processBlock(blockNumber uint64) {
	// ... (block fetching logic is the same) ...
	block, err := s.externalClient.BlockByNumber(
		context.Background(),
		big.NewInt(int64(blockNumber)),
	)
	if err != nil {
		slog.Error("Error fetching block", "block", blockNumber, "error", err)
		return
	}

	slog.Info(
		"Processing block",
		"number",
		block.NumberU64(),
		"tx_count",
		len(block.Transactions()),
	)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, tx := range block.Transactions() {
		to := tx.To()
		if to == nil {
			continue
		}

		// Case 1: Native Currency Deposit (e.g., BNB)
		// Check if the destination is one of our deposit addresses.
		if _, ok := s.depositAddresses[common.BytesToAddress(to.Bytes())]; ok {
			// Ensure this wallet is for native currency (TokenContractAddress is empty)
			// and that there's actual value being transferred.
			// if wallet.TokenContractAddress == "" && tx.Value().Cmp(big.NewInt(0)) > 0 {
			// 	slog.Info(
			// 		"Detected native deposit",
			// 		"to",
			// 		to.Hex(),
			// 		"amount",
			// 		tx.Value().String(),
			// 		"tx",
			// 		tx.Hash().Hex(),
			// 	)
			// 	s.updateWalletBalance(wallet, tx.Value())
			// }
			// TODO: may skip native and force to use wrapped token
		}

		// Case 2: ERC20 Token Deposit
		// Check if the transaction is to a token contract we are watching.
		if _, isTokenContract := s.tokenContracts[common.BytesToAddress(to.Bytes())]; isTokenContract {
			s.processERC20Logs(tx)
		}
	}
}

// ... (processERC20Logs, updateWalletBalance, and file IO methods are unchanged) ...
func (s *BlockScannerService) processERC20Logs(tx *types.Transaction) {
	receipt, err := s.externalClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		slog.Error("Error getting transaction receipt", "tx", tx.Hash().Hex(), "error", err)
		return
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return // Ignore failed transactions
	}

	tokenName := s.tokenContracts[common.BytesToAddress(tx.To().Bytes())]
	// may need to convert fixed point, but do it later
	for _, logEntry := range receipt.Logs {
		if logEntry.Address == *tx.To() && len(logEntry.Topics) == 3 &&
			logEntry.Topics[0] == ERC20TransferTopic {
			recipientAddress := common.BytesToAddress(logEntry.Topics[2].Bytes())

			if wallet, ok := s.depositAddresses[recipientAddress]; ok &&
				wallet.TokenName == tokenName {
				amount := new(big.Int).SetBytes(logEntry.Data)
				slog.Info(
					"Detected ERC20 deposit",
					"token",
					tokenName,
					"to",
					recipientAddress.Hex(),
					"amount",
					amount.String(),
					"tx",
					tx.Hash().Hex(),
				)
				s.updateWalletBalance(wallet, amount)
				// let mint it

				fichainAddress := common.BytesToAddress(wallet.FichainAddress)

				err := s.minter.MintToken(
					wallet.TokenName,
					fichainAddress,
					common.BytesToHash(logEntry.TxHash.Bytes()),
					amount,
				)
				if err != nil {
					logger.Error(
						"Failed to complete mint for user %s: %v",
						fichainAddress.Hex(),
						err,
					)
				}
			}
		}
	}
}

func (s *BlockScannerService) updateWalletBalance(wallet *models.DepositWallet, amount *big.Int) {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var currentWallet models.DepositWallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("fichain_address = ? AND token_name = ?", wallet.FichainAddress, wallet.TokenName).
			First(&currentWallet).Error; err != nil {
			return err
		}

		currentBalance := new(big.Int)
		if len(currentWallet.Balance) > 0 {
			currentBalance.SetBytes(currentWallet.Balance)
		}
		newBalance := new(big.Int).Add(currentBalance, amount)

		currentWallet.Balance = newBalance.Bytes()
		now := time.Now()
		currentWallet.LastBalanceSyncAt = &now

		return tx.Save(&currentWallet).Error
	})
	if err != nil {
		logger.Error(
			"Failed to update wallet balance",
			"fichain_address",
			common.BytesToAddress(wallet.FichainAddress).Hex(),
			"token",
			wallet.TokenName,
			"error",
			err,
		)
	}
}

func (c *BlockScannerService) saveLastBlock(blockNumber uint64) {
	data, _ := json.Marshal(map[string]uint64{"last_block": blockNumber})
	_ = os.WriteFile(c.lastBlockNumberSavePath, data, 0644)
}

func (c *BlockScannerService) loadLastBlock() uint64 {
	data, err := os.ReadFile(c.lastBlockNumberSavePath)
	if err != nil {
		return 0
	}
	var savedData map[string]uint64
	if json.Unmarshal(data, &savedData) != nil {
		return 0
	}
	return savedData["last_block"]
}
