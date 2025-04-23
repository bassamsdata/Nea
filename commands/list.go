package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nvm_manager_go/utils"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type Tag struct {
	Name string `json:"name"`
}

// Define a struct to hold version info for the table
type VersionData struct {
	Version string
	Status  string
}

var ListCmd = &cobra.Command{
	Use:   "ls [local|remote] [count]",
	Short: "List all Neovim versions",
	Long: `List Neovim versions.

Usage:
  nvm ls local [count]  - List local installed versions
                         (By default shows all stable versions and 5 most recent nightly versions)
                         (Optional: specify max count of nightly versions to show)
                         (Use -1 to show ALL nightly versions)
  nvm ls remote [count] - List remote available versions
                         (By default shows 7 most recent versions)
                         (Optional: specify max count to show)
                         (Use -1 to show ALL available versions)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: You must specify local or remote")
			return
		}

		// Default count
		count := 0 // 0 means all versions

		// Check if a count argument is provided
		if len(args) >= 2 {
			requestedCount, err := strconv.Atoi(args[1])
			if err == nil {
				// Allow special value -1 to show all versions
				if requestedCount == -1 {
					count = -1
				} else if requestedCount > 0 {
					count = requestedCount
				}
			}
		}

		switch args[0] {
		case "local":
			listHandlerLocal(count)
		case "remote":
			listHandler(args[1:])
		}
	},
}

// Display remote available versions
func listHandler(args []string) {
	numVersions := 7 // Default number of versions to list

	// Parse count from args if present
	if len(args) >= 1 {
		requestedVersions, err := strconv.Atoi(args[0])
		if err == nil {
			// Special value -1 means show all versions
			if requestedVersions == -1 {
				numVersions = -1 // No limit
			} else if requestedVersions > 0 {
				numVersions = requestedVersions
			}
		}
	}

	// Fetch the tags from GitHub API
	resp, err := http.Get("https://api.github.com/repos/neovim/neovim/tags")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch Neovim versions:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var tags []Tag
	err = json.Unmarshal(body, &tags)
	if err != nil {
		fmt.Println("Failed to decode JSON:", err)
		return
	}

	// Create a nice table for the output
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Version", "Type"})

	// Add nightly version at the top
	table.Append([]string{"nightly", "Development build"})

	// Add the stable versions
	displayCount := 0
	for _, tag := range tags {
		// If numVersions is -1, show all versions
		// Otherwise, break if we've shown enough versions
		if numVersions != -1 && displayCount >= numVersions {
			break
		}

		// Only show versions that match semantic versioning pattern (e.g., v0.9.2)
		if strings.HasPrefix(tag.Name, "v0.") || strings.HasPrefix(tag.Name, "v1.") {
			table.Append([]string{tag.Name, "Stable release"})
			displayCount++
		}
	}

	// Show total count if limited and not showing all
	stableCount := 0
	for _, tag := range tags {
		if strings.HasPrefix(tag.Name, "v0.") || strings.HasPrefix(tag.Name, "v1.") {
			stableCount++
		}
	}

	if numVersions != -1 && stableCount > displayCount {
		table.Append([]string{"", ""})
		table.Append([]string{fmt.Sprintf("Showing %d of %d available versions", displayCount, stableCount), ""})
	}

	table.Render()
	fmt.Println(tableString.String())
}

func listHandlerLocal(count int) {
	versions, err := utils.ReadVersionsInfo()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read versions info:", err)
		return
	}

	// If count wasn't explicitly specified, default to showing only 5 nightly versions
	if count == 0 {
		count = 7
	}

	// BUG: this doesn't work with the current state of it
	currentVersion, err := utils.DetermineCurrentVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to determine current version:", err)
		return
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Version", "Created At", "Rollback Step", "Status"})

	// Read the stable versions from the directory
	stableVersions, err := os.ReadDir(targetDirStable)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: Failed to read stable versions directory:", err)
	} else {
		// Sort the stable version names from latest to oldest using semver.
		sort.Slice(stableVersions, func(i, j int) bool {
			vi := stableVersions[i].Name()
			vj := stableVersions[j].Name()
			// Prepend "v" to comply with semver.Compare requirements.
			return semver.Compare("v"+vi, "v"+vj) > 0
		})
		for _, version := range stableVersions {
			versionName := version.Name()
			status := "stable"
			if versionName == currentVersion {
				status = "used"
			}
			table.Append([]string{versionName, "", "N/A", status})
		}
	}

	// Group versions by date to detect multiple nightlies on the same day
	dateMap := make(map[string]int)

	// Limit the number of nightly versions to display
	displayVersions := versions

	// Special value -1 means show all
	if count == -1 {
		displayVersions = versions
	} else if count > 0 && count < len(versions) {
		displayVersions = versions[:count]
	}

	for _, version := range displayVersions {
		status := "installed"
		createdAt := ""
		versionName := "nightly"
		t, err := time.Parse(time.RFC3339, version.CreatedAt)
		if err == nil {
			date := t.Format("2006-01-02")
			dateMap[date]++

			// If there are multiple nightlies on the same day, show a more detailed timestamp
			if dateMap[date] > 1 {
				createdAt = t.Format("2006-01-02 15:04")
			} else {
				createdAt = date
			}
		}
		if versionName == currentVersion {
			status = "used"
		}

		table.Append([]string{versionName, createdAt, fmt.Sprint(version.UniqueNumber), status})
	}

	table.Append([]string{"", "", "", ""})
	table.Append([]string{"Total\n(nightlies)", fmt.Sprintf("%d", len(versions)), "", ""})

	table.Render()

	fmt.Println(tableString.String())
}
