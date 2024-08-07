package main

import (
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Types   []string `json:"types"`
	Scopes  []string `json:"scopes"`
}

func loadConfig() (*Config, error) {
    data, err := ioutil.ReadFile("convcom.json")
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
	fileName := "convcom.json"
	// Check if the file already exists
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return fmt.Errorf("config file %s already exists", fileName)
	}

	// Define the configuration data
	config := Config{
		Types:  []string{"build", "ci", "chore", "docs", "feat", "fix", "perf", "refactor", "revert", "style", "test"},
		Scopes: []string{},
	}

	// Open the file for writing
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	// Encode the config to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print with indent
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	fmt.Printf("Config file %s created successfully.\n", fileName)
	return nil
}