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
	if isVersionInstalled(latestRelease.NodeId) {
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

	// 5. Extract Archive
	err = utils.ExtractTarGz(archivePath, targetDir)
	if err != nil {
		return fmt.Errorf("failed to extract Neovim: %w", err)
	}

	// 6. Remove Archive
	err = os.Remove(archivePath)
	if err != nil {
		fmt.Println("Warning: failed to remove archive:", err)
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
	return nil
}

// NOTE: this is just for nightly veriosn currently
func getArchiveFilename(version string) (string, error) {
	os := runtime.GOOS
	arch := runtime.GOARCH

	if os != "darwin" {
		return "", fmt.Errorf("unsupported operating system: %s. currently only support macOS but PR is welcome", os)
	}

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
}

func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	for i := 0; i < maxLen; i++ {
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

func isVersionInstalled(nodeId string) bool {
	versionsInfo, err := utils.ReadVersionsInfo()
	if err != nil {
		// TODO: handle error
		return false
	}

	for _, version := range versionsInfo {
		if version.NodeID == nodeId {
			return true
		}
	}
	return false
}
