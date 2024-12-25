package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nvm_manager_go/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
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
	Use:   "ls",
	Short: "List all Neovim versions",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: You must specify local or remote")
			return
		}
		switch args[0] {
		case "local":
			listHandlerLocal()
		case "remote":
			listHandler(args[1:])
		}
	},
}

// TODO: we can use the fetchreleases() function and then call a table
func listHandler(args []string) {
	numVersions := 7 // Default number of versions to list

	if len(args) >= 1 {
		requestedVersions, err := strconv.Atoi(args[0])
		if err == nil {
			numVersions = requestedVersions
		} else {
			fmt.Fprintln(os.Stderr, "no argument. Using default number of versions which is 7")
		}
	}
	resp, err := http.Get("https://api.github.com/repos/neovim/neovim/tags")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch Neovim versions:", err) // Changed fmt.Println to fmt.Fprintln
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

	for i, tag := range tags {
		if i >= numVersions {
			break
		}
		fmt.Println(tag.Name)
	}
}

func listHandlerLocal() {
	versions, err := utils.ReadVersionsInfo()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read versions info:", err)
		return
	}

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
		for _, version := range stableVersions {
			versionName := version.Name()
			status := "stable"
			if versionName == currentVersion {
				status = "used"
			}
			table.Append([]string{versionName, "", "N/A", status})
		}
	}

	for _, version := range versions {
		status := "installed"
		createdAt := ""
		versionName := "nightly"
		t, err := time.Parse(time.RFC3339, version.CreatedAt)
		if err == nil {
			createdAt = t.Format("2006-01-02")
		}
		if strings.Contains(version.Directory, currentVersion) {
			status = "used"
		}

		table.Append([]string{versionName, createdAt, fmt.Sprint(version.UniqueNumber), status})
	}

	table.Render()

	fmt.Println(tableString.String())
}
