package dbw

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CreateFileTable creates the file table
func CreateFileTable(db *sql.DB) error {
	create := `CREATE TABLE file(
			name TEXT NOT NULL,
			path TEXT NOT NULL,
			ext TEXT NOT NULL,
			size INTEGER,
			is_dir INTEGER,
			error INTEGER,
			error_message TEXT,
			is_notexist INTEGER,
			is_permission INTEGER,
			is_timeout INTEGER,
			mode TEXT NOT NULL,
			modified TEXT NOT NULL
		);`

	_, err := db.Exec(create)
	return err
}

// InsertFileRowArgs wraps the args for InsertFileRow
type InsertFileRowArgs struct {
	Path      string
	FileInfo  os.FileInfo
	AccessErr error
}

// InsertFileRow inserts a row into the file table
func InsertFileRow(db *sql.DB, args *InsertFileRowArgs) error {
	if args == nil {
		return fmt.Errorf("`args` must not be nil")
	}

	errorMessage := ""
	if args.AccessErr != nil {
		errorMessage = args.AccessErr.Error()
	}

	_, err := db.Exec(`
	INSERT INTO
		file (
		name,
		path,
		ext,
		size,
		is_dir,
		error,
		error_message,
		is_notexist,
		is_permission,
		is_timeout,
		mode,
		modified
	) VALUES (
		?,?,?,?,?,?,?,?,?,?,?,?
	)`,
		args.FileInfo.Name(),
		args.Path,
		filepath.Ext(args.Path),
		args.FileInfo.Size(),
		args.FileInfo.IsDir(),
		args.AccessErr != nil,
		errorMessage,
		os.IsNotExist(args.AccessErr),
		os.IsPermission(args.AccessErr),
		os.IsTimeout(args.AccessErr),
		args.FileInfo.Mode().String(),
		args.FileInfo.ModTime().Format(time.RFC3339),
	)

	return err
}
