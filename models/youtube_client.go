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
	"strconv"
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

func (ya *YoutubeAuthenticator) newClient(state string, r *http.Request) (interface{}, error) {
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

func addPropertyToResource(ref map[string]interface{}, keys []string, value string, count int) map[string]interface{} {
	for k := count; k < (len(keys) - 1); k++ {
		switch val := ref[keys[k]].(type) {
		case map[string]interface{}:
			ref[keys[k]] = addPropertyToResource(val, keys, value, (k + 1))
		case nil:
			next := make(map[string]interface{})
			ref[keys[k]] = addPropertyToResource(next, keys, value, (k + 1))
		}
	}
	// Only include properties that have values.
	if count == len(keys)-1 && value != "" {
		valueKey := keys[len(keys)-1]
		if valueKey[len(valueKey)-2:] == "[]" {
			ref[valueKey[0:len(valueKey)-2]] = strings.Split(value, ",")
		} else if len(valueKey) > 4 && valueKey[len(valueKey)-4:] == "|int" {
			ref[valueKey[0:len(valueKey)-4]], _ = strconv.Atoi(value)
		} else if value == "true" {
			ref[valueKey] = true
		} else if value == "false" {
			ref[valueKey] = false
		} else {
			ref[valueKey] = value
		}
	}
	return ref
}

func createResource(properties map[string]string) string {
	resource := make(map[string]interface{})
	for key, value := range properties {
		keys := strings.Split(key, ".")
		ref := addPropertyToResource(resource, keys, value, 0)
		resource = ref
	}
	propJson, err := json.Marshal(resource)
	if err != nil {
		log.Fatal("cannot encode to JSON ", err)
	}
	return string(propJson)
}
