package testing

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

func TestConfigFile(t *testing.T) {
	configPath := "C:/WorkSpace/GoGoChain/Development/GoGoChain/testing/config/config.toml" // Update with the path to your TOML config file
	if err := PrintConfigValues(configPath); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
	}
}

type AppConfig struct {
	Bootnodes []string `toml:"bootnodes"`
	// Add more configuration fields here.
}

func PrintConfigValues(configPath string) error {
	config, err := toml.LoadFile(configPath)
	if err != nil {
		return err
	}

	var appConfig AppConfig
	if err := config.Unmarshal(&appConfig); err != nil {
		return err
	}

	fmt.Println("Bootnodes:", appConfig.Bootnodes)
	// Print other configuration values here if needed

	return nil
}

func TestConfigFileusingViper(t *testing.T) {
	// Set the path to your config.toml file
	viper.SetConfigFile("C:/WorkSpace/GoGoChain/Development/GoGoChain/testing/config/config.toml")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	// Retrieve values from the configuration
	bootnodes := viper.GetStringSlice("Node.P2P.bootnodes")

	// Print the values
	fmt.Println("Bootnodes:")
	for _, node := range bootnodes {
		fmt.Println(node)
	}
}
