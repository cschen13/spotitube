package models

import (
	"encoding/gob"
	"net/http"

	"golang.org/x/oauth2"
)

// type User struct {
// 	state  string
// 	tokens map[string]*oauth2.Token
// }

type Authenticator interface {
	BuildAuthURL(string) string
	GetType() string
	GenerateToken(string, *http.Request) (*oauth2.Token, error)
	NewClient(*oauth2.Token) (interface{}, error)
}

// func NewUser(state string, r *http.Request, auth Authenticator) (*User, error) {
// 	user := &User{
// 		state:  state,
// 		tokens: make(map[string]*oauth2.Token),
// 	}

// 	err := user.AddToken(r, auth)
// 	if err != nil {
// 		log.Printf("Error adding new user with state %s: %s", state, err.Error())
// 		return nil, err
// 	}
// 	return user, nil
// }

// func (user *User) GetState() string {
// 	return user.state
// }

// func (user *User) AddToken(r *http.Request, auth Authenticator) error {
// 	tok, err := auth.GenerateToken(user.state, r)
// 	if err != nil {
// 		return err
// 	}

// 	user.tokens[auth.GetType()] = tok
// 	return nil
// }

// func (user *User) GetToken(key string) *oauth2.Token {
// 	if token, present := user.tokens[key]; present {
// 		return token
// 	}
// 	return nil
// }

// // TODO: DB
// var (
// 	users      = make(map[string]*User)
// 	addChan    = make(chan addReq)
// 	getChan    = make(chan getReq)
// 	deleteChan = make(chan string)
// )

// type addReq struct {
// 	state string
// 	user  *User
// }

// func (user *User) Add() {
// 	addChan <- addReq{user.state, user}
// }

// type getReq struct {
// 	state string
// 	res   chan *User
// }

// func GetUser(state string) *User {
// 	res := make(chan *User)
// 	getChan <- getReq{state, res}
// 	return <-res
// }

// func DeleteUser(state string) {
// 	deleteChan <- state
// }

// func HandleUsers() {
// 	for {
// 		select {
// 		case session := <-addChan:
// 			if _, present := users[session.state]; !present {
// 				log.Printf("New session with state: %s", session.state)
// 				users[session.state] = session.user
// 			} else {
// 				log.Printf("User with state %s already exists: skipping addition")
// 			}
// 		case req := <-getChan:
// 			if user, present := users[req.state]; present {
// 				req.res <- user
// 			} else {
// 				req.res <- nil
// 			}
// 		case state := <-deleteChan:
// 			delete(users, state)
// 		}
// 	}
// }

func init() {
	gob.Register(&oauth2.Token{})
}
