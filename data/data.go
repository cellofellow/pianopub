package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Email string
	Name  string
	Nick  string
	HashedPassword
	Rep   uint32
	dbmap *gorp.DbMap `db:"-"`
}

func (u *User) String() string {
	return fmt.Sprintf(
		"⟨User » Email: %s, Name: \"%s\", Nick: %s, Rep: %d⟩",
		u.Email,
		u.Name,
		u.Nick,
		u.Rep,
	)
}

func (u *User) Authenticate(password string) bool {
	return u.HashedPassword == CheckPassword(password, u.Salt)
}

func (u *User) IncrementRep() error {
	_, err := u.dbmap.Exec(`UPDATE user SET rep = rep + 1 WHERE email = ?`, u.Email)
	if err == nil {
		u.Rep++
	}
	return err
}

func (u *User) DecrementRep() error {
	_, err := u.dbmap.Exec(`UPDATE user SET rep = rep - 1 WHERE email = ?`, u.Email)
	if err == nil {
		u.Rep--
	}
	return err
}

func (u *User) ChangePassword(password string) error {
	pw := HashPassword(password)
	u.HashedPassword = pw
	_, err := u.dbmap.Update(u)
	return err
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

	// Create table for User. Fields set with struct above.
	usertable := dbmap.AddTableWithName(User{}, "user")
	usertable.SetKeys(false, "Email")
	usertable.ColMap("Email").SetNotNull(true)
	usertable.ColMap("Name").SetNotNull(true)
	usertable.ColMap("Nick").SetNotNull(true).SetUnique(true)
	usertable.ColMap("Salt").SetNotNull(true)
	usertable.ColMap("Hash").SetNotNull(true)
	usertable.ColMap("Rep").SetNotNull(true)

	_, err = dbmap.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			Email TEXT NOT NULL PRIMARY KEY,
			Name TEXT NOT NULL,
			Nick TEXT NOT NULL UNIQUE,
			Salt BLOB NOT NULL,
			Hash BLOB NOT NULL,
			Rep INTEGER NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &Database{dbmap}, nil
}

func (d *Database) AddUser(email, name, nick, password string) (*User, error) {
	pw := HashPassword(password)
	user := &User{
		Email:          email,
		Name:           name,
		Nick:           nick,
		HashedPassword: pw,
		Rep:            0,
	}
	err := d.dbmap.Insert(user)
	if err != nil {
		return nil, err
	}
	user.dbmap = d.dbmap
	return user, nil
}

func (d *Database) GetUser(email string) (*User, error) {
	var user User
	err := d.dbmap.SelectOne(&user, `SELECT * FROM user WHERE email = ?`, email)
	if err != nil {
		return nil, err
	}
	user.dbmap = d.dbmap
	return &user, nil
}
