package models

import (
	"log"
	"net/http"
)

type User struct {
	state   string
	clients map[string]Client
}

func NewUser(state string, r *http.Request, auth Authenticator) (*User, error) {
	user := &User{
		state:   state,
		clients: make(map[string]Client),
	}

	err := user.AddClient(r, auth)
	if err != nil {
		log.Printf("Error adding new user with state %s: %s", state, err.Error())
		return nil, err
	}
	return user, nil
}

func (user *User) GetState() string {
	return user.state
}

func (user *User) AddClient(r *http.Request, auth Authenticator) error {
	client, err := auth.newClient(user.state, r)
	if err != nil {
		return err
	}

	user.clients[auth.GetType()] = client
	return nil
}

func (user *User) GetClient(key string) Client {
	if client, present := user.clients[key]; present {
		return client
	}
	return nil
}

// TODO: DB
var (
	users      = make(map[string]*User)
	addChan    = make(chan addReq)
	getChan    = make(chan getReq)
	deleteChan = make(chan string)
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

func DeleteUser(state string) {
	deleteChan <- state
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
		case state := <-deleteChan:
			delete(users, state)
		}
	}
}
