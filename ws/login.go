package ws

import (
	"log"

	"github.com/cellofellow/pianopub/data"
	"github.com/mattbaird/turnpike"
)

type login struct {
	clientsLoggedIn map[string]*data.User
	db              *data.Database
}

func newlogin(db *data.Database) *login {
	return &login{
		clientsLoggedIn: make(map[string]*data.User),
		db:              db,
	}
}

func (l *login) HandleRPC(clientID string, topicURI string, args ...interface{}) (interface{}, error) {
	var email, password string
	var err error

	log.Printf("rpc: %s, %s", clientID, topicURI)

	if len(args) != 2 {
		return nil, turnpike.RPCError{
			URI:         topicURI,
			Description: "Invalid Call",
			Details:     "Incorrect number of arguments. Must have 2: email and password.",
		}
	}

	email, err = argToString(args[0], topicURI, "Email")
	if err != nil {
		return nil, err
	}

	password, err = argToString(args[1], topicURI, "Password")
	if err != nil {
		return nil, err
	}

	user, err := l.db.GetUser(email)
	if err != nil {
		return nil, turnpike.RPCError{
			URI:         topicURI,
			Description: "Login Failed",
			Details:     "Invalid email or password",
		}
	}

	if user.HashedPassword != data.CheckPassword(password, user.Salt) {
		return nil, turnpike.RPCError{
			URI:         topicURI,
			Description: "Login Failed",
			Details:     "Invalid email or password",
		}
	}

	l.clientsLoggedIn[clientID] = user
	return user, nil
}
