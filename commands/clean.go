package commands

import (
	"encoding/json"
	"fmt"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
	"strings"
)

var (
	targetDirNightly = filepath.Join(homeDir, ".local", "share", "nv_manager", "nightly")
	versionsFilePath = targetDirNightly + "versions_info.json" // Use your actual path
)

type VersionInfo struct {
	NodeID       string `json:"node_id"`
	CreatedAt    string `json:"created_at"`
	Directory    string `json:"directory"`
	UniqueNumber int    `json:"unique_number"`
}

func clean(target string, options []string) error {
	switch {
	case target == "nightly" && len(options) == 0: // clean nightly (last version)
		return cleanLatestNightly()

	case target == "nightly" && options[0] == "all": // clean nightly all
		return cleanAllNightly()

	case target == "stable" && options[0] == "all": // clean stable all
		return cleanAllStable()

	case target == "all": // clean all
		err1 := cleanAllNightly()
		err2 := cleanAllStable()
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
		return nil

	default: // clean specific version (nightly or stable)
		if strings.HasPrefix(target, "0.") {
			return cleanSpecificStable(target)
		} else {
			return cleanSpecificNightly(target)
		}
	}
}
func cleanLatestNightly() error {
	versions, err := utils.ReadVersionsInfo()
	if err != nil {
		return fmt.Errorf("failed to read versions info: %w", err)
	}

	if len(versions) == 0 {
		return fmt.Errorf("no nightly versions installed")
	}

	// Get the last version
	lastVersion := versions[len(versions)-1]

	// Delete the directory
	err = os.RemoveAll(lastVersion.Directory)
	if err != nil {
		return fmt.Errorf("failed to delete directory %s: %w", lastVersion.Directory, err)
	}

	// Remove the entry from 'versions' (TODO)
	versions = versions[:len(versions)-1] // Remove the last element

	// Update the versions_info.json file (TODO)
	updatedJson, err := json.Marshal(versions)
	if err != nil {
		return fmt.Errorf("failed to marshal updated versions info: %w", err)
	}

	err = os.WriteFile(versionsFilePath, updatedJson, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated versions info: %w", err)
	}

	fmt.Printf("Deleted nightly version %s\n", lastVersion.CreatedAt)
	return nil
}

func cleanAllNightly() error {
	versions, err := utils.ReadVersionsInfo()
	if err != nil {
		return fmt.Errorf("failed to read versions info: %w", err)
	}

	for _, version := range versions {
		err = os.RemoveAll(version.Directory)
		if err != nil {
			return fmt.Errorf("failed to delete directory %s: %w", version.Directory, err)
		}
	}

	// Clear and update JSON
	versions = []VersionInfo{} // Empty the slice
	updatedJson, err := json.Marshal(versions)
	if err != nil {
		return fmt.Errorf("failed to marshal updated versions: %w", err)
	}
	err = os.WriteFile(versionsFilePath, updatedJson, 0644)
	if err != nil {
		return fmt.Errorf("failed to update versions info file: %w", err)
	}

	fmt.Println("Deleted all nightly versions.")
	return nil
}

func cleanAllStable() error {
	// Consider adding a confirmation prompt here
	err := os.RemoveAll(targetDirStable)
	if err != nil {
		return fmt.Errorf("failed to delete stable directory: %w", err)
	}
	fmt.Println("Deleted all stable versions.")
	return nil
}

// Helper to locate a nightly version by creation date
func findNightlyVersion(versions []VersionInfo, dateStr string) (int, bool) {
	for i, version := range versions {
		if version.Directory == dateStr {
			return i, true
		}
	}
	return -1, false // Not found
}

func cleanSpecificNightly(target string) error {
	versions, err := utils.ReadVersionsInfo()
	if err != nil {
		return fmt.Errorf("failed to read versions info: %w", err)
	}

	// FIX: why does this not work?
	index, found := findNightlyVersion(versions, target)
	if !found {
		return fmt.Errorf("nightly version %s not found", target)
	}

	// Delete the directory
	versionDir := versions[index].Directory
	err = os.RemoveAll(versionDir) // Error checking
	if err != nil {
		return err
	}

	// Update versions slice
	versions = append(versions[:index], versions[index+1:]...)

	// Update  versions_info.json
	updatedJson, _ := json.Marshal(versions)
	err = os.WriteFile(versionsFilePath, updatedJson, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted nightly version %s\n", target)
	return nil
}

func cleanSpecificStable(versionStr string) error {
	targetDir := filepath.Join(targetDirStable, versionStr)

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("stable version %s not found", versionStr)
	}

	// Consider a confirmation prompt here
	err := os.RemoveAll(targetDir)
	if err != nil {
		return fmt.Errorf("failed to delete stable version %s: %w", versionStr, err)
	}

	fmt.Printf("Deleted stable version %s\n", versionStr)
	return nil
}
