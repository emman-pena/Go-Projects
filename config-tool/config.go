package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	AppName string `json:"app_name"`
	Port    int    `json:"port"`
	Debug   bool   `json:"debug"`
}

// LoadConfig loads and merges default and environment-specific configs.
func LoadConfig(env string) (*Config, error) {
	basePath := "./config"
	defaultConfigPath := filepath.Join(basePath, "default.json")
	envConfigPath := filepath.Join(basePath, fmt.Sprintf("%s.json", env))

	config := &Config{}

	// Load default config
	if err := loadFile(defaultConfigPath, config); err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	// Load environment-specific config
	if err := loadFile(envConfigPath, config); err != nil {
		return nil, fmt.Errorf("failed to load %s config: %w", env, err)
	}

	return config, nil
}

/**
Parameters:

filePath string: The path to the JSON file you want to load.
config *Config: A pointer to a Config struct where the data will be loaded.
Return Value:

It returns an error. If something goes wrong, it provides details about the issue.

*/
// Helper to load a file into the config struct
func loadFile(filePath string, config *Config) error {

	// Opens the file at filePath for reading using os.Open.
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	/** Creates a new JSON decoder for the file. The decoder will read the JSON
	data from the file stream.

	Uses the Decode method to parse the JSON data and populate the config struct.
	If decoding fails (e.g., due to invalid JSON structure), it returns the decoding error.
	*/
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return err
	}

	return nil
}
