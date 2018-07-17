package models

import (
	"encoding/gob"
	"net/http"

	"github.com/zmb3/spotify"
)

const (
	SPOTIFY_SERVICE = "spotify"
)

var spotifyPermissions = []string{
	spotify.ScopePlaylistReadPrivate,
	spotify.ScopePlaylistReadCollaborative,
}

type SpotifyAuthenticator struct {
	auth spotify.Authenticator
}

func NewSpotifyAuthenticator(addr string) *SpotifyAuthenticator {
	return &SpotifyAuthenticator{
		auth: spotify.NewAuthenticator(addr+"/callback/"+SPOTIFY_SERVICE, spotifyPermissions...),
	}
}

func (sa *SpotifyAuthenticator) BuildAuthURL(state string) string {
	return sa.auth.AuthURL(state)
}

func (sa *SpotifyAuthenticator) NewClient(state string, r *http.Request) (interface{}, error) {
	// acquire access token (also checks state parameter)
	tok, err := sa.auth.Token(state, r)
	if err != nil {
		return nil, err
	}

	client := sa.auth.NewClient(tok)
	return &spotifyClient{&client}, nil
}

func (sa *SpotifyAuthenticator) GetType() string {
	return SPOTIFY_SERVICE
}

type spotifyClient struct {
	*spotify.Client
	// token *oauth2.Token
}

func init() {
	gob.Register(&spotifyClient{})
}
