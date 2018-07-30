package models

import (
	"net/http"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
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

func (sa *SpotifyAuthenticator) GenerateToken(state string, r *http.Request) (*oauth2.Token, error) {
	return sa.auth.Token(state, r)
}

func (sa *SpotifyAuthenticator) NewClient(tok *oauth2.Token) (interface{}, error) {
	client := sa.auth.NewClient(tok)
	return &spotifyClient{&client}, nil
}

func (sa *SpotifyAuthenticator) GetType() string {
	return SPOTIFY_SERVICE
}

type spotifyClient struct {
	*spotify.Client
}
