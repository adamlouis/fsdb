package hasher

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/adamlouis/fsdb/internal/pkg/md5"
	"github.com/adamlouis/fsdb/internal/pkg/store"
)

// Options are options for `hash`
type Options struct {
	Verbose   bool
	BatchSize int
}

// Hash writes the md5 hash of all files in the db
func Hash(db *sql.DB, opts *Options) error {
	if opts == nil {
		return fmt.Errorf("`opts` may not be nil")
	}

	bsz := 10
	if opts.BatchSize > 0 {
		bsz = opts.BatchSize
	}

	i := 0
	cur := int64(0)
	for {
		result, err := store.ListPathsForUnsetMD5(db, bsz, cur)
		if err != nil {
			return err
		}
		i++
		if len(result.Paths) == 0 {
			break
		}
		cur = result.Cursor

		var wg sync.WaitGroup
		for _, p := range result.Paths {
			path := p
			wg.Add(1)
			go func() {
				defer wg.Done()

				sum, err := md5.GetMD5(path)
				if err != nil {
					return
				}

				err = store.UpdateMD5Sum(db, &store.UpdateMD5SumArgs{
					Path: path,
					MD5:  sum,
				})

				i++
				if opts.Verbose {
					fmt.Printf("\rindex | files hashed: %d", i)
				}
			}()
		}
		wg.Wait()
	}
	fmt.Println("")
	return nil
}

func ptrS(s string) *string {
	return &s
}
func ptrB(b bool) *bool {
	return &b
}
