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
func ListHandler(cmd *cobra.Command, args []string) {
	numVersions := 7 // Default number of versions to list

	if len(args) >= 1 {
		requestedVersions, err := strconv.Atoi(args[0])
		if err == nil {
			numVersions = requestedVersions
		} else {
			fmt.Fprintln(os.Stderr, "Invalid argument. Using default number of versions.")
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

func ListHandlerLocal(cmd *cobra.Command, args []string) { // Or no args if you want to list everything
	versions, err := utils.ReadVersionsInfo()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read versions info:", err)
		return
	}

	// Determine the current version
	currentVersion, err := utils.DetermineCurrentVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to determine current version:", err)
		return
	}

	// Prepare data for the table
	var tableData []VersionData
	for _, version := range versions {
		status := "installed"
		if strings.HasPrefix(version.Directory, "nightly") {
			if version.Directory == currentVersion {
				status = "used"
			} else {
				// Display the creation date
				t, err := time.Parse(time.RFC3339, version.CreatedAt)
				if err == nil {
					status = t.Format("2006-01-02")
				}
			}
		} else if version.Directory == currentVersion { // Stable version
			status = "used"
		}

		tableData = append(tableData, VersionData{Version: version.Directory, Status: status})
	}

	// Create the table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Version", "Status"})
	// TODO: I'm not sure what I'm doing here
	// maybe table.BulkAppend(tableData) but we need it as a []string
	for _, row := range tableData {
		table.Append([]string{row.Version, row.Status})
	}
	table.Render()
}
