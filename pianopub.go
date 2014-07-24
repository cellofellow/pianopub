package main

import (
	"flag"
	"log"

	"github.com/cellofellow/pianopub/data"
)

func main() {
	var databaseFile string
	flag.StringVar(&databaseFile, "db", ":memory:", "database file name")
	flag.Parse()

	db, err := data.NewDatabase(databaseFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(db)
}
