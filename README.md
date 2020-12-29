# fsdb

## what

`fsdb` is a small utility for indexing & exploring file metadata with sql.

`fsdb` traverses the file system & writes data to a local `sqlite3` db file.

`sqlite3` schema for the file table is:

```
CREATE TABLE file(
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
);
```

## how

See the `Makefile` for commands

Run the index command with: `go run main.go index -v`

## why

I wrote this utility to help manage disk space & disk organization on my personal computer. I needed a tool to get organized, free up space, and locate important documents.

The `sqlite3` backend makes it easy to answer ad-hoc questions about filesystem contents after initial indexing.

## examples

10 largest files:

`SELECT path, size FROM file ORDER BY size DESC LIMIT 10;`

10 file extensions that consume most space:

`SELECT ext, SUM(size)/1e9 AS GB FROM file GROUP BY ext ORDER BY GB DESC LIMIT 10;`

## todo

here are some things I may considering doing:

- specific file support - start w/ images ... write filetype-specific columns (or tables) with additional metadata attributes (e.g., image dimensions, general color, camera, etc.)

- content id - store content hash (md5/sha256/etc) in table to dedupe files

- `miner` - publish & consider merging code. I have another small utility "`miner`" to help with photo/file organization & depuplication... considering merging or make `miner` and `fsdb` work nicely together.

- ancestors - add table to associate file with all parent path prefixes of a file to enable efficient queries of dir size & more

- time estimate - show estimated time to complete indexing when running with `--verbose`
