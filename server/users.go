package server

import (
	"github.com/zmb3/spotify"
	"log"
)

var (
	spotifyClients = make(map[string]*spotify.Client)
	addSpotifyChan = make(chan spotifySession)
	getSpotifyChan = make(chan getSpotifyReq)
)

type spotifySession struct {
	state  string
	client *spotify.Client
}

type getSpotifyReq struct {
	state string
	res   chan *spotify.Client
}

func handleUsers() {
	for {
		select {
		case session := <-addSpotifyChan:
			if _, present := spotifyClients[session.state]; !present {
				log.Printf("New session with state: %s", session.state)
				spotifyClients[session.state] = session.client
			}
		case req := <-getSpotifyChan:
			if client, present := spotifyClients[req.state]; present {
				req.res <- client
			} else {
				req.res <- nil
			}
		}
	}
}

func getSpotify(state string) *spotify.Client {
	response := make(chan *spotify.Client)
	getSpotifyChan <- getSpotifyReq{state: state, res: response}
	return <-response
}
