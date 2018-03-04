package models

import (
	"encoding/json"
	"errors"
	youtube "google.golang.org/api/youtube/v3"
	"log"
	"strings"
)

func (client *youtubeClient) GetTrackByID(videoId string) (*Track, error) {
	return nil, errors.New("Unimplemented")
}

func (client *youtubeClient) GetTracks(playlist *Playlist) (Tracks, error) {
	return nil, errors.New("Unimplemented")
}

func (client *youtubeClient) InsertTrack(playlist *Playlist, track Track) (bool, error) {
	playlistItem := &youtube.PlaylistItem{}
	videoId, err := client.searchForMatchingVideo(track)
	if err != nil {
		log.Printf("youtube: Error searching for track")
		return false, err
	} else if videoId == "" {
		log.Printf("youtube: Zero search results for track %s - %s", track.Artist, track.Title)
		return false, nil
	}

	properties := map[string]string{
		"snippet.playlistId":         playlist.ID,
		"snippet.resourceId.kind":    "youtube#video",
		"snippet.resourceId.videoId": videoId,
	}
	res := createResource(properties)
	if err := json.NewDecoder(strings.NewReader(res)).Decode(&playlistItem); err != nil {
		log.Printf("youtube: Failed to decode JSON into playlist item resource")
		return false, err
	}

	call := client.PlaylistItems.Insert("id", playlistItem)
	_, err = call.Do()
	return true, err
}

func (client *youtubeClient) searchForMatchingVideo(track Track) (videoId string, err error) {
	call := client.Search.List("snippet").Q(track.Artist + " " + track.Title + " video").Type("video")
	response, err := call.Do()
	if err != nil {
		return
	}

	if len(response.Items) > 0 {
		videoId = response.Items[0].Id.VideoId
	}
	return
}
