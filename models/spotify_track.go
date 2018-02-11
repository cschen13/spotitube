package models

import (
	"errors"
	"github.com/zmb3/spotify"
)

const SPOTIFY_PLAYLIST_TRACKS_PAGE_LIMIT = 100

type spotifyTrack struct {
	obj *spotify.SimpleTrack
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

func (client *spotifyClient) GetTracks(playlist *Playlist, page string) (Tracks, error) {
	ownerId := playlist.OwnerID
	playlistId := spotify.ID(playlist.ID)
	trackPage, err := client.GetPlaylistTracksOpt(ownerId, playlistId, nil, "total")
	if err != nil {
		return nil, err
	}

	tracks := make(Tracks, trackPage.Total)
	limit := SPOTIFY_PLAYLIST_TRACKS_PAGE_LIMIT
	offset := 0
	for page := 0; offset < len(tracks); page++ {
		offset = (page - 1) * limit
		options := spotify.Options{Limit: &limit, Offset: &offset}
		trackPage, err = client.GetPlaylistTracksOpt(ownerId, playlistId, &options, "")
		if err != nil {
			return nil, err
		}

		for i, track := range trackPage.Tracks {
			newTrack := track.Track.SimpleTrack
			tracks[offset+i] = NewTrack(&spotifyTrack{&newTrack})
		}
	}

	return tracks, nil
}

func (client *spotifyClient) InsertTrack(playlist *Playlist, track Track) (bool, error) {
	return false, errors.New("Unimplemented")
}
