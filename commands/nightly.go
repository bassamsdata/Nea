package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Release struct {
	NodeId    string `json:"node_id"`
	CreatedAt string `json:"created_at"`
}

func installNightly() error {
	// 0.get system and arch for file name and stop if os is not supported
	filename, err := getArchiveFilename("nightly")
	if err != nil {
		return err
	}

	// 1. Fetch Release Information
	latestRelease, err := fetchLatestNightlyRelease()
	if err != nil {
		return fmt.Errorf("failed to fetch release information: %w", err)
	}

	// 2. Check if Already Installed
	if isVersionInstalled(latestRelease.NodeId, latestRelease.CreatedAt) {
		color.Yellow("The latest nightly version is already installed.")
		color.Yellow("Use 'nvm use nightly' to switch to it.")
		return nil
	}

	// 3. Create Target Directory
	targetDir, err := utils.CreateTargetDirectory(latestRelease.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// 4. Download Archive
	archivePath := filepath.Join(targetDir, filename)
	buildURL := nvm_night_url + "/" + filename
	err = utils.DownloadArchive(buildURL, archivePath)
	if err != nil {
		return fmt.Errorf("failed to download Neovim: %w", err)
	}
	// SUG: we can change dir name here to nvim-macos

	osType := runtime.GOOS
	// 5. Extract Archive or set executable for AppImage
	switch osType {
	case "darwin":
		err = utils.ExtractTarGz(archivePath, targetDir)
		if err != nil {
			return fmt.Errorf("failed to extract Neovim: %w", err)
		}
	case "linux":
		err = os.Chmod(archivePath, 0755)
		if err != nil {
			return fmt.Errorf("failed to set executable permission: %w", err)
		}

		// Move the AppImage to the bin directory
		binDir := filepath.Join(appDir, "bin")
		finalPath := filepath.Join(binDir, "nvim") // rename to nvim
		err = os.Rename(archivePath, finalPath)
		if err != nil {
			return fmt.Errorf("failed to move AppImage to bin directory: %w", err)
		}

		// Create a symlink in /usr/local/bin
		//err = os.Symlink(finalPath, "/usr/local/bin/nvim")
		//if err != nil {
		//	fmt.Println("Warning: failed to create symlink in /usr/local/bin/nvim:", err)
		//}

	}

	// 6. Remove Archive
	if osType == "darwin" {
		err = os.Remove(archivePath)
		if err != nil {
			fmt.Println("Warning: failed to remove archive:", err)
		}
	}

	// 8. Update versions_info.json
	err = updateVersionsInfo(latestRelease, targetDir)
	if err != nil {
		return fmt.Errorf("failed to update versions info: %w", err)
	}

	// 9. Call useVersion function
	err = useVersion("nightly", nil)
	if err != nil {
		return fmt.Errorf("failed to use nightly version: %w", err)
	}

	// 10. Success message
	color.Green("Neovim nightly installed successfully!")
	color.Green("and you are using Neovim nightly created on %s", latestRelease.CreatedAt)

	// Create symlink
	// Reuse osType from earlier in the function
	binDir := filepath.Join(appDir, "bin")

	// For macOS, we need to find the actual nvim binary inside the extracted app
	var nvimBinaryPath string
	if osType == "darwin" {
		// After extraction, the structure is likely nvim-osx64/bin/nvim or similar
		// First try to locate the executable using a glob pattern
		matches, err := filepath.Glob(filepath.Join(targetDir, "*/bin/nvim"))
		if err == nil && len(matches) > 0 {
			// Use the first match
			nvimBinaryPath = matches[0]
		} else {
			// Try alternative paths
			alternativePaths := []string{
				filepath.Join(targetDir, "nvim-macos", "bin", "nvim"),
				filepath.Join(targetDir, "bin", "nvim"),
			}

			for _, path := range alternativePaths {
				if _, err := os.Stat(path); err == nil {
					nvimBinaryPath = path
					break
				}
			}

			// If we still couldn't find it, log a detailed error
			if nvimBinaryPath == "" {
				// Print directory contents for debugging
				fmt.Println("DEBUG: Could not find nvim binary. Directory contents:")
				printDirContents(targetDir)
				return fmt.Errorf("could not locate nvim binary in extracted directory")
			}
		}
	} else {
		// For Linux, the binary is directly in the bin directory
		nvimBinaryPath = filepath.Join(binDir, "nvim")
	}

	// Remove existing symlink
	symlinkPath := filepath.Join(appDir, "bin", "nvim")
	err = os.Remove(symlinkPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing symlink: %w", err)
	}

	// Create new symlink
	err = os.Symlink(nvimBinaryPath, symlinkPath)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// NOTE: this is just for nightly veriosn currently
func getArchiveFilename(version string) (string, error) {
	osType := runtime.GOOS
	arch := runtime.GOARCH

	switch osType {
	case "darwin":
		if version == "nightly" || compareVersions(version, "0.10.0") >= 0 {
			switch arch {
			case "amd64":
				return "nvim-macos-x86_64.tar.gz", nil
			case "arm64":
				return "nvim-macos-arm64.tar.gz", nil
			default:
				return "", fmt.Errorf("unsupported architecture: %s", arch)
			}
		} else {
			return "nvim-macos.tar.gz", nil
		}
	case "linux":
		switch arch {
		case "amd64":
			// For nightly we use AppImage, but for stable we should use tar.gz
			if version == "nightly" {
				return "nvim-linux-x86_64.appimage", nil
			} else {
				return "nvim-linux64.tar.gz", nil
			}
		case "arm64":
			// Assuming arm64 exists. May need to verify
			if version == "nightly" {
				return "nvim-linux-arm64.appimage", nil
			} else {
				return "nvim-linux-arm64.tar.gz", nil
			}
		default:
			return "", fmt.Errorf("unsupported architecture: %s", arch)
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %s. currently only support macOS and Linux but PR is welcome", osType)
	}
}

func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	maxLen := max(len(v2Parts), len(v1Parts))

	for i := range maxLen {
		partV1, partV2 := 0, 0
		if i < len(v1Parts) {
			partV1, _ = strconv.Atoi(v1Parts[i])
		}
		if i < len(v2Parts) {
			partV2, _ = strconv.Atoi(v2Parts[i])
		}

		if partV1 < partV2 {
			return -1
		} else if partV1 > partV2 {
			return 1
		}
	}

	return 0
}

func updateVersionsInfo(latestRelease Release, targetDir string) error {
	versionsInfo, err := utils.ReadVersionsInfo()
	if err != nil {
		return fmt.Errorf("failed to read versions info: %w", err)
	}

	// Read the version limit from the config file
	config, err := utils.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Check if the number of versions exceeds the limit
	if len(versionsInfo) >= config.RollbackLimit {
		// Remove the oldest version -- we don't need to sort since ReadVersionsInfo already does that
		oldestVersion := versionsInfo[len(versionsInfo)-1]
		versionsInfo = versionsInfo[:len(versionsInfo)-1]

		// Delete the corresponding directory
		// FIX: here is something wrong, it should use the oldestVersion.targetDir
		dirToDelete := filepath.Join(targetDir, oldestVersion.CreatedAt[:10])
		if err = os.RemoveAll(dirToDelete); err != nil {
			return fmt.Errorf("failed to delete directory: %w", err)
		}
	}

	// Create the new VersionInfo
	newVersion := utils.VersionInfo{
		NodeID:       latestRelease.NodeId,
		CreatedAt:    latestRelease.CreatedAt,
		Directory:    targetDir,
		UniqueNumber: 0,
	}

	// Append the new entry to the slice
	versionsInfo = append(versionsInfo, newVersion)

	utils.SortVersionsDesc(versionsInfo)

	// re-assign unique number
	for i := range versionsInfo {
		versionsInfo[i].UniqueNumber = i
	}

	// Marshal back into JSON
	updatedJson, err := json.Marshal(versionsInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal versions info: %w", err)
	}

	// Write to versions_info.json
	err = os.WriteFile(versionFilePath, updatedJson, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write versions info: %w", err)
	}

	return nil
}

func fetchLatestNightlyRelease() (Release, error) {
	const tagsNightlyUrl = "https://api.github.com/repos/neovim/neovim/releases/tags/nightly"

	resp, err := http.Get(tagsNightlyUrl)
	if err != nil {
		return Release{}, err // Could wrap the error here
	}
	defer resp.Body.Close()

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return Release{}, err
	}

	return release, nil
}

func isVersionInstalled(nodeId, createdAt string) bool {
	versionsInfo, err := utils.ReadVersionsInfo()
	if err != nil {
		return false
	}

	for _, version := range versionsInfo {
		// Check if the exact same node ID is installed
		if version.NodeID == nodeId {
			return true
		}
	}
	return false
}
