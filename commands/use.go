package commands

import (
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
)

func useVersion(version string, optionalDir *string) error {
	symlinkPath := "/usr/local/bin/nvim"

	var neovimBinary string
	if version == "nightly" {
		switch {
		case optionalDir != nil:
			neovimBinary = *optionalDir
		default:
			versions, _ := utils.ReadVersionsInfo() // already sorted
			if len(versions) > 0 {
				neovimBinary = versions[0].Directory
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
