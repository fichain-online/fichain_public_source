package models

// TransactionDB represents a transaction stored in the database.
// It maps to the 'transactions' table.
type TransactionDB struct {
	// Core Identifiers
	Hash             []byte `gorm:"primaryKey;type:bytea"`
	BlockHash        []byte `gorm:"not null;type:bytea"`
	BlockHeight      uint64 `gorm:"not null"`
	TransactionIndex uint32 `gorm:"not null"`

	// Core Transaction Data
	FromAddress []byte `gorm:"not null;type:bytea;index"` // Indexed for faster lookups
	ToAddress   []byte `gorm:"type:bytea;index"`          // Indexed for faster lookups
	Nonce       uint64 `gorm:"not null"`
	Amount      []byte `gorm:"type:bytea;not null"`
	GasLimit    uint64 `gorm:"column:gas_limit;not null"`
	GasPrice    []byte `gorm:"type:bytea;not null"`
	Data        []byte `gorm:"column:data;type:bytea"`
	Message     string

	Signature []byte `gorm:"type:bytea;not null"`

	// Relationships
	Logs    []*LogDB   `gorm:"foreignKey:TransactionHash;references:Hash"`
	Receipt *ReceiptDB `gorm:"foreignKey:TransactionHash;references:Hash"`
}

func (TransactionDB) TableName() string {
	return "transactions"
}
