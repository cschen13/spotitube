package models

import (
	"log"
	"net/http"
)

type User struct {
	state   string
	spotify *spotifyClient
}

func NewUser(state string, r *http.Request, auth *SpotifyAuthenticator) (*User, error) {
	client, err := auth.newSpotifyClient(state, r)
	if err != nil {
		return nil, err
	}

	return &User{state: state, spotify: client}, nil
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

func (user *User) GetSpotifyPlaylists(pageNumber int) ([]SpotifyPlaylist, error) {
	return user.spotify.getPlaylists(pageNumber)
}
