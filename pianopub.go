package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/cellofellow/pianopub/data"
	"github.com/cellofellow/pianopub/ws"
)

func main() {
	var databaseFile, listen string
	flag.StringVar(&databaseFile, "db", ":memory:", "database file name")
	flag.StringVar(&listen, "listen", ":8080", "host:port to listen on")
	flag.Parse()

	db, err := data.NewDatabase(databaseFile)
	if err != nil {
		log.Fatal(err)
	}

	if !db.AddFirstUser() {
		log.Fatal("No admin user added.")
	}

	s := ws.Server(db)
	http.Handle("/ws", s.Handler)
	http.Handle("/", http.FileServer(http.Dir("http")))
	log.Println("Listening at", listen)
	err = http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}
