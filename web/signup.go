package web

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/gocraft/web"
)

type signupRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Nick     string `json:"nick"`
	Password string `json:"password"`
}

var nowhitespace = regexp.MustCompile(`^\S$`)

func (c *Context) Signup(w web.ResponseWriter, r *web.Request) {
	var err error
	decoder := json.NewDecoder(r.Body)
	s := &signupRequest{}
	err = decoder.Decode(s)
	r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	s.Email = strings.TrimSpace(s.Email)
	s.Name = strings.TrimSpace(s.Name)
	s.Nick = strings.TrimSpace(s.Nick)
	s.Password = strings.TrimSpace(s.Password)

	required := make([]string, 0, 4)
	if s.Email == "" {
		required = append(required, "email")
	}
	if s.Name == "" {
		required = append(required, "name")
	}
	if s.Nick == "" {
		required = append(required, "nick")
	}
	if s.Password == "" {
		required = append(required, "password")
	}
	if len(required) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(strings.Join(required, ",") + " are required fields"))
		return
	}

	addr, err := mail.ParseAddress(s.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email error: " + err.Error()))
		return
	}

	s.Email = addr.Address

	if nowhitespace.Match([]byte(s.Nick)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Nick should be 1 word."))
		return
	}

	_, err = c.db.AddUser(s.Email, s.Name, s.Nick, s.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
