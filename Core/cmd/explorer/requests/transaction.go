package requests

import (
	"math/big"

	"FichainCore/cmd/explorer/models"
	"FichainCore/common"
)

// TransactionResponse is the DTO for a transaction, ready for JSON serialization.
type TransactionResponse struct {
	Hash             string           `json:"hash"`
	BlockHash        string           `json:"blockHash"`
	BlockHeight      uint64           `json:"blockHeight"`
	TransactionIndex uint32           `json:"transactionIndex"`
	FromAddress      string           `json:"fromAddress"`
	ToAddress        string           `json:"toAddress,omitempty"` // omitempty for contract creation
	Nonce            uint64           `json:"nonce"`
	Amount           *big.Int         `json:"amount"`
	GasLimit         uint64           `json:"gasLimit"`
	GasPrice         *big.Int         `json:"gasPrice"`
	Data             string           `json:"data"`
	Message          string           `json:"message"`
	Signature        string           `json:"signature"`
	Logs             []*LogResponse   `json:"logs"`
	Receipt          *ReceiptResponse `json:"receipt"`
}

// toTransactionResponse converts a database model to its corresponding response DTO.
func ToTransactionResponse(tx *models.TransactionDB) *TransactionResponse {
	if tx == nil {
		return nil
	}

	// Convert ToAddress, handling the nil case for contract creation
	var toAddrHex string
	toAddrHex = common.BytesToAddress(tx.ToAddress).Hex()

	// Convert FromAddress (should not be nil)
	var fromAddrHex string
	fromAddrHex = common.BytesToAddress(tx.FromAddress).Hex()

	// Convert related Logs
	logResponses := make([]*LogResponse, len(tx.Logs))
	for i, log := range tx.Logs {
		logResponses[i] = ToLogResponse(log)
	}

	// Convert related Receipt
	var receiptResponse *ReceiptResponse
	if tx.Receipt != nil {
		receiptResponse = ToReceiptResponse(tx.Receipt)
	}
	amount := new(big.Int)
	amount.SetBytes(tx.Amount)

	gasPrice := new(big.Int)
	gasPrice.SetBytes(tx.GasPrice)

	return &TransactionResponse{
		Hash:             common.Bytes2Hex(tx.Hash),
		BlockHash:        common.Bytes2Hex(tx.BlockHash),
		BlockHeight:      tx.BlockHeight,
		TransactionIndex: tx.TransactionIndex,
		FromAddress:      fromAddrHex,
		ToAddress:        toAddrHex,
		Nonce:            tx.Nonce,
		Amount:           amount,
		GasLimit:         tx.GasLimit,
		GasPrice:         gasPrice,
		Data:             common.Bytes2Hex(tx.Data),
		Message:          tx.Message,
		Signature:        common.Bytes2Hex(tx.Signature),
		Logs:             logResponses,
		Receipt:          receiptResponse,
	}
}
