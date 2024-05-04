package commands

import (
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
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
	symlinkPath := "/usr/local/bin/nvim"

	var neovimBinary string
	if version == "nightly" {
		switch {
		// TODO: I'm not sure why I did this, it's creating a bug where the
		// symlink won't work in install nightly
		case optionalDir != nil:
			neovimBinary = *optionalDir
		default:
			versions, _ := utils.ReadVersionsInfo() // already sorted
			if len(versions) > 0 {
				// Get the archive filename based on the current architecture
				archiveFilename, err := getArchiveFilename()
				if err != nil {
					return fmt.Errorf("failed to determine archive filename: %w", err)
				}
				// Extract the directory name from the archive filename
				dirName := strings.TrimSuffix(archiveFilename, ".tar.gz")
				neovimBinary = filepath.Join(versions[0].Directory, dirName, "bin/nvim")
			} else {
				return fmt.Errorf("no nightly versions installed")
			}
		}
	} else {
		// Build the path for the specific version
		neovimBinary = filepath.Join(targetDirStable, version) + "/nvim-macos/bin/nvim"
	}

	if _, err := os.Stat(neovimBinary); err != nil {
		return fmt.Errorf("version %s is not installed: %w", version, err)
	}

	if err := os.Remove(symlinkPath); err != nil {
		return fmt.Errorf("failed to remove existing symlink: %w", err)
	}

	if err := os.Symlink(neovimBinary, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	// TODO: use another function for use, keep this as a helper function
	// so no need to print any message
	// fmt.Printf("Currently using Neovim version %s\n", version)
	return nil
}
