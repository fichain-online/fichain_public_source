package models

// LogDB represents an event log stored in the database.
// It maps to the 'logs' table.
type LogDB struct {
	// Composite Primary Key
	BlockHash []byte `gorm:"primaryKey;type:bytea"`
	LogIndex  uint32 `gorm:"primaryKey"`

	// Foreign key to transaction
	TransactionHash []byte `gorm:"not null;type:bytea"`

	// Log Data
	EmitterAddress []byte `gorm:"not null;type:bytea"`
	Data           []byte `gorm:"column:data;type:bytea"`
	Removed        bool   `gorm:"not null"`

	// Indexed Topics (pointers to handle nil/missing topics)
	Topic0 []byte `gorm:"type:bytea"`
	Topic1 []byte `gorm:"type:bytea"`
	Topic2 []byte `gorm:"type:bytea"`
	Topic3 []byte `gorm:"type:bytea"`
}

func (LogDB) TableName() string {
	return "logs"
}
