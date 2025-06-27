// models/deposit_log.go
package models

import (
	"gorm.io/gorm"
)

// LogStatus defines the state of a minting operation.
type LogStatus string

const (
	// LogStatusPending means the deposit has been detected and is waiting to be minted.
	LogStatusPending LogStatus = "pending"
	// LogStatusProcessing means a minter has picked up the log and is attempting the transaction.
	LogStatusProcessing LogStatus = "processing"
	// LogStatusSuccess means the minting transaction on Fichain was successful.
	LogStatusSuccess LogStatus = "success"
	// LogStatusFailed means the minting transaction failed and may need manual review or a retry.
	LogStatusFailed LogStatus = "failed"
)

// DepositLog tracks a single deposit event from the source chain (e.g., BSC)
// and its corresponding minting status on the destination chain (Fichain).
type DepositLog struct {
	gorm.Model // Includes ID, CreatedAt, UpdatedAt

	// --- Source Chain Info ---

	// The transaction hash of the deposit on the source chain (e.g., BSC).
	// This MUST be unique to prevent double-processing of the same deposit.
	SourceChainTxHash []byte `gorm:"type:bytea;uniqueIndex"`

	// --- Destination Chain Info ---

	// The Fichain address of the user who should receive the minted tokens.
	FichainAddress []byte `gorm:"type:bytea;index"`

	// The name of the token to be minted.
	TokenName string `gorm:"index"`

	// The amount of tokens to be minted (in the smallest unit, e.g., WEI).
	Amount []byte `gorm:"type:bytea"`

	// The transaction hash of the minting operation on Fichain.
	DestChainTxHash []byte `gorm:"type:bytea;index"`

	// --- Process State ---

	// The current status of this minting operation.
	Status LogStatus `gorm:"type:varchar(20);default:'pending';index"`

	// Stores error messages if the minting process fails.
	ErrorMessage string

	// The number of times a retry has been attempted for this log.
	RetryCount uint `gorm:"default:0"`
}

// TableName sets the custom table name for the model.
func (DepositLog) TableName() string {
	return "deposit_logs"
}
