package models

// ReceiptDB represents a transaction receipt stored in the database.
// It maps to the 'receipts' table.
type ReceiptDB struct {
	TransactionHash   []byte `gorm:"primaryKey;type:bytea"`
	Status            uint32 `gorm:"not null"` // TRUE for success (1), FALSE for failure (0)
	CumulativeGasUsed uint64 `gorm:"not null"`
	GasUsed           uint64 `gorm:"not null"`
	ContractAddress   []byte `gorm:"type:bytea"`          // Pointer for nil when not a contract creation
	LogsBloom         []byte `gorm:"not null;type:bytea"` // 256 bytes
}

func (ReceiptDB) TableName() string {
	return "receipts"
}
