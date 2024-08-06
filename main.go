package main

import (
	"fmt"
	"nvm_manager_go/commands"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	homeDir          = os.Getenv("HOME")
	appDir           = filepath.Join(homeDir, ".local", "share", "neoManager")
	neovimURL        = "https://github.com/neovim/neovim/releases/download/nightly/nvim-macos.tar.gz"
	targetNightly    = filepath.Join(homeDir, ".local", "share", "neoManager", "nightly")
	stableBaseURL    = "https://github.com/neovim/neovim/releases/download/v"
	targetDirNightly = filepath.Join(homeDir, ".local", "share", "neoManager", "nightly")
	targetDirStable  = filepath.Join(homeDir, ".local", "share", "neoManager", "stable")
	tagsURL          = "https://api.github.com/repos/neovim/neovim/tags"
	tagsNightlyURL   = "https://api.github.com/repos/neovim/neovim/releases/tags/nightly"
)

type VersionInfo struct {
	NodeID       string `json:"node_id"`
	CreatedAt    string `json:"created_at"`
	Directory    string `json:"directory"`
	UniqueNumber int    `json:"unique_number"`
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "nvm",
		Short: "Neovim Version Manager (Go)",
	}

	rootCmd.AddCommand(commands.InstallCmd)
	rootCmd.AddCommand(commands.UseCmd)
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.RollbackCmd)
	rootCmd.AddCommand(commands.CleanCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
