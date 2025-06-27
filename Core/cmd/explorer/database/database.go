package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"FichainCore/cmd/explorer/models"
)

// DB is a package-level variable to hold the database connection pool.
var DB *gorm.DB

// Config holds all the configuration for the database connection.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Init initializes the database connection and runs migrations.
// It returns the gorm.DB instance and an error if one occurred.
func Init(config Config) (*gorm.DB, error) {
	// 1. Construct the Data Source Name (DSN) string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	// 2. Configure the GORM logger
	// This provides more detailed logs during development.
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: true,                   // Don't log "record not found" errors
			Colorful:                  true,                   // Enable colorful logging
		},
	)

	// 3. Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Database connection established successfully.")

	// 4. Assign the connection to the package-level variable
	DB = db

	// 5. Run database migrations
	if err := migrate(db); err != nil {
		log.Printf("Failed to run database migrations: %v", err)
		return nil, err
	}

	return DB, nil
}

// migrate runs the GORM AutoMigrate function.
// It will CREATE or UPDATE tables to match the Go structs.
func migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// GORM will automatically create tables, columns, and foreign keys.
	// It will ONLY add missing columns, not remove/change existing ones.
	err := db.AutoMigrate(
		&models.BlockDB{},
		&models.TransactionDB{},
		&models.ReceiptDB{},
		&models.LogDB{},
	)

	if err == nil {
		log.Println("Database migration completed successfully.")
	}

	return err
}
