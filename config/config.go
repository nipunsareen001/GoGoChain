package config

import (
	"GoGoChain/enum"
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// IMPROVEMENTS
// More fields can be added here as configuration grows.
// AppConfig holds the configuration for the application.
type AppConfig struct {
	Node struct {
		P2P struct {
			Bootnodes []string `mapstructure:"bootnodes"`
		} `mapstructure:"P2P"`
	} `mapstructure:"Node"`
	// Add more configuration fields here.
}

// appConfig holds the application's configuration.
var appConfig AppConfig

// configMutex ensures that the configuration is only loaded once.
var configMutex sync.Once

// LoadConfig reads the configuration from the given file path.
// It should only be called once at the start of your application.
func LoadConfig() {
	configMutex.Do(func() {
		configPath := string(enum.ConfigFilePath)

		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error reading config file: %w", err))
		}

		if err := viper.Unmarshal(&appConfig); err != nil {
			fmt.Printf("Error unmarshaling config: %v\n", err)
		}
	})
}

// GetConfig returns the loaded configuration.
// This function can be called at any time to retrieve the configuration.
func GetConfig() AppConfig {
	return appConfig
}
