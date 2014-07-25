package data

import (
	"database/sql"
	"log"
	"os"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

type ddler interface {
	ddl (*gorp.DbMap) error
}

type Database struct {
	dbmap *gorp.DbMap
}

func NewDatabase(file string) (*Database, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.TraceOn("[gorp]", log.New(os.Stderr, "pianopub:", log.Lmicroseconds))
	for _, table := range []ddler{User{}, Config{}} {
		err = table.ddl(dbmap)
		if err != nil {
			return nil, err
		}
	}

	return &Database{dbmap}, nil
}
