package ws

import (
	"github.com/cellofellow/pianopub/data"
	"github.com/mattbaird/turnpike"
)

var server *turnpike.Server

func Server(db *data.Database) *turnpike.Server {
	s := turnpike.NewServer(false)
	server = s

	l := newlogin(db)
	s.RegisterRPCFunc("rpc:login", l)
	
	sp := newsignup(db)
	s.RegisterRPCFunc("rpc:signup", sp)
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
