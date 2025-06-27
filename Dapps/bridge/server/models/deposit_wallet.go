// models/deposit_wallet.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// WalletStatus defines the state of the deposit wallet.
type WalletStatus string

const (
	// StatusPending means the wallet is created but not yet active for use.
	StatusPending WalletStatus = "pending"
	// StatusActive means the wallet is ready to receive deposits.
	StatusActive WalletStatus = "active"
	// StatusLocked means the wallet is temporarily disabled for deposits/withdrawals.
	StatusLocked WalletStatus = "locked"
	// StatusArchived means the wallet is no longer in use.
	StatusArchived WalletStatus = "archived"
)

// DepositWallet represents a unique deposit wallet on a blockchain (like BSC)
// assigned to a user of your bridge system.
type DepositWallet struct {
	gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt

	// --- Foreign Key ---
	// Links this wallet to a user in your main application.
	FichainAddress []byte `gorm:"primaryKey;type:bytea"`

	// --- Blockchain Credentials ---
	// The public address (e.g., "0x..."). This is what users will send funds to.
	Address []byte `gorm:"type:bytea;uniqueIndex;not null"`

	TokenName string

	// !!! SECURITY WARNING !!!
	// Stores the private key after it has been encrypted by a secure,
	// external key management service (KMS, Vault, etc.).
	EncryptedPrivateKey []byte `gorm:"type:bytea;not null"`

	// --- Balance Management ---
	// Balances are stored as byte slices. This is the raw big-endian binary
	// representation of the number (from `*big.Int.Bytes()`).
	// It's more storage-efficient than a string for uint256 values.
	// A nil or empty slice represents a balance of 0.
	// GORM will map this to a BLOB/BYTEA type in the database.

	// Current on-chain balance, updated by a background worker.
	Balance []byte `gorm:"type:bytea"`

	// Total amount that has been successfully withdrawn/bridged from this wallet.
	TotalWithdrawn []byte `gorm:"type:bytea"`

	// Balance that is locked during a pending withdrawal process.
	LockedBalance []byte `gorm:"type:bytea"`

	// --- Metadata and State ---
	// The current status of the wallet.
	Status WalletStatus `gorm:"type:varchar(20);default:'pending';index"`

	// The transaction nonce for this wallet, critical for sending transactions.
	Nonce uint64 `gorm:"default:0"`

	// Timestamp of the last time a background service synced the balance.
	LastBalanceSyncAt *time.Time
}

// TableName sets the custom table name for the model.
func (DepositWallet) TableName() string {
	return "deposit_wallets"
}
