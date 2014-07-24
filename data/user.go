package data

import (
	"fmt"

	"github.com/coopernurse/gorp"
)

type User struct {
	Email string
	Name  string
	Nick  string
	HashedPassword
	Rep   uint32
	Admin bool
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

func (d *Database) GetUser(email string) (*User, error) {
	var user User
	err := d.dbmap.SelectOne(&user, `SELECT * FROM user WHERE email = ?`, email)
	if err != nil {
		return nil, err
	}
	user.dbmap = d.dbmap
	return &user, nil
}
