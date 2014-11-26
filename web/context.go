package web

import (
	"github.com/gocraft/web"

	"github.com/cellofellow/pianopub/data"
)

type Context struct {
	db *data.Database
}

func PersistantMiddleware(db *data.Database) func(*Context, web.ResponseWriter, *web.Request, web.NextMiddlewareFunc) {
	return func(c *Context, w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		c.db = db
		next(w, r)
	}
}
