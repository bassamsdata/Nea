package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"nvm_manager_go/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	targetDirNightly = filepath.Join(homeDir, ".local", "share", "neoManager", "nightly")
	versionsFilePath = targetDirNightly + "/versions_info.json"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up old versions",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: You must specify 'nightly', 'stable', or 'all'")
			return
		}
		versionType := args[0]
		additionalArgs := args[1:]
		clean(versionType, additionalArgs)
	},
}

func clean(target string, options []string) error {
	switch {
	case target == "nightly" && len(options) == 0:
		return cleanLatestNightly()

	case target == "nightly" && len(options) > 0 && options[0] == "all":
		return cleanAllNightly()

	case target == "stable" && len(options) > 0 && options[0] == "all":
		return cleanAllStable()

	case target == "stable" && len(options) == 0:
		return cleanSpecificStable("stable")

	case target == "all":
		err1 := cleanAllNightly()
		err2 := cleanAllStable()
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
		return nil

	case strings.HasPrefix(target, "0."):
		return cleanSpecificStable(target)

	default:
		if strings.HasPrefix(target, "20") {
			return cleanSpecificNightly(target)
		} else {
			// TODO:: add proper error
			return fmt.Errorf("invalid version: %s", target)
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
	lastVersion := versions[0]

	// Delete the directory
	err = os.RemoveAll(lastVersion.Directory)
	if err != nil {
		return fmt.Errorf("failed to delete directory %s: %w", lastVersion.Directory, err)
	}
	// Remove the first element
	versions = append(versions[:0], versions[1:]...)

	// renumber the unique numbers again
	for i := range versions {
		versions[i].UniqueNumber = i
	}

	// Update the versions_info.json file
	updatedJson, err := json.MarshalIndent(versions, "", "  ")
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
	versions = versions[:0] // Empty the slice
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
func findNightlyVersion(versions []utils.VersionInfo, dateStr string) (int, bool) {
	for i, version := range versions {
		t, err := time.Parse(time.RFC3339, version.CreatedAt)
		if err != nil {
			log.Fatalf("Error parsing date for version %s: %v", version.CreatedAt, err)
		}
		if t.Format("2006-01-02") == dateStr {
			return i, true
		}
	}
	return -1, false // Not found
}

func cleanSpecificNightly(target string) error { // target is a date like 2022-12-07
	versions, err := utils.ReadVersionsInfo()
	if err != nil {
		return fmt.Errorf("failed to read versions info: %w", err)
	}

	index, found := findNightlyVersion(versions, target)
	if !found {
		fmt.Printf("Nightly version %s was not found.\n", target)
		return fmt.Errorf("nightly version %s not found", target)
	}
	fmt.Printf("Found version %s at index %d\n", target, index)

	// Delete the directory
	versionDir := versions[index].Directory
	err = os.RemoveAll(versionDir)
	if err != nil {
		return err
	}

	versions = append(versions[:index], versions[index+1:]...)

	// Update  versions_info.json
	updatedJson, _ := json.Marshal(versions)
	err = os.WriteFile(versionsFilePath, updatedJson, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Nightly version deleted successfully")
	return nil
}

func cleanSpecificStable(versionStr string) error {
	// if version is stable -> get the version no
	versionStr, err := utils.ResolveVersion(versionStr)
	if err != nil {
		return err
	}
	targetDir := filepath.Join(targetDirStable, versionStr)

	if _, err = os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("stable version %s not found", versionStr)
	}

	// Consider a confirmation prompt here
	err = os.RemoveAll(targetDir)
	if err != nil {
		return fmt.Errorf("failed to delete stable version %s: %w", versionStr, err)
	}

	fmt.Printf("Deleted stable version %s\n", versionStr)
	return nil
}
