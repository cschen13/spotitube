package models

type ClientPlaylist interface {
	GetID() string
	GetOwnerID() string
	GetName() string
	GetURL() string
	GetCoverURL() string
}

type Playlist struct {
	ID       string `json:"id"`
	OwnerID  string `json:"ownerId"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	CoverURL string `json:"coverUrl"`
}

func NewPlaylist(playlist ClientPlaylist) *Playlist {
	return &Playlist{
		playlist.GetID(),
		playlist.GetOwnerID(),
		playlist.GetName(),
		playlist.GetURL(),
		playlist.GetCoverURL(),
	}
}

type Playlists []*Playlist
