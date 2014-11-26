package web

import (
	"github.com/gocraft/web"

	"github.com/cellofellow/pianopub/data"
)

func Router(db *data.Database) *web.Router {
	router := web.New(Context{})
	router.Middleware(web.LoggerMiddleware)
	router.Middleware(web.StaticMiddleware("static"))
	router.Middleware(PersistantMiddleware(db))
	router.Post("/user", (*Context).Signup)
	return router
}
