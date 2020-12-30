package indexer

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adamlouis/fsdb/internal/pkg/hasher"
	"github.com/adamlouis/fsdb/internal/pkg/store"
)

// Options are options for `Index`
type Options struct {
	Root    string
	Verbose bool
	Hash    bool
}

// Index crawls the filesysten starting at `fs.Root` and writes metadata to `db`
func Index(db *sql.DB, opts *Options) error {
	if opts == nil {
		return fmt.Errorf("`opts` may not be nil")
	}
	if opts.Verbose {
		fmt.Printf("\rindex | root=%s hash=%v\n", opts.Root, opts.Hash)
	}

	err := store.CreateFileTable(db)
	if err != nil {
		return err
	}

	i := 0
	err = filepath.Walk(
		opts.Root,
		func(path string, info os.FileInfo, err error) error {
			insertErr := store.InsertFileRow(db, &store.InsertFileRowArgs{
				Path:      path,
				FileInfo:  info,
				AccessErr: err,
			})
			if opts.Verbose {
				fmt.Printf("\rindex | files indexed: %d", i)
			}
			i++
			return insertErr
		})
	if err != nil {
		return err
	}

	fmt.Println("")
	if opts.Hash {
		err = hasher.Hash(db, &hasher.Options{
			Verbose: opts.Verbose,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
