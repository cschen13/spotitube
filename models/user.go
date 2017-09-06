package models

import (
	"log"
	"net/http"
)

const (
	SPOTIFY_CLIENT = "spotify"
)

type User struct {
	state   string
	clients map[string]Client
}

func NewUser(state string, r *http.Request, auth Authenticator, clientType string) (*User, error) {
	client, err := auth.newClient(state, r)
	if err != nil {
		return nil, err
	}

	return &User{
		state:   state,
		clients: map[string]Client{clientType: client},
	}, nil
}

var (
	users   = make(map[string]*User)
	addChan = make(chan addReq)
	getChan = make(chan getReq)
)

type addReq struct {
	state string
	user  *User
}

func (user *User) Add() {
	addChan <- addReq{user.state, user}
}

type getReq struct {
	state string
	res   chan *User
}

func GetUser(state string) *User {
	res := make(chan *User)
	getChan <- getReq{state, res}
	return <-res
}

func HandleUsers() {
	for {
		select {
		case session := <-addChan:
			if _, present := users[session.state]; !present {
				log.Printf("New session with state: %s", session.state)
				users[session.state] = session.user
			} else {
				log.Printf("User with state %s already exists: skipping addition")
			}
		case req := <-getChan:
			if user, present := users[req.state]; present {
				req.res <- user
			} else {
				req.res <- nil
			}
		}
	}
}

func (user *User) GetClient(key string) Client {
	if client, present := user.clients[key]; present {
		return client
	}
	return nil
}
