package models

import (
	"time"
)

// BlockDB represents a block header stored in the database.
// It maps to the 'blocks' table.
type BlockDB struct {
	// Core Identifiers
	Hash   []byte `gorm:"primaryKey;type:bytea"`
	Height uint64 `gorm:"not null;unique"`

	// Header Fields
	ParentHash       []byte    `gorm:"not null;type:bytea"`
	StateRoot        []byte    `gorm:"not null;type:bytea"`
	TransactionsRoot []byte    `gorm:"not null;type:bytea"`
	ReceiptRoot      []byte    `gorm:"not null;type:bytea"`
	UncleHash        []byte    `gorm:"not null;type:bytea"`
	Proposer         []byte    `gorm:"not null;type:bytea"`
	Prevrandao       []byte    `gorm:"not null;type:bytea"`
	Timestamp        time.Time `gorm:"column:timestamp;not null"`
	GasUsed          uint64    `gorm:"not null"`
	Bloom            []byte    `gorm:"not null;type:bytea"` // 256 bytes
	ExtraData        []byte    `gorm:"type:bytea"`
	Signature        []byte    `gorm:"type:bytea"`

	// Denormalized count for convenience
	TransactionCount int `gorm:"not null"`

	// Relationships
	Transactions []*TransactionDB `gorm:"foreignKey:BlockHash;references:Hash"`
}

func (BlockDB) TableName() string {
	return "blocks"
}
