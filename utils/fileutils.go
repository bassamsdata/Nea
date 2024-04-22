package utils

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var (
	homeDir          = os.Getenv("HOME")
	appDir           = filepath.Join(homeDir, ".local", "share", "neoManager")
	configPath       = filepath.Join(appDir, "config.json")
	neovimURL        = "https://github.com/neovim/neovim/releases/download/nightly/nvim-macos.tar.gz"
	targetNightly    = filepath.Join(homeDir, ".local", "share", "neoManager", "nightly")
	StableBaseURL    = "https://github.com/neovim/neovim/releases/download/v"
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

// Struct to represent release info
type Release struct {
	Name string `json:"name"`
}

// TODO: add number of releases
// Fetch stable neovim releases
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
	// Check if symlinked - I don't like this in Go
	if fi.Mode()&os.ModeSymlink == 0 {
		return "", fmt.Errorf("neovim is not symlinked")
	}

	symlinkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink target: %w", err)
	}

	parts := strings.Split(symlinkTarget, "/")

	// the main logic
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
	versionsFilePath := targetDirNightly + "versions_info.json"

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
	// Sort versions in descending order by CreatedAt
	slices.SortFunc(versions, func(a, b VersionInfo) int {
		timeA, _ := time.Parse(time.RFC3339, a.CreatedAt)
		timeB, _ := time.Parse(time.RFC3339, b.CreatedAt)
		return timeB.Compare(timeA) // NOTE: DESC order
	})

	return versions, nil
}

func CreateTargetDirectory(createdAt string) (string, error) {
	t, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return "", err
	}

	formattedDate := t.Format("2006-01-02")
	targetDir := filepath.Join(targetDirNightly, formattedDate)

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return "", err
	}

	return targetDir, nil
}

// TODO: Move to utils/fileutils.go
func DownloadArchive(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// Helper to extract a tar.gz archive
func ExtractTarGz(filePath, targetDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for header, err := tarReader.Next(); err == nil; header, err = tarReader.Next() {
		targetPath := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create target file: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close() // Close on error
				return fmt.Errorf("failed to extract file: %w", err)
			}
			if err := outFile.Close(); err != nil {
				return fmt.Errorf("failed to close file: %w", err)
			}
		default:
			return fmt.Errorf("unknown type: %b in %s", header.Typeflag, header.Name)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("error reading archive: %w", err)
	}
	return nil
}
