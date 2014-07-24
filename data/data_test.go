package data

import "testing"

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
	}
	users, err := db.dbmap.Select(User{}, `SELECT * FROM user`)
	if err != nil {
		t.Error(err)
	}
	if len(users) > 0 {
		t.Error("Shouldn't have any users, we didn't create one.")
	}
}

func TestAddUser(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
	}
	if !(user.Email == email && user.Name == name && user.Nick == nick) {
		t.Error("Database did not save identical data.")
	}
}

func TestAuthUser(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
	}
	if !(user.Authenticate(password)) {
		t.Error("Did not authenticate with correct password.")
	}
}

func TestNoAuthUser(t *testing.T) {
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Error(err)
	}
	user, err := db.AddUser(email, name, nick, password)
	if err != nil {
		t.Error(err)
	}
	if user.Authenticate(wrongpassword) {
		t.Error("Authenticated with incorrect password.")
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
