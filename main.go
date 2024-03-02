package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
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

type Tag struct {
	Name string `json:"name"`
}

type VersionInfo struct {
	NodeID       string `json:"node_id"`
	CreatedAt    string `json:"created_at"`
	Directory    string `json:"directory"`
	UniqueNumber int    `json:"unique_number"`
}
