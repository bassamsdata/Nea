package commands

import (
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var RollbackCmd = &cobra.Command{
	Use:   "rollback [steps]",
	Short: "Rollback to a previous version",
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		// Parse the first argument as an integer
		rollbackStep, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: The argument must be a number")
			return
		}

		// Call RollbackVersion with the parsed integer
		err = RollbackVersion(rollbackStep)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Rollback successful")
	},
}

// either use int number or date
func RollbackVersion(rollbackStep int) error {
	// symlinkPath := "/usr/local/bin/nvim"
	var neovimBinary string

	// 1. Read versions_info.json
	versionsInfo, err := utils.ReadVersionsInfo()
	if err != nil {
		return fmt.Errorf("failed to read versions info: %w", err)
	}

	// 2. Check if rollbackStep is valid
	if rollbackStep < 1 {
		return fmt.Errorf("rollback steps must be at least 1. Or to use latest nightly version, use 'use nightly'")
	}

	// 3. Check if rollbackStep is valid
	if rollbackStep >= len(versionsInfo) {
		return fmt.Errorf("cannot rollback %d versions, not enough versions installed", rollbackStep)
	}

	// 4. Get the version to rollback to
	rollbackTarget := versionsInfo[rollbackStep]

	// 5. detect filename of the archive
	binaryNameArch, err := getArchiveFilename("nightly")
	if err != nil {
		return fmt.Errorf("failed to get filename: %w", err)
	}
	dirName := strings.TrimSuffix(binaryNameArch, ".tar.gz")

	// why not using usevwrsion function
	// 5. Use the version
	neovimBinary = filepath.Join(rollbackTarget.Directory, dirName) + "/bin/nvim"
	if _, err := os.Stat(neovimBinary); err != nil {
		return fmt.Errorf("version %s is not installed: %w", rollbackTarget.CreatedAt, err)
	}
	if err := os.Remove(utils.SymlinkPath); err != nil {
		return fmt.Errorf("failed to remove existing symlink: %w", err)
	}
	if err := os.Symlink(neovimBinary, utils.SymlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}
