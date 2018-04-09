package models

import (
	"github.com/zmb3/spotify"
	"log"
)

const SPOTIFY_PLAYLIST_TRACKS_PAGE_LIMIT = 100

type spotifyTrack struct {
	obj *spotify.SimpleTrack
}

func (track *spotifyTrack) GetID() string {
	return track.obj.ID.String()
}

func (track *spotifyTrack) GetTitle() string {
	return track.obj.Name
}

func (track *spotifyTrack) GetArtist() string {
	if len(track.obj.Artists) > 0 {
		return track.obj.Artists[0].Name
	}
	return ""
}

func (client *spotifyClient) GetTrackByID(id string) (*Track, error) {
	fullTrack, err := client.GetTrack(spotify.ID(id))
	if err != nil {
		return nil, err
	}

	return NewTrack(&spotifyTrack{&fullTrack.SimpleTrack}), nil
}

func (client *spotifyClient) GetTracks(playlist *Playlist) (Tracks, error) {
	ownerId := playlist.OwnerID
	playlistId := spotify.ID(playlist.ID)
	trackPage, err := client.GetPlaylistTracksOpt(ownerId, playlistId, nil, "total")
	if err != nil {
		log.Printf("Error getting playlist tracks")
		return nil, err
	}

	tracks := make(Tracks, trackPage.Total)
	limit := SPOTIFY_PLAYLIST_TRACKS_PAGE_LIMIT
	offset := 0
	for page := 0; offset < len(tracks); page++ {
		offset = page * limit
		options := spotify.Options{Limit: &limit, Offset: &offset}
		trackPage, err = client.GetPlaylistTracksOpt(ownerId, playlistId, &options, "")
		if err != nil {
			log.Printf("Error getting playlist tracks")
			return nil, err
		}

		for i, track := range trackPage.Tracks {
			newTrack := track.Track.SimpleTrack
			tracks[offset+i] = NewTrack(&spotifyTrack{&newTrack})
		}
	}

	return tracks, nil
}
