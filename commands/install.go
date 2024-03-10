package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/progress"
)

var (
	homeDir         = os.Getenv("HOME")
	appDir          = filepath.Join(homeDir, ".local", "share", "nv_manager")
	neovimURL       = "https://github.com/neovim/neovim/releases/download/nightly/nvim-macos.tar.gz"
	targetNightly   = filepath.Join(homeDir, ".local", "share", "nv_manager", "nightly")
	stableBaseURL   = "https://github.com/neovim/neovim/releases/download/v"
	targetDirStable = filepath.Join(homeDir, ".local", "share", "nv_manager", "stable")
	tagsURL         = "https://api.github.com/repos/neovim/neovim/tags"
	versionFilePath = filepath.Join(targetNightly, "versions_info.json")
	nvm_night_url   = "https://github.com/neovim/neovim/releases/download/nightly/nvim-macos.tar.gz"
)

func InstallSpecificStable(version string) error {
	stableURL := stableBaseURL + version + "/nvim-macos.tar.gz"
	targetDir := targetDirStable + version + "/"

	// Create the target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Download the Neovim archive
	resp, err := http.Get(stableURL)
	if err != nil {
		return fmt.Errorf("failed to download Neovim: %w", err)
	}
	defer resp.Body.Close()

	// Initialize progress bar
	progress := progress.New(progress.WithDefaultGradient())
	err = progress.Start()
	if err != nil {
		return fmt.Errorf("failed to initialize progress bar: %w", err)
	}

	// Save the downloaded file & track progress
	file, err := os.Create(filepath.Join(targetDir, "nvim-macos.tar.gz"))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	reader := progress.WrapReader(resp.Body) // Track download progress
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to save Neovim archive: %w", err)
	}

	// Complete the progress bar
	err = progress.Full()
	if err != nil {
		return fmt.Errorf("failed to complete progress bar: %w", err)
	}

	// Extract the archive
	extractCommand := fmt.Sprintf("tar xzvf %s -C %s", filepath.Join(targetDir, "nvim-macos.tar.gz"), targetDir)
	result, err := executeCommand(extractCommand) //  need a helper function for this
	if err != nil || result.ExitCode != 0 {
		return fmt.Errorf("failed to extract Neovim: %v", result.Output)
	}

	// Remove the downloaded archive
	err = os.Remove(filepath.Join(targetDir, "nvim-macos.tar.gz"))
	if err != nil {
		fmt.Println("Warning (non-fatal): Failed to remove Neovim archive:", err)
	}

	// Call your use_version function
	err = useVersion(version)
	if err != nil {
		return fmt.Errorf("failed to switch to version %s: %w", version, err)
	}

	fmt.Printf("Neovim version %s installed successfully!\n", version)
	return nil
}
