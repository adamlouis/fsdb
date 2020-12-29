package index

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adamlouis/fsdb/internal/dbw"
)

// Options are options for `Index`
type Options struct {
	Root    string
	Verbose bool
}

// Index crawls the filesysten starting at `fs.Root` and writes metadata to `db`
func Index(db *sql.DB, opts *Options) error {
	if opts == nil {
		return fmt.Errorf("`opts` may not be nil")
	}
	if opts.Verbose {
		fmt.Printf("\rindex | root=%s\n", opts.Root)
	}

	err := dbw.CreateFileTable(db)
	if err != nil {
		return err
	}

	i := 0
	return filepath.Walk(
		opts.Root,
		func(path string, info os.FileInfo, err error) error {
			insertErr := dbw.InsertFileRow(db, &dbw.InsertFileRowArgs{
				Path:      path,
				FileInfo:  info,
				AccessErr: err,
			})
			if opts.Verbose {
				fmt.Printf("\rindex | files processed: %d", i)
				i++
			}
			return insertErr
		})
}
