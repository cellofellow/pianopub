package data

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/coopernurse/gorp"
)

type User struct {
	Email string      `json:"email"`
	Name  string      `json:"name"`
	Nick  string      `json:"nick"`
	HashedPassword
	Rep   uint32      `json:"rep"`
	Admin bool        `json:"admin"`
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

func (u User) ddl(dbmap *gorp.DbMap) error {
	table := dbmap.AddTableWithName(User{}, "user")
	table.SetKeys(false, "Email")
	table.ColMap("Email").SetNotNull(true)
	table.ColMap("Name").SetNotNull(true)
	table.ColMap("Nick").SetNotNull(true).SetUnique(true)
	table.ColMap("Salt").SetNotNull(true)
	table.ColMap("Hash").SetNotNull(true)
	table.ColMap("Rep").SetNotNull(true)
	table.ColMap("Admin").SetNotNull(true)

	_, err := dbmap.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			Email TEXT NOT NULL PRIMARY KEY,
			Name  TEXT NOT NULL,
			Nick  TEXT NOT NULL UNIQUE CHECK (Nick NOT LIKE '% %'),
			Salt  TEXT NOT NULL CHECK (length(Salt) = 16),
			Hash  TEXT NOT NULL CHECK (length(Hash) = 44),
			Rep   INTEGER NOT NULL CHECK (Rep >= 0),
			Admin INTEGER NOT NULL CHECK (Admin in (0,1))
		)
	`)

	return err
}

func (u *User) Update() (err error) {
	_, err = u.dbmap.Update(u)
	return
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

func (u *User) SetAdmin(admin bool) error {
	u.Admin = admin
	_, err := u.dbmap.Update(u)
	return err
}

func (d *Database) AddUser(email, name, nick, password string) (*User, error) {
	pw := HashPassword(password)
	user := &User{
		Email:          email,
		Name:           name,
		Nick:           nick,
		HashedPassword: pw,
		Rep:            0,
		Admin:          false,
	}
	err := d.dbmap.Insert(user)
	if err != nil {
		return nil, err
	}
	user.dbmap = d.dbmap
	return user, nil
}

var nowhitespace = regexp.MustCompile(`^\S$`)

func (d *Database) AddFirstUser() bool {
	var email, name, nick, password string
	var user *User

	exists, err := d.dbmap.SelectInt(`SELECT EXISTS (SELECT email FROM user WHERE admin = 1)`)
	if exists == 1 || err != nil {
		return true
	}

	fmt.Println("There are no users. At least one admin user is required.")
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Email Address: ")
	email, err = reader.ReadString('\n')
	if err != nil {
		goto fail
	}
	email = strings.TrimSpace(email)
	if !strings.Contains(email, "a") {
		log.Println("Email should have an \"@\"")
		return false
	}

	fmt.Printf("Name: ")
	name, err = reader.ReadString('\n')
	if err != nil {
		goto fail
	}
	name = strings.TrimSpace(name)

	fmt.Printf("Nick: ")
	nick, err = reader.ReadString('\n')
	if err != nil {
		goto fail
	}
	nick = strings.TrimSpace(nick)
	if nowhitespace.Match([]byte(nick)) {
		log.Println("Nick should be 1 word.")
		return false
	}

	fmt.Printf("Password: ")
	password, err = reader.ReadString('\n')
	if err != nil {
		goto fail
	}
	password = strings.TrimSpace(password)

	user, err = d.AddUser(email, name, nick, password)
	if err != nil {
		goto fail
	}
	err = user.SetAdmin(true)
	if err != nil {
		goto fail
	}

	return true

	fail:
		log.Println(err)
		return false
}

func (d *Database) GetUser(email string) (*User, error) {
	var user *User
	u, err := d.dbmap.Get(user, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("No user with email " + email)
	}
	user = u.(*User)
	user.dbmap = d.dbmap
	return user, nil
}
