package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func Setup() error {
	// Create the base application directory
	err := os.MkdirAll(appDir, 0755) // 0755 permissions
	if err != nil {
		return fmt.Errorf("failed to create app directory: %w", err)
	}
	// Create nightly and stable directories
	err = os.MkdirAll(targetNightly, 0755)
	if err != nil {
		return fmt.Errorf("failed to create nightly directory: %w", err)
	}
	err = os.MkdirAll(targetDirStable, 0755)
	if err != nil {
		return fmt.Errorf("failed to create stable directory: %w", err)
	}
	// Create versions_info.json if it doesn't exist
	_, err = os.Stat(versionFilePath)
	if os.IsNotExist(err) {
		// If it doesn't exist, create it and initialize with an empty JSON array
		err = os.WriteFile(versionFilePath, []byte("[]"), 0644) // 0644 permissions
		if err != nil {
			return fmt.Errorf("failed to create versions_info.json: %w", err)
		}
	} else if err != nil {
		// Handle other errors
		return fmt.Errorf("error checking versions_info.json: %w", err)
	}

	err = createConfigFile()
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	return nil
}

// CreateConfigFile creates a config file with default values if it only doesn't exist.
func createConfigFile() error {
	if _, err := os.Stat(configPath); err == nil {
		return nil
	}
	defaultConfig := Config{RollbackLimit: 7}
	configJson, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}
	if err := os.WriteFile(configPath, configJson, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
