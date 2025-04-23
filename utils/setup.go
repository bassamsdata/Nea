package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"
)

func Setup() error {
	// Existing directories
	dirs := []string{
		appDir,
		targetNightly,
		targetDirStable,
		filepath.Join(appDir, "bin"),
	}

	// Create all directories
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	var err error
	// Create versions_info.json if it doesn't exist
	_, err = os.Stat(versionFilePath)
	if os.IsNotExist(err) {
		// If it doesn't exist, create it and initialize with an empty JSON array
		err = os.WriteFile(versionFilePath, []byte("[]"), 0o644) // 0644 permissions
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

	// Check and notify about PATH setup
	binDir := filepath.Join(appDir, "bin")
	if !isInPath(binDir) {
		color.Yellow("\nImportant: neomanager bin directory is not in your PATH")
		fmt.Printf("\nAdd this line to your shell configuration file (.zshrc, .bashrc, etc.):\n")
		fmt.Printf("export PATH=\"%s:$PATH\"\n", binDir)
		fmt.Printf("\nThen restart your terminal or run:\n")
		fmt.Printf("source ~/.zshrc  # or your shell's config file\n\n")
	}

	return nil
}

// CreateConfigFile creates a config file with default values if it only doesn't exist.
func createConfigFile() error {
	if _, err := os.Stat(SymlinkPath); err == nil {
		return nil
	}
	defaultConfig := Config{RollbackLimit: 7}
	configJson, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}
	if err := os.WriteFile(SymlinkPath, configJson, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func isInPath(dir string) bool {
	path := os.Getenv("PATH")
	paths := strings.Split(path, ":")
	return slices.Contains(paths, dir)
}
