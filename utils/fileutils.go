package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// Struct to represent release info
type Release struct {
	Name string `json:"name"`
}

func fetchReleases() ([]Release, error) {
	// 1. Make the HTTP Request
	resp, err := http.Get(tagsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	// 2. Read the Response Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 3. Parse JSON
	var releases []Release
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return releases, nil
}

// NOTE: Prod-ready function
func DetermineCurrentVersion() (string, error) {
	symlinkPath := "/usr/local/bin/nvim"

	fi, err := os.Lstat(symlinkPath)
	if err != nil {
		return "", fmt.Errorf("neovim is not symlinked: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return "", fmt.Errorf("neovim is not symlinked")
	}

	symlinkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink target: %w", err)
	}

	parts := strings.Split(symlinkTarget, "/")

	if strings.Contains(symlinkTarget, "nightly") {
		for _, part := range parts {
			if strings.HasPrefix(part, "20") {
				return part, nil // Return the nightly directory name as of `2022-03-31`
			}
		}
	} else if strings.Contains(symlinkTarget, "stable") {
		for _, part := range parts {
			if strings.HasPrefix(part, "0.") {
				return part, nil // Return the stable version number as of `0.9.5`
			}
		}
	}

	return "", fmt.Errorf("failed to parse version from symlink target")
}

// read nightly versions info
func ReadVersionsInfo() ([]VersionInfo, error) {
	versionsFilePath := targetDirNightly + "versions_info.json" // Use your actual path

	// Read the file contents
	data, err := os.ReadFile(versionsFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions info file: %w", err)
	}

	// Parse the JSON
	var versions []VersionInfo
	err = json.Unmarshal(data, &versions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse versions info JSON: %w", err)
	}

	return versions, nil
}
