package models

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	youtube "google.golang.org/api/youtube/v3"
	"log"
	"net/http"
)

const YOUTUBE_SERVICE = "youtube"

var youtubePermissions = youtube.YoutubeScope

type YoutubeAuthenticator struct {
	config  *oauth2.Config
	context context.Context
}

func NewYoutubeAuthenticator(json string) *YoutubeAuthenticator {
	config, err := google.ConfigFromJSON([]byte(json), youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return &YoutubeAuthenticator{
		config:  config,
		context: context.Background(),
	}
}

func (ya *YoutubeAuthenticator) BuildAuthURL(state string) string {
	return ya.config.AuthCodeURL(state)
}

func (ya *YoutubeAuthenticator) newClient(state string, r *http.Request) (Client, error) {
	token, err := ya.token(state, r)
	if err != nil {
		return nil, err
	}

	client := ya.config.Client(ya.context, token)
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	call := service.Channels.List("snippet")
	call.Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	func(response *youtube.ChannelListResponse) {
		for _, item := range response.Items {
			log.Println(item.Id, ": ", item.Snippet.Title)
		}
	}(response)

	return &youtubeClient{
		service,
		token,
	}, nil
}

// Token pulls an authorization code from an HTTP request and attempts to exchange
// it for an access token.
func (ya *YoutubeAuthenticator) token(state string, r *http.Request) (*oauth2.Token, error) {
	values := r.URL.Query()
	if e := values.Get("error"); e != "" {
		return nil, errors.New("youtube: auth failed - " + e)
	}
	code := values.Get("code")
	if code == "" {
		return nil, errors.New("youtube: didn't get access code")
	}
	actualState := values.Get("state")
	if actualState != state {
		return nil, errors.New("youtube: redirect state parameter doesn't match")
	}
	return ya.config.Exchange(ya.context, code)
}

func (ya *YoutubeAuthenticator) GetType() string {
	return YOUTUBE_SERVICE
}

type youtubeClient struct {
	*youtube.Service
	token *oauth2.Token
}

func (client *youtubeClient) GetPlaylists(pageToken string) (playlistPage *PlaylistsPage, err error) {
	return nil, errors.New("youtube: GetPlaylists is unimplemented.")
}
