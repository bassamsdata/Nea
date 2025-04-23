package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

// printDirContents prints the directory structure for debugging purposes
func printDirContents(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			rel = path
		}

		if rel == "." {
			return nil
		}

		type_indicator := " "
		if info.IsDir() {
			type_indicator = "/"
		} else if info.Mode()&os.ModeSymlink != 0 {
			type_indicator = "@"
		} else if info.Mode()&0111 != 0 {
			type_indicator = "*"
		}

		fmt.Printf("%s%s\n", rel, type_indicator)
		return nil
	})
}
