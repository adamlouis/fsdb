package fsdb

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/adamlouis/fsdb/internal/pkg/hasher"
	_ "github.com/mattn/go-sqlite3" // import sql driver
	"github.com/spf13/cobra"
)

var (
	hashFlagInput   string
	hashFlagVerbose bool
)

var hashCommand = &cobra.Command{
	Use:   "hash",
	Short: "Fill in the md5 hash of files in the db",
	RunE: func(cmd *cobra.Command, args []string) error {
		if hashFlagVerbose {
			fmt.Printf("fsdb | start=%v\n", time.Now())
		}

		if !exists(hashFlagInput) {
			return fmt.Errorf("file `%s` does not exist", hashFlagInput)
		}

		db, err := sql.Open("sqlite3", hashFlagInput)
		if err != nil {
			return err
		}

		err = hasher.Hash(db, &hasher.Options{
			Verbose: true,
		})

		if hashFlagVerbose {
			fmt.Printf("fsdb | end=%v\n", time.Now())
		}
		return err
	},
}

func init() {
	hashCommand.Flags().StringVarP(&hashFlagInput, "input", "i", "", "the db file to read from")
	hashCommand.Flags().BoolVarP(&hashFlagVerbose, "verbose", "v", false, "print verbose execution logs")

	root.AddCommand(hashCommand)
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
