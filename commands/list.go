package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type Tag struct {
	Name string `json:"name"`
}

// TODO: this is working but probably there is a better version
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
