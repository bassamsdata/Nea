package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type Release struct {
	NodeId    string `json:"node_id"`
	CreatedAt string `json:"created_at"`
	// Other release fields if you need them
}

func installNightly() error {
	// 1. Fetch Release Information
	latestRelease, err := fetchLatestNightlyRelease()
	if err != nil {
		return fmt.Errorf("failed to fetch release information: %w", err)
	}

	// 2. Check if Already Installed
	if isVersionInstalled(latestRelease.NodeId) {
		color.Yellow("The latest nightly version is already installed.")
		return nil
	}

	// 3. Create Target Directory
	targetDir, err := utils.CreateTargetDirectory(latestRelease.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// 4. Download Archive
	archivePath := filepath.Join(targetDir, "nvim-macos.tar.gz")
	err = utils.DownloadArchive(nvm_night_url, archivePath)
	if err != nil {
		return fmt.Errorf("failed to download Neovim: %w", err)
	}

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

	// 7. Update versions_info.json
	err = updateVersionsInfo(latestRelease, targetDir)
	if err != nil {
		return fmt.Errorf("failed to update versions info: %w", err)
	}

	// 8. Success message
	color.Green("Neovim nightly installed successfully!")
	return nil
}

func updateVersionsInfo(latestRelease Release, targetDir string) error {
	versionsInfo, err := utils.ReadVersionsInfo()
	if err != nil {
		return fmt.Errorf("failed to read versions info: %w", err)
	}

	// Add logic to remove oldest versions if count is >= 3
	if len(versionsInfo) >= 3 {
		// Remove oldest version entry from versionsInfo slice
		versionsInfo = versionsInfo[1:]

		// TODO: I'm not sure but we should sort here probably
		// Renumber the remaining entries
		for i := range versionsInfo {
			versionsInfo[i].UniqueNumber = i + 1
		}
	}

	// Create the new VersionInfo
	newVersion := utils.VersionInfo{
		NodeID:       latestRelease.NodeId,
		CreatedAt:    latestRelease.CreatedAt,
		Directory:    targetDir,
		UniqueNumber: len(versionsInfo) + 1,
	}

	// Append the new entry to the slice
	versionsInfo = append(versionsInfo, newVersion)

	// Marshal back into JSON
	updatedJson, err := json.Marshal(versionsInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal versions info: %w", err)
	}

	// Write to versions_info.json
	// TODO: I'm not sure if this is the path
	err = os.WriteFile(versionFilePath, updatedJson, 0644)
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
