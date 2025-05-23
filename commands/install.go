package commands

import (
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
	"strings"
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
	Long: `Install a Neovim version. Valid formats:
- nightly: Latest nightly build
- stable: Latest stable version
- x.y.z: Specific version (e.g., 0.9.5)`,
	Args: cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.Setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

		var err error
		if version == "nightly" {
			err = installNightly()
			if err != nil {
				fmt.Println("Failed to install nightly:", err)
				return
			}
		} else {
			// Validate and resolve version before proceeding
			resolvedVersion, err := utils.ResolveVersion(version)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			err = InstallSpecificStable(resolvedVersion)
			if err != nil {
				fmt.Printf("Failed to install version %s: %v\n", resolvedVersion, err)
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

	// 3. Extract the archive or set executable permissions for AppImage
	if strings.HasSuffix(archiveFilename, ".tar.gz") {
		// Extract tar.gz archive
		err = utils.ExtractTarGz(archivePath, targetDir)
		if err != nil {
			return fmt.Errorf("failed to extract Neovim: %w", err)
		}

		// 4. Remove the downloaded archive
		err = os.Remove(archivePath)
		if err != nil {
			fmt.Println("Warning (non-fatal): Failed to remove Neovim archive:", err)
		}
	} else if strings.HasSuffix(archiveFilename, ".appimage") {
		// Set executable permissions for AppImage
		err = os.Chmod(archivePath, 0755)
		if err != nil {
			return fmt.Errorf("failed to set executable permission: %w", err)
		}

		// Move the AppImage to the bin directory
		binDir := filepath.Join(appDir, "bin")
		finalPath := filepath.Join(binDir, "nvim") // rename to nvim
		err = os.Rename(archivePath, finalPath)
		if err != nil {
			return fmt.Errorf("failed to move AppImage to bin directory: %w", err)
		}
	}

	err = useVersion(version, nil)
	if err != nil {
		return fmt.Errorf("failed to switch to version %s: %w", version, err)
	}
	green := color.New(color.FgCyan).PrintfFunc()
	green("Neovim version %s installed successfully!\n", version)
	return nil
}
