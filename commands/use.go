package commands

import (
	"fmt"
	"os"
)

func useVersion(version string) error {
	symlinkPath := "/usr/local/bin/nvim" // Assuming this path

	var neovimBinary string

	if version == "nightly" {
		// based on sort order, the latest nightly version will be the last one in the list
		neovimBinary = "helloniihgtly"
	} else {
		// Build the path for the specific version
		neovimBinary = targetDirStable + version + "/nvim-macos/bin/nvim"
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
