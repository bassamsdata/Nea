package commands

import (
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

var (
	homeDir       = os.Getenv("HOME")
	appDir        = filepath.Join(homeDir, ".local", "share", "neoManager")
	neovimURL     = "https://github.com/neovim/neovim/releases/download/nightly/nvim-macos.tar.gz"
	targetNightly = filepath.Join(homeDir, ".local", "share", "neoManager", "nightly")
	// stableBaseURL   = "https://github.com/neovim/neovim/releases/download/v"
	targetDirStable = filepath.Join(homeDir, ".local", "share", "neoManager", "stable/")
	tagsURL         = "https://api.github.com/repos/neovim/neovim/tags"
	versionFilePath = filepath.Join(targetNightly, "versions_info.json")
	nvm_night_url   = "https://github.com/neovim/neovim/releases/download/nightly/nvim-macos.tar.gz"
)
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a Neovim version",
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.Setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "nightly":
			err = installNightly()
			if err != nil {
				fmt.Println("Failed to install nightly:", err)
				return
			}
			// fmt.Println("You are now using the nightly version")

		case "stable":
			err = InstallSpecificStable("stable")
			if err != nil {
				fmt.Println("Failed to install latest stable:", err)
				return
			}
			fmt.Println("You are now using the latest stable version")

		default:
			err = InstallSpecificStable(args[0])
			if err != nil {
				fmt.Println("Failed to install version", args[0], ":", err)
				return
			}
			fmt.Println("Neovim version", args[0], "installed successfully")
		}
	},
}

func InstallSpecificStable(version string) error {
	stableURL := utils.StableBaseURL + version + "/nvim-macos.tar.gz"
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

	err = useVersion(version, nil)
	if err != nil {
		return fmt.Errorf("failed to switch to version %s: %w", version, err)
	}
	green := color.New(color.FgGreen).PrintfFunc()
	green("Neovim version %s installed successfully!\n", version)
	return nil
}
