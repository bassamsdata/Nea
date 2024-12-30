package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nvm_manager_go/utils"

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

	// Resolve version (stable -> actual version number)
	version, err = utils.ResolveVersion(version)
	if err != nil {
		return err
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
			versions, err := utils.ReadVersionsInfo()
			if err != nil || len(versions) == 0 {
				return fmt.Errorf("no nightly versions installed")
			}
			neovimBinary = filepath.Join(versions[0].Directory, dirName, "bin/nvim")
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
