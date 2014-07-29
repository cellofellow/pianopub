package ws

import (
	"log"
	"regexp"
	"strings"

	"github.com/cellofellow/pianopub/data"
	"github.com/mattbaird/turnpike"
)

type signup struct {
	db *data.Database
}

func newsignup(db *data.Database) *signup {
	return &signup{db}
}

var nowhitespace = regexp.MustCompile(`^\S$`)

func (s *signup) HandleRPC(clientID string, topicURI string, args ...interface{}) (interface{}, error) {
	var email, name, nick, password string
	var err error

	log.Printf("rpc: %s, %s", clientID, topicURI)

	if len(args) != 4 {
		return nil, turnpike.RPCError{
			URI:         topicURI,
			Description: "Invalid Call",
			Details:     "Incorrect number of arguments. Must have 4: email, name, nick, and password.",
		}
	}

	email, err = argToString(args[0], topicURI, "Email")
	if err != nil {
		return nil, err
	}

	name, err = argToString(args[1], topicURI, "Name")
	if err != nil {
		return nil, err
	}

	nick, err = argToString(args[2], topicURI, "Nick")
	if err != nil {
		return nil, err
	}

	password, err = argToString(args[3], topicURI, "Password")
	if err != nil {
		return nil, err
	}

	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)
	nick = strings.TrimSpace(nick)
	password = strings.TrimSpace(password)

	if !strings.Contains(email, "@") {
		return nil, turnpike.RPCError{
			URI:         topicURI,
			Description: "Invalid Data",
			Details:     "Email should have an \"@\"",
		}
	}

	if nowhitespace.Match([]byte(nick)) {
		return nil, turnpike.RPCError{
			URI:         topicURI,
			Description: "Invalid Data",
			Details:     "Nick should be 1 word.",
		}
	}

	user, err := s.db.AddUser(email, name, nick, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
