package models

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	youtube "google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"strings"
)

const YOUTUBE_SERVICE = "youtube"

var youtubePermissions = youtube.YoutubeScope

type YoutubeAuthenticator struct {
	config  *oauth2.Config
	context context.Context
}

func NewYoutubeAuthenticator(json string, addr string, isDev bool) *YoutubeAuthenticator {
	config, err := google.ConfigFromJSON([]byte(json), youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	if isDev {
		config.RedirectURL = addr + "/callback/" + YOUTUBE_SERVICE
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

func (client *youtubeClient) GetPlaylistInfo(playlistId string) (Playlist, error) {
	return nil, errors.New("Unimplemented")
}

func (client *youtubeClient) CreatePlaylist(name string) (Playlist, error) {
	resource := make(map[string]interface{})
	resource["snippet"] = make(map[string]interface{})
	snippet := resource["snippet"].(map[string]interface{})
	snippet["title"] = name
	jsonObj, err := json.Marshal(resource)
	if err != nil {
		log.Printf("youtube: Failed to encode JSON for playlist resource")
		return nil, err
	}

	playlist := &youtube.Playlist{}
	if err := json.NewDecoder(strings.NewReader(string(jsonObj))).Decode(&playlist); err != nil {
		log.Printf("youtube: Failed to decode JSON into playlist resource")
		return nil, err
	}

	call := client.Playlists.Insert("id,snippet", playlist)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	return &youtubePlaylist{response}, nil
}

func (client *youtubeClient) GetPlaylistTracks(playlist Playlist, page string) ([]PlaylistTrack, bool, error) {
	return nil, true, errors.New("Unimplemented")
}

func (client *youtubeClient) InsertTrack(playlist Playlist, track PlaylistTrack) error {
	return errors.New("Unimplemented")
}

type youtubePlaylist struct {
	obj *youtube.Playlist
}

func (playlist *youtubePlaylist) GetID() string {
	return playlist.obj.Id
}

func (playlist *youtubePlaylist) GetName() string {
	return playlist.obj.Snippet.Title
}

func (playlist *youtubePlaylist) GetURL() string {
	return "https://www.youtube.com/playlist?list=" + playlist.GetID()
}

func (playlist *youtubePlaylist) GetCoverURL() string {
	if thumbnails := playlist.obj.Snippet.Thumbnails; thumbnails != nil {
		return thumbnails.Default.Url
	}
	return ""
}
