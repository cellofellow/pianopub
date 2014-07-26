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
	return s
}
