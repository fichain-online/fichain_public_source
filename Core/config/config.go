package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"

	"FichainCore/common"
)

type Config struct {
	NodeID     string
	Version    uint32
	PrivateKey string

	TCPServerAddress string
	BootAddress      string
	WsServerAddress  string
	EkycApiUrl       string

	// storage
	StatesDBPath               string
	AuthorityValidatorDBPath   string
	AuthorityObserverDBPath    string
	AuthorityFiatReserveDBPath string

	// explorers
	ExplorerAddresses []string
}

var (
	instance *Config
	once     sync.Once
)

// InitConfig initializes the config only once
func InitConfig(configPath string) {
	once.Do(func() {
		// Set default values
		viper.SetDefault("node_id", "default-node")
		viper.SetDefault("version", 1)
		viper.SetDefault("EkycApiUrl", "http://127.0.0.1:8080/api/user/ekyc/")

		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Printf("Warning: could not read config file, using defaults. Error: %v", err)
		}

		instance = &Config{}
		if err = viper.Unmarshal(instance); err != nil {
			err = fmt.Errorf("failed to unmarshal config: %w", err)
		}
	})
}

func SetConfig(c *Config) {
	instance = c
}

// GetConfig returns the singleton config instance
func GetConfig() *Config {
	if instance == nil {
		log.Fatal("Config is not initialized. Call InitConfig() first.")
	}
	return instance
}

func (c *Config) GetExplorerAddresses() []common.Address {
	adds := []common.Address{}
	for _, v := range c.ExplorerAddresses {
		adds = append(adds, common.HexToAddress(v))
	}
	return adds
}
