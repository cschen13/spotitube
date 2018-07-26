package models

import (
	"github.com/zmb3/spotify"
	"log"
)

const SPOTIFY_PLAYLISTS_PAGE_LIMIT = 50

type spotifyPlaylist struct {
	obj *spotify.SimplePlaylist
}

func (playlist *spotifyPlaylist) GetID() string {
	return playlist.obj.ID.String()
}

func (playlist *spotifyPlaylist) GetOwnerID() string {
	return playlist.obj.Owner.ID
}

func (playlist *spotifyPlaylist) GetName() string {
	return playlist.obj.Name
}

func (playlist *spotifyPlaylist) GetURL() string {
	if url, present := playlist.obj.ExternalURLs["spotify"]; present {
		return url
	}
	return ""
}

func (playlist *spotifyPlaylist) GetCoverURL() string {
	if len(playlist.obj.Images) > 0 {
		return playlist.obj.Images[0].URL
	}
	return ""
}

func (client *spotifyClient) GetPlaylists() (Playlists, error) {
	page := 0
	limit := SPOTIFY_PLAYLISTS_PAGE_LIMIT
	var playlists Playlists

	for {
		offset := page * (limit - 1)
		options := spotify.Options{Limit: &limit, Offset: &offset}
		simplePlaylistPage, err := client.CurrentUsersPlaylistsOpt(&options)
		if err != nil {
			return nil, err
		}

		// Only grab the first 49 playlists, use the last one for peeking
		for i, playlist := range simplePlaylistPage.Playlists {
			if i < limit {
				next := playlist
				playlists = append(playlists, NewPlaylist(&spotifyPlaylist{&next}))
			}
		}

		// Peeking to see if we should make another API call
		if len(simplePlaylistPage.Playlists) == limit {
			page++
		} else {
			return playlists, nil
		}
	}
}

func (client *spotifyClient) GetOwnPlaylistInfo(playlistId string) (*Playlist, error) {
	user, err := client.CurrentUser()
	if err != nil {
		return nil, err
	}

	return client.GetPlaylistInfo(user.ID, playlistId)
}

func (client *spotifyClient) GetPlaylistInfo(ownerId, playlistId string) (*Playlist, error) {
	log.Printf("Finding playlist %s belonging to user %s", spotify.ID(playlistId), ownerId)
	fullPlaylist, err := client.GetPlaylist(ownerId, spotify.ID(playlistId))
	if err != nil {
		return nil, err
	}

	return NewPlaylist(&spotifyPlaylist{&fullPlaylist.SimplePlaylist}), nil
}
