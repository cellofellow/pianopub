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

func newlogin (db *data.Database) *login {
	return &login{
		clientsLoggedIn: make(map[string]*data.User),
		db:              db,
	}
}

func (l *login) HandleRPC(clientID string, topicURI string, args ...interface{}) (interface{}, error) {
	var email, password string
	var err error
	var ok bool

	log.Printf("rpc: %s, %s, %v", clientID, topicURI)

	if len(args) != 2 {
		return nil, turnpike.RPCError{
			URI: topicURI,
			Description: "Invalid Call",
			Details: "Incorrect number of arguments. Must have 2: username and password.",
		}
	}
	if email, ok = args[0].(string); !ok {
		return nil, turnpike.RPCError{
			URI: topicURI,
			Description: "Invalid Call",
			Details: "Email must be a string.",
		}

	}
	if password, ok = args[1].(string); !ok {
		return nil, turnpike.RPCError{
			URI: topicURI,
			Description: "Invalid Call",
			Details: "Password must be a string.",
		}
	}

	user, err := l.db.GetUser(email)
	if err != nil {
		return nil, turnpike.RPCError{
			URI: topicURI,
			Description: "Login Failed",
			Details: "Invalid email or password",
		}
	}

	if user.HashedPassword != data.CheckPassword(password, user.Salt) {
		return nil, turnpike.RPCError{
			URI: topicURI,
			Description: "Login Failed",
			Details: "Invalid email or password",
		}
	}

	l.clientsLoggedIn[clientID] = user
	return user, nil
}
