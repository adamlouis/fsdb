package fsdb

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/adamlouis/fsdb/internal/pkg/indexer"
	_ "github.com/mattn/go-sqlite3" // import sql driver
	"github.com/spf13/cobra"
)

var (
	indexFlagOutput  string
	indexFlagRoot    string
	indexFlagVerbose bool
	indexFlagHash    bool
)

var indexCommand = &cobra.Command{
	Use:   "index",
	Short: "Index the file system",
	RunE: func(cmd *cobra.Command, args []string) error {
		if indexFlagVerbose {
			fmt.Printf("fsdb | start=%v\n", time.Now())
		}

		// set db
		out := indexFlagOutput
		if out == "" {
			out = fmt.Sprintf("./fsdb.%d.db", time.Now().Unix())
		}

		db, err := sql.Open("sqlite3", out)
		if err != nil {
			return err
		}

		// set root
		root := indexFlagRoot
		if root == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			root = home
		}

		// run the indexer
		err = indexer.Index(db, &indexer.Options{
			Root:    root,
			Verbose: indexFlagVerbose,
			Hash:    indexFlagHash,
		})
		if indexFlagVerbose {
			fmt.Printf("fsdb | end=%v\n", time.Now())
		}
		return err
	},
}

func init() {
	indexCommand.Flags().StringVarP(&indexFlagOutput, "output", "o", "", "the destination file to write to")
	indexCommand.Flags().StringVarP(&indexFlagRoot, "root", "r", "", "the root directory to index")
	indexCommand.Flags().BoolVarP(&indexFlagVerbose, "verbose", "v", false, "print verbose execution logs")
	indexCommand.Flags().BoolVar(&indexFlagHash, "hash", false, "after indexing, fill in the hash of each file")

	root.AddCommand(indexCommand)
}
