package ws

import (
	"gopkg.in/jcelliott/turnpike.v1"

	"github.com/cellofellow/pianopub/data"
)

var server *turnpike.Server

func Server(db *data.Database) *turnpike.Server {
	s := turnpike.NewServer()
	server = s

	l := newLogin(db, s)
	s.RegisterRPC("rpc:login", l.HandleRPC)
	return s

}

func argToString(arg interface{}, topicURI, name string) (string, error) {
	if s, ok := arg.(string); !ok {
		return "", turnpike.RPCError{
			URI:         topicURI,
			Description: "Invalid Call",
			Details:     name + " must be a string.",
		}
	} else {
		return s, nil
	}
	panic("unreachable")
}
