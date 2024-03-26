package commands

import (
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
)

var (
	homeDir         = os.Getenv("HOME")
	appDir          = filepath.Join(homeDir, ".local", "share", "nv_manager")
	neovimURL       = "https://github.com/neovim/neovim/releases/download/nightly/nvim-macos.tar.gz"
	targetNightly   = filepath.Join(homeDir, ".local", "share", "nv_manager", "nightly")
	stableBaseURL   = "https://github.com/neovim/neovim/releases/download/v"
	targetDirStable = filepath.Join(homeDir, ".local", "share", "nv_manager", "stable/")
	tagsURL         = "https://api.github.com/repos/neovim/neovim/tags"
	versionFilePath = filepath.Join(targetNightly, "versions_info.json")
	nvm_night_url   = "https://github.com/neovim/neovim/releases/download/nightly/nvim-macos.tar.gz"
)

func InstallSpecificStable(version string) error {
	stableURL := stableBaseURL + version + "/nvim-macos.tar.gz"
	targetDir := filepath.Join(targetNightly, version)

	// Create the target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Download the Neovim archive
	archivePath := filepath.Join(targetDir, "nvim-macos.tar.gz")
	if err := utils.DownloadArchive(stableURL, archivePath); err != nil {
		return fmt.Errorf("failed to download Neovim: %w", err)
	}

	// Extract the archive
	err := utils.ExtractTarGz(archivePath, targetDir)
	if err != nil {
		return fmt.Errorf("failed to extract Neovim: %w", err)
	}

	// Remove the downloaded archive
	err = os.Remove(archivePath)
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
