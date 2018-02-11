package models

type ClientPlaylist interface {
	GetID() string
	GetOwnerID() string
	GetName() string
	GetURL() string
	GetCoverURL() string
}

type Playlist struct {
	ID       string
	OwnerID  string
	Name     string
	URL      string
	CoverURL string
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

// func (playlists Playlists) MarshalJSON() ([]byte, error) {
// 	m := make(map[string]interface{})
// 	m["playlists"] = make([]map[string]string, len(playlists))
// 	mPlaylists := m["playlists"].([]map[string]string)
// 	for i, playlist := range playlists {
// 		mPlaylists[i] = make(map[string]string)
// 		mPlaylists[i]["id"] = playlist.ID
// 		mPlaylists[i]["ownerId"] = playlist.OwnerID
// 		mPlaylists[i]["name"] = playlist.Name
// 		mPlaylists[i]["url"] = playlist.URL
// 		mPlaylists[i]["coverUrl"] = playlist.CoverURL
// 	}

// 	return json.Marshal(m)
// }
