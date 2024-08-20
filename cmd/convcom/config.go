package main

import (
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

const configfile = "convcom.json"

type Config struct {
	Types   []string `json:"types"`
	Scopes  []string `json:"scopes"`
}

// getDefaultConfig returns the default configuration with predefined types and scopes
func getDefaultConfig() Config {
	return Config{
		Types:  []string{"build", "ci", "chore", "docs", "feat", "fix", "perf", "refactor", "revert", "style", "test"},
		Scopes: []string{},
	}
}

// loadConfig loads the local convcom.json containing the available types and scopes
func loadConfig() (*Config, error) {
    data, err := ioutil.ReadFile(configfile)
    if err != nil {
        return nil, err
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}

// createConfigFile creates a config file with the specified name if it does not already exist.
func createConfigFile() error {

	// Check if the file already exists
	if _, err := os.Stat(configfile); !os.IsNotExist(err) {
		return fmt.Errorf("config file %s already exists", configfile)
	}

	// Open the file for writing
	file, err := os.Create(configfile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	// Define the configuration data
	config := getDefaultConfig()

	// Encode the config to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print with indent
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	fmt.Printf("Config file %s created successfully.\n", configfile)
	return nil
}