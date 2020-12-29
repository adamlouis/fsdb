package main

import (
	"log"

	"github.com/adamlouis/fsdb/cmd/fsdb"
)

func main() {
	if err := fsdb.Execute(); err != nil {
		log.Fatal(err)
	}
}
