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
	appDir           = filepath.Join(homeDir, ".local", "share", "nv_manager")
	neovimURL        = "https://github.com/neovim/neovim/releases/download/nightly/nvim-macos.tar.gz"
	targetNightly    = filepath.Join(homeDir, ".local", "share", "nv_manager", "nightly")
	stableBaseURL    = "https://github.com/neovim/neovim/releases/download/v"
	targetDirNightly = filepath.Join(homeDir, ".local", "share", "nv_manager", "nightly")
	targetDirStable  = filepath.Join(homeDir, ".local", "share", "nv_manager", "stable")
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
	var rootCmd = &cobra.Command{
		Use:   "nvimv",
		Short: "Neovim version manager",
		Long:  `Neovim version manager allows you to install, update, and switch between different versions of Neovim.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Main command logic here
			fmt.Println("Use \"nvimv install nightly\" to install the latest nightly build.")
		},
	}

	var installCmd = &cobra.Command{
		Use:   "install [version]",
		Short: "Install the latest nightly or a specific version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// version := args[0]
			// if version == "nightly" {
			// installNightly()
			// } else {
			// installSpecificStable(version)
			// }
		},
	}
	rootCmd.AddCommand(installCmd)

	// Update Commands
	updateCmd := &cobra.Command{
		Use:   "update [nightly]",
		Short: "Update a Neovim version",
		Args:  cobra.MaximumNArgs(1),
		// Run:   updateHandler(),
	}
	rootCmd.AddCommand(updateCmd)

	// Update Commands
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup directories and files",
		Args:  cobra.ExactArgs(0),
		// Run:   setup,
	}
	rootCmd.AddCommand(setupCmd)

	// List Commands
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed Neovim versions",
		Run:   commands.ListHandlerLocal,
	}
	rootCmd.AddCommand(listCmd)

	// Use Command
	useCmd := &cobra.Command{
		Use:   "use [nightly | latest | version]",
		Short: "Switch to a specific Neovim version",
		Args:  cobra.ExactArgs(1),
		// Run:   useHandler(),
	}
	rootCmd.AddCommand(useCmd)

	// Clean Commands
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove Neovim versions",
	}

	cleanNightlyCmd := &cobra.Command{
		Use:   "nightly [all]",
		Short: "Clean nightly versions",
		Args:  cobra.MaximumNArgs(1),
		// Run:   cleanNightlyHandler,
	}
	cleanCmd.AddCommand(cleanNightlyCmd)

	cleanStableCmd := &cobra.Command{
		Use:   "stable [version | all]",
		Short: "Clean stable versions",
		Args:  cobra.MinimumNArgs(1),
		// Run:   cleanStableHandler,
	}
	cleanCmd.AddCommand(cleanStableCmd)

	rootCmd.AddCommand(cleanCmd)

	// Other Commands
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Check the currently active Neovim version",
		// Run:   checkHandler,
	}
	rootCmd.AddCommand(checkCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Define your functions here, such as installNightly() and installSpecificStable(version string)
// These functions should contain the logic to install Neovim versions as per your Vlang code.

// TODO: need review and probably, should create all dirs here
func setup() {
	versionFilePath := filepath.Join(targetNightly, "versions_info.json")
	if _, err := os.Stat(versionFilePath); !os.IsNotExist(err) {
		fmt.Println("The versions_info.json file already exists.")
		return
	}

	// Define the initial content for the versions_info.json file
	initialContent := "[]" // Start with an empty JSON array

	// Write the initial content to the file
	err := os.WriteFile(versionFilePath, []byte(initialContent), 0644)
	if err != nil {
		fmt.Printf("Failed to create versions_info.json file: %s\n", err)
		return
	}

	fmt.Println("The versions_info.json file has been created successfully.")
}
