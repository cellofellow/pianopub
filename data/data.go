package data

import (
	"database/sql"
	"log"
	"os"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)


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

	// Create table for User. Fields set with struct above.
	usertable := dbmap.AddTableWithName(User{}, "user")
	usertable.SetKeys(false, "Email")
	usertable.ColMap("Email").SetNotNull(true)
	usertable.ColMap("Name").SetNotNull(true)
	usertable.ColMap("Nick").SetNotNull(true).SetUnique(true)
	usertable.ColMap("Salt").SetNotNull(true)
	usertable.ColMap("Hash").SetNotNull(true)
	usertable.ColMap("Rep").SetNotNull(true)
	usertable.ColMap("Admin").SetNotNull(true)

	_, err = dbmap.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			Email TEXT NOT NULL PRIMARY KEY,
			Name TEXT NOT NULL,
			Nick TEXT NOT NULL UNIQUE CHECK (Nick NOT LIKE '% %'),
			Salt TEXT NOT NULL CHECK (length(Salt) = 16),
			Hash TEXT NOT NULL CHECK (length(Hash) = 44),
			Rep INTEGER NOT NULL CHECK (Rep >= 0),
			Admin INTEGER NOT NULL CHECK (Admin in (0,1))
		)
	`)
	if err != nil {
		return nil, err
	}

	return &Database{dbmap}, nil
}
