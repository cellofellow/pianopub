package ws

import (
	"log"
    "time"

	"github.com/cellofellow/pianopub/data"
	"gopkg.in/jcelliott/turnpike.v1"
)

type clientUser struct {
	ClientId string
	User     *data.User
}

type login struct {
	*data.Database
    server          *turnpike.Server
	clientsLoggedIn map[string]*data.User
	ch              chan clientUser
}

func newLogin(db *data.Database, server *turnpike.Server) *login {
	l := &login{
		db,
        server,
		make(map[string]*data.User),
		make(chan clientUser, 1000),
	}
	l.startManager()
	return l
}

func (l *login) startManager() {
    go func () {
        for cu := range l.ch {
            l.clientsLoggedIn[cu.ClientId] = cu.User
        }
    }()

    // Log out clients if no longer connected.
    go func() {
        time.Sleep(time.Second * 10)
        clientIDs := make(map[string]struct{})
        for _, clientID := range l.server.ConnectedClients() {
            clientIDs[clientID] = struct{}{}
        }
        for clientID, _ := range l.clientsLoggedIn {
            if _, ok := clientIDs[clientID]; !ok {
                delete(l.clientsLoggedIn, clientID)
            }
        }
    }()
}

func (l *login) addClient(cu clientUser) {
	l.ch <- cu
}

func (l *login) GetClient(clientID string) (user *data.User, ok bool) {
	user, ok = l.clientsLoggedIn[clientID]
	return
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

	user, err := l.GetUser(email)
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

	l.addClient(clientUser{clientID, user})
	return user, nil
}
