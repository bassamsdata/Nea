package commands

import (
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
	nvm_night_url   = "https://github.com/neovim/neovim/releases/download/nightly/"
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
		}
	},
}

func InstallSpecificStable(version string) error {
	startTime := time.Now()
	defer func() { fmt.Printf("Total execution time: %v\n", time.Since(startTime)) }()

	// if version is stable -> get the version no
	version, err := utils.ResolveVersion(version)
	if err != nil {
		return err
	}

	targetDir := filepath.Join(targetDirStable, version)
	_, err = os.Stat(targetDir)
	if err == nil || !os.IsNotExist(err) {
		fmt.Println("Version", version, "is already installed.")
		return nil
	}

	// Determine the correct archive filename based on the version
	archiveFilename, err := getArchiveFilename(version)
	if err != nil {
		return fmt.Errorf("failed to determine archive filename: %w", err)
	}

	stableURL := utils.StableBaseURL + version + "/" + archiveFilename

	// 1. Create the target directory
	if err = os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// 2. Download the Neovim archive
	archivePath := filepath.Join(targetDir, archiveFilename)
	if err = utils.DownloadArchive(stableURL, archivePath); err != nil {
		return fmt.Errorf("failed to download Neovim: %w", err)
	}

	// 3. Extract the archive
	err = utils.ExtractTarGz(archivePath, targetDir)
	if err != nil {
		return fmt.Errorf("failed to extract Neovim: %w", err)
	}

	// 4. Remove the downloaded archive
	err = os.Remove(archivePath)
	if err != nil {
		fmt.Println("Warning (non-fatal): Failed to remove Neovim archive:", err)
	}

	err = useVersion(version, nil)
	if err != nil {
		return fmt.Errorf("failed to switch to version %s: %w", version, err)
	}
	green := color.New(color.FgCyan).PrintfFunc()
	green("Neovim version %s installed successfully!\n", version)
	return nil
}
