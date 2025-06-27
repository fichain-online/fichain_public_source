// requests/deposit_log_response.go
package requests

import (
	"math/big"
	"time"

	"FichainCore/common" // Assuming FichainCore/common is used for hex conversion

	"FichainBridge/models"
)

// DepositLogResponse is the DTO for a DepositLog, formatted for API responses.
// It converts byte slices into user-friendly strings (hex) and big.Ints.
type DepositLogResponse struct {
	ID                uint             `json:"id"`
	SourceChainTxHash string           `json:"sourceChainTxHash,omitempty"`
	FichainAddress    string           `json:"fichainAddress"`
	TokenName         string           `json:"tokenName"`
	Amount            *big.Int         `json:"amount"`
	DestChainTxHash   string           `json:"destChainTxHash,omitempty"`
	Status            models.LogStatus `json:"status"`
	ErrorMessage      string           `json:"errorMessage,omitempty"`
	RetryCount        uint             `json:"retryCount"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
}

// ToDepositLogResponse converts a database model of a DepositLog into a
// client-facing DepositLogResponse DTO.
func ToDepositLogResponse(log *models.DepositLog) *DepositLogResponse {
	// A robust function should handle nil input gracefully.
	if log == nil {
		return nil
	}

	// Helper function to safely convert a byte slice to a *big.Int.
	bytesToBigInt := func(b []byte) *big.Int {
		i := new(big.Int)
		if len(b) > 0 {
			i.SetBytes(b)
		}
		return i
	}

	// Helper function to safely convert a byte slice to a 0x-prefixed hex string for hashes.
	bytesToHexHash := func(b []byte) string {
		if len(b) == 0 {
			return ""
		}
		return common.BytesToHash(b).Hex()
	}

	return &DepositLogResponse{
		ID:                log.ID,
		SourceChainTxHash: bytesToHexHash(log.SourceChainTxHash),
		FichainAddress:    common.BytesToAddress(log.FichainAddress).Hex(),
		TokenName:         log.TokenName,
		Amount:            bytesToBigInt(log.Amount),
		DestChainTxHash:   bytesToHexHash(log.DestChainTxHash),
		Status:            log.Status,
		ErrorMessage:      log.ErrorMessage,
		RetryCount:        log.RetryCount,
		CreatedAt:         log.CreatedAt,
		UpdatedAt:         log.UpdatedAt,
	}
}

// ToDepositLogListResponse converts a slice of DepositLog models into a slice of DTOs.
func ToDepositLogListResponse(logs []*models.DepositLog) []*DepositLogResponse {
	if logs == nil {
		return nil
	}
	responseList := make([]*DepositLogResponse, len(logs))
	for i, log := range logs {
		responseList[i] = ToDepositLogResponse(log)
	}
	return responseList
}
