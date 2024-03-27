package commands

import (
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
)

func useVersion(version string) error {
	symlinkPath := "/usr/local/bin/nvim"

	var neovimBinary string

	// build file path for latest installed nightly version
	if version == "nightly" {
		// TODO: Move it to its own function
		nightlyVersions, _ := utils.ReadVersionsInfo()
		if len(nightlyVersions) > 0 {
			VersionCreatedAt := nightlyVersions[0].CreatedAt
			neovimBinary = filepath.Join(targetDirNightly, VersionCreatedAt) + "/nvim-macos/bin/nvim"
		} else {
			return fmt.Errorf("no nightly versions installed")
		}
	} else {
		// Build the path for the stable version
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

	fmt.Printf("Currently using Neovim version %s\n", version)
	return nil
}
