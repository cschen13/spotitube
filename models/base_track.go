package models

type Track struct {
	ID     string
	Title  string
	Artist string
}

type Tracks []*Track

type ClientTrack interface {
	GetID() string
	GetTitle() string
	GetArtist() string
}

func NewTrack(track ClientTrack) *Track {
	return &Track{
		track.GetID(),
		track.GetTitle(),
		track.GetArtist(),
	}
}
