package requests

import (
	"FichainCore/cmd/explorer/models"
	"FichainCore/common"
)

// LogResponse is the DTO for a log entry.
type LogResponse struct {
	BlockHash       string   `json:"blockHash"`
	LogIndex        uint32   `json:"logIndex"`
	TransactionHash string   `json:"transactionHash"`
	EmitterAddress  string   `json:"emitterAddress"`
	Data            string   `json:"data"`
	Removed         bool     `json:"removed"`
	Topics          []string `json:"topics"`
}

// ToLogResponse converts a LogDB model to its response DTO.
func ToLogResponse(log *models.LogDB) *LogResponse {
	if log == nil {
		return nil
	}

	// Aggregate all non-nil topics into a single slice.
	topics := make([]string, 0)
	if log.Topic0 != nil {
		topics = append(topics, common.Bytes2Hex(log.Topic0))
	}
	if log.Topic1 != nil {
		topics = append(topics, common.Bytes2Hex(log.Topic1))
	}
	if log.Topic2 != nil {
		topics = append(topics, common.Bytes2Hex(log.Topic2))
	}
	if log.Topic3 != nil {
		topics = append(topics, common.Bytes2Hex(log.Topic3))
	}

	return &LogResponse{
		BlockHash:       common.Bytes2Hex(log.BlockHash),
		LogIndex:        log.LogIndex,
		TransactionHash: common.Bytes2Hex(log.TransactionHash),
		EmitterAddress:  common.Bytes2Hex(log.EmitterAddress),
		Data:            common.Bytes2Hex(log.Data),
		Removed:         log.Removed,
		Topics:          topics,
	}
}
