package fsdb

import (
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "fsdb",
	Short: "fsdb is a utility for indexing & exploring file metadata with sql.",
}

// Execute executes the root fsdb command
func Execute() error {
	return root.Execute()
}
