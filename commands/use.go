package commands

import (
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var UseCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a Neovim version",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: You must specify 'nightly', a version number, or 'stable'")
			return
		}
		useFinalHandler(args)
	},
}

func useFinalHandler(args []string) {
	switch args[0] {
	case "nightly":
		err := useVersion("nightly", nil)
		if err != nil {
			fmt.Println("Failed to use version:", err)
			return
		} else {
			fmt.Println("Nightly version used successfully")
		}
	case "stable":
		err := useVersion("stable", nil)
		if err != nil {
			fmt.Println("Failed to use version:", err)
			return
		} else {
			fmt.Println("Latest stable version used successfully")
		}
	default:
		err := useVersion(args[0], nil)
		if err != nil {
			fmt.Println("Failed to use version:", err)
			return
		} else {
			fmt.Println("Version", args[0], "used successfully")
		}
	}
}

func useVersion(version string, optionalDir *string) error {
	// Get base paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	baseDir := filepath.Join(homeDir, ".local", "share", "neoManager")
	binDir := filepath.Join(baseDir, "bin")
	symlinkPath := filepath.Join(binDir, "nvim")

	// Ensure bin directory exists
	if mkdirErr := os.MkdirAll(binDir, 0o755); mkdirErr != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Handle stable version differently
	if version == "stable" {
		// First check local versions
		localVersions, localErr := utils.GetLocalStableVersions()
		if localErr != nil {
			return fmt.Errorf("failed to get local versions: %w", localErr)
		}

		if len(localVersions) == 0 {
			return fmt.Errorf("no stable versions installed. Run 'nvm install stable' first")
		}

		// Use the latest local version
		version = localVersions[0]

		// Check if newer version is available online
		latestOnline, onlineErr := utils.FetchLatestStable()
		if onlineErr == nil && version != latestOnline { // Only show warning if we can fetch latest version
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("%s\n", yellow(fmt.Sprintf("Note: A newer version (%s) is available. Run 'nvm install stable' to get it.", latestOnline)))
		}
	} else if version != "nightly" {
		// For specific versions, resolve version number
		resolvedVersion, resolveErr := utils.ResolveVersion(version)
		if resolveErr != nil {
			return resolveErr
		}
		version = resolvedVersion
	}

	// Get the correct archive filename for the architecture
	archiveFilename, err := getArchiveFilename(version)
	if err != nil {
		return fmt.Errorf("failed to determine archive filename: %w", err)
	}

	// Extract the directory name from the archive filename
	dirName := strings.TrimSuffix(archiveFilename, ".tar.gz")

	var neovimBinary string
	if version == "nightly" {
		if optionalDir != nil {
			neovimBinary = filepath.Join(*optionalDir, "bin/nvim")
		} else {
			versions, readErr := utils.ReadVersionsInfo()
			if readErr != nil || len(versions) == 0 {
				return fmt.Errorf("no nightly versions installed")
			}
			
			// Check if the directory exists directly
			versionPath := versions[0].Directory
			nvimBinPath := filepath.Join(versionPath, dirName, "bin/nvim")
			
			// Check if the binary exists at the expected path
			if _, err := os.Stat(nvimBinPath); err == nil {
				neovimBinary = nvimBinPath
			} else {
				// Handle case where directory structure might be different
				// This could happen with the new timestamp-based directories
				entries, err := os.ReadDir(versionPath)
				if err != nil {
					return fmt.Errorf("failed to read nightly directory: %w", err)
				}
				
				// Look for nvim directory within the version directory
				for _, entry := range entries {
					if entry.IsDir() && strings.Contains(entry.Name(), "nvim") {
						possibleBin := filepath.Join(versionPath, entry.Name(), "bin/nvim")
						if _, err := os.Stat(possibleBin); err == nil {
							neovimBinary = possibleBin
							break
						}
					}
				}
				
				// If we still don't have a binary path, return an error
				if neovimBinary == "" {
					return fmt.Errorf("neovim binary not found in %s", versionPath)
				}
			}
		}
	} else {
		neovimBinary = filepath.Join(targetDirStable, version, dirName, "bin", "nvim")
	}

	// Verify the binary exists
	if _, err := os.Stat(neovimBinary); err != nil {
		return fmt.Errorf("neovim binary not found at %s: %w", neovimBinary, err)
	}

	// Remove existing symlink if it exists
	if err := os.RemoveAll(symlinkPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing symlink: %w", err)
	}

	// Create new symlink
	if err := os.Symlink(neovimBinary, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}
