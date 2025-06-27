package config

import (
	"log"
	"sync"

	core_config "FichainCore/config"

	"github.com/spf13/viper"

	"FichainBridge/database"
)

// Config holds all configuration for the explorer application.
// It nests the core and database configurations.
type Config struct {
	Core     core_config.Config `mapstructure:"core"`
	Database database.Config    `mapstructure:"database"`

	APIAddress              string
	TokenMap                map[string]string `mapstructure:"token_map"`
	FichainTokenMap         map[string]string `mapstructure:"fichain_token_map"`
	NetworkConnectionString string
}

var (
	instance *Config
	once     sync.Once
)

// InitConfig initializes the explorer and its components from a single config file.
func InitConfig(configPath string) {
	once.Do(func() {
		// Set defaults for all nested configurations using their full path.
		viper.SetDefault("APIAddress", ":8080")

		// Core defaults
		viper.SetDefault("core.node_id", "simple-bridge-server")
		viper.SetDefault("core.version", 1)
		viper.SetDefault("core.tcpserveraddress", ":3000")
		viper.SetDefault("core.wsserveraddress", ":8080")

		// Database defaults
		viper.SetDefault("database.path", "./bridge_db")
		viper.SetDefault("database.user", "user")

		viper.SetConfigFile(configPath)
		viper.AutomaticEnv() // Recommended: allows overriding config with environment variables

		if err := viper.ReadInConfig(); err != nil {
			log.Printf(
				"Warning: could not read config file '%s', using defaults. Error: %v",
				configPath,
				err,
			)
		}

		instance = &Config{}
		if err := viper.Unmarshal(instance); err != nil {
			log.Fatalf("Failed to unmarshal config: %v", err)
		}
	})
}

// GetConfig returns the singleton config instance.
func GetConfig() *Config {
	if instance == nil {
		log.Fatal("Config is not initialized. Call InitConfig() first.")
	}
	return instance
}
