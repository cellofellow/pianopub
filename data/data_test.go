package data

import "testing"
import "strings"
import "fmt"

const (
	email = "test@test.email"
	name = "Test User"
	nick = "tester"
	password = "test password"
	wrongpassword = "wrong password"
)

func TestNewDatabase(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	users, err := db.dbmap.Select(User{}, `SELECT * FROM user`)
	if err != nil {
		t.Error(err)
		return
	}
	if len(users) > 0 {
		t.Error("Shouldn't have any users, we didn't create one.")
	}
}

func TestAddUser(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}
	if !(user.Email == email && user.Name == name && user.Nick == nick) {
		t.Error("Database did not save identical data.")
	}
}

func TestAddUserUniqueConflict(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = db.AddUser("another@test.email", "New Name", "tester", wrongpassword)
	if err == nil {
		t.Error("Should have returned a Unique on User.Nick error.")
	}
}

func TestGetUser(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	userMade, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}

	userGot, err := db.GetUser(email)
	if err != nil {
		t.Error(err)
		return
	}

	if *userMade != *userGot {
		t.Error("Users are not the same.")
	}
}

func TestGetNonExistantUser(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = db.GetUser("tester2@test.email")
	if err == nil {
		t.Error("Selected a user that doesn't exist.")
	}
}

func TestAuthUser(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}
	if !(user.Authenticate(password)) {
		t.Error("Did not authenticate with correct password.")
	}
}

func TestNoAuthUser(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}
	if user.Authenticate(wrongpassword) {
		t.Error("Authenticated with incorrect password.")
	}
}

func TestUserChangePassword(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}

	err = user.ChangePassword(wrongpassword)
	if err != nil {
		t.Error(err)
		return
	}

	if !user.Authenticate(wrongpassword) {
		t.Error("Password not successfully changed.")
	}
}

func TestUserSetAdmin(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}

	err = user.SetAdmin(true)
	if err != nil {
		t.Error(err)
		return
	}

	if !user.Admin {
		t.Error("User is not admin.")
	}
}

func TestUserString(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}

	if strings.Index(fmt.Sprintln(user), "⟨User »") != 0 {
		t.Error("String of user does not begin properly.")
	}
}

func TestUserIncrementRep(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}

	err = user.IncrementRep()
	if err != nil {
		t.Error(err)
		return
	}

	if user.Rep != 1 {
		t.Error("Rep not successfully incremented.")
	}
}

func TestUserDecrementRep(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
		return
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
		return
	}

	user.IncrementRep()
	user.IncrementRep()
	user.IncrementRep()
	err = user.DecrementRep()
	if err != nil {
		t.Error(err)
		return
	}

	if user.Rep != 2 {
		t.Log(user.Rep)
		t.Error("Rep not successfully decremented.")
	}
}

func TestHashPassword(t *testing.T) {
	var uninitializedstring string
	pw := HashPassword(password)
	if pw.Salt == uninitializedstring {
		t.Error("Salt is empty string.")
	}
	if pw.Hash == uninitializedstring {
		t.Error("Password is empty string.")
	}
}

func TestCheckPassword(t *testing.T) {
	pw := HashPassword(password)
	if pw != CheckPassword(password, pw.Salt) {
		t.Error("Check password failed.")
	}
}

func TestCheckPasswordFailWrongPassword(t *testing.T) {
	pw := HashPassword(password)
	if pw == CheckPassword(wrongpassword, pw.Salt) {
		t.Error("Password hash collision.")
	}
}

func TestCheckPasswordFailWrongSalt(t *testing.T) {
	pw := HashPassword(password)
	if pw == CheckPassword(password, "fakeSalt") {
		t.Error("Password and salt hash collision.")
	}
}

func TestHashPasswordRandomSalts(t *testing.T) {
	pws := make([]HashedPassword, 100)
	for i := 0; i < 100; i++ {
		pws[i] = HashPassword(password)
	}
	var p HashedPassword
	pw := pws[0]
	for i := 1; i < 100; i++ {
		p = pws[i]
		if p == pw {
			t.Error("Salt collision.")
			return
		}
		pw = p
	}
}
