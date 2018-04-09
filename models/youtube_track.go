package models

import (
	"encoding/json"
	youtube "google.golang.org/api/youtube/v3"
	"log"
	"strings"
)

func (client *youtubeClient) InsertTrack(playlist *Playlist, track *Track) (bool, error) {
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

	// Don't actually need snippet part, but API throws an unexpectedPart error
	// for just "id"
	call := client.PlaylistItems.Insert("id,snippet", playlistItem)
	_, err = call.Do()
	return true, err
}

func (client *youtubeClient) searchForMatchingVideo(track *Track) (videoId string, err error) {
	call := client.Search.List("snippet").Q(track.Artist + " " + track.Title + " music video").Type("video")
	response, err := call.Do()
	if err != nil {
		return
	}

	if len(response.Items) > 0 {
		videoId = response.Items[0].Id.VideoId
	}
	return
}
