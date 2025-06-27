// requests/deposit_wallet_response.go
package requests

import (
	"math/big"
	"time"

	"FichainCore/common"

	"FichainBridge/models" // Assumes this is your project's module path
)

// DepositWalletResponse is the DTO (Data Transfer Object) for a DepositWallet.
// It formats the data for API responses, ensuring sensitive information like the
// private key is never exposed and data formats are client-friendly.
type DepositWalletResponse struct {
	ID             uint                `json:"id"`
	FichainAddress string              `json:"fichainAddress"`
	Address        string              `json:"address"`
	TokenName      string              `json:"tokenName"`
	Balance        *big.Int            `json:"balance"`
	TotalWithdrawn *big.Int            `json:"totalWithdrawn"`
	LockedBalance  *big.Int            `json:"lockedBalance"`
	Status         models.WalletStatus `json:"status"`
	Nonce          uint64              `json:"nonce"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
	// Use a pointer so it can be omitted if null (omitempty)
	LastBalanceSyncAt *time.Time `json:"lastBalanceSyncAt,omitempty"`
}

// ToDepositWalletResponse converts a database model of a DepositWallet into a
// client-facing DepositWalletResponse DTO.
func ToDepositWalletResponse(wallet *models.DepositWallet) *DepositWalletResponse {
	// A robust function should handle nil input gracefully.
	if wallet == nil {
		return nil
	}

	// Helper function to safely convert a byte slice to a *big.Int.
	// An empty or nil slice will correctly result in a big.Int with a value of 0.
	bytesToBigInt := func(b []byte) *big.Int {
		i := new(big.Int)
		if len(b) > 0 {
			i.SetBytes(b)
		}
		return i
	}

	return &DepositWalletResponse{
		ID:                wallet.ID,
		FichainAddress:    common.BytesToAddress(wallet.FichainAddress).Hex(),
		Address:           common.BytesToAddress(wallet.Address).Hex(),
		TokenName:         wallet.TokenName,
		Balance:           bytesToBigInt(wallet.Balance),
		TotalWithdrawn:    bytesToBigInt(wallet.TotalWithdrawn),
		LockedBalance:     bytesToBigInt(wallet.LockedBalance),
		Status:            wallet.Status,
		Nonce:             wallet.Nonce,
		CreatedAt:         wallet.CreatedAt,
		UpdatedAt:         wallet.UpdatedAt,
		LastBalanceSyncAt: wallet.LastBalanceSyncAt,
	}
}
