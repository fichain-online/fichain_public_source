package requests

import (
	"FichainCore/cmd/explorer/models"
	"FichainCore/common"
)

// ReceiptResponse is the DTO for a transaction receipt.
type ReceiptResponse struct {
	TransactionHash   string `json:"transactionHash"`
	Status            uint32 `json:"status"` // 1 for success, 0 for failure
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed"`
	GasUsed           uint64 `json:"gasUsed"`
	ContractAddress   string `json:"contractAddress,omitempty"`
	LogsBloom         string `json:"logsBloom"`
}

// ToReceiptResponse converts a ReceiptDB model to its response DTO.
func ToReceiptResponse(receipt *models.ReceiptDB) *ReceiptResponse {
	if receipt == nil {
		return nil
	}

	return &ReceiptResponse{
		TransactionHash:   common.Bytes2Hex(receipt.TransactionHash),
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		GasUsed:           receipt.GasUsed,
		ContractAddress:   common.Bytes2Hex(receipt.ContractAddress), // Helper handles nil
		LogsBloom:         common.Bytes2Hex(receipt.LogsBloom),
	}
}
