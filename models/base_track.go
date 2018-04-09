package models

type ClientTrack interface {
	GetID() string
	GetTitle() string
	GetArtist() string
}

type Track struct {
	ID     string
	Title  string
	Artist string
}

func NewTrack(track ClientTrack) *Track {
	return &Track{
		track.GetID(),
		track.GetTitle(),
		track.GetArtist(),
	}
}

type Tracks []*Track
