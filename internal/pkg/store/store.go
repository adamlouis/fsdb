package store

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
			id INTEGER PRIMARY KEY,
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
			modified TEXT NOT NULL,
			md5 TEXT,
			CONSTRAINT file_path_uniq UNIQUE(path)
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

	// TODO: use sqlx if this grows: https://github.com/jmoiron/sqlx
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
		modified,
		md5
	) VALUES (
		?,?,?,?,?,?,?,?,?,?,?,?,?
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
		"",
	)

	return err
}

// UpdateMD5SumArgs wraps the args for UpdateMD5Sum
type UpdateMD5SumArgs struct {
	Path string
	MD5  string
}

// UpdateMD5Sum updates the md5 value of the file
func UpdateMD5Sum(db *sql.DB, args *UpdateMD5SumArgs) error {
	if args == nil {
		return fmt.Errorf("`args` must not be nil")
	}

	_, err := db.Exec(`UPDATE file SET md5 = ? WHERE path = ?`, args.MD5, args.Path)

	return err
}

// ListPathsForUnsetMD5Result is the result of ListPathsForUnsetMD5
type ListPathsForUnsetMD5Result struct {
	Paths  []string
	Cursor int64
}

// ListPathsForUnsetMD5 returns paths whose md5 hashes are not set
// TODO: awkward api ... consider generic list fn if needed
func ListPathsForUnsetMD5(db *sql.DB, count int, cur int64) (*ListPathsForUnsetMD5Result, error) {
	rows, err := db.Query(`
		SELECT
			id, path
		FROM
			file
		WHERE
			(md5 IS NULL OR md5 = '')
			AND is_dir = false
			AND size > 0
			AND id > ?
		ORDER BY SIZE
		LIMIT ?`,
		cur, count)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	id := int64(0)
	for rows.Next() {
		path := ""
		rows.Scan(&id, &path)
		paths = append(paths, path)
	}
	return &ListPathsForUnsetMD5Result{
		Paths:  paths,
		Cursor: id,
	}, nil
}
