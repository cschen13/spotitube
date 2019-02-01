package models

type ClientTrack interface {
	GetID() string
	GetTitle() string
	GetArtist() string
}

type Track struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

func NewTrack(track ClientTrack) *Track {
	return &Track{
		track.GetID(),
		track.GetTitle(),
		track.GetArtist(),
	}
}

type Tracks []*Track
