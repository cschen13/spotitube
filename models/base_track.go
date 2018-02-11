package models

type Track struct {
	Title  string
	Artist string
}

type Tracks []*Track

type ClientTrack interface {
	GetTitle() string
	GetArtist() string
}

func NewTrack(track ClientTrack) *Track {
	return &Track{
		track.GetTitle(),
		track.GetArtist(),
	}
}

// func (tracks Tracks) MarshalJSON() ([]byte, error) {
// 	m := make(map[string]interface{})
// 	m["tracks"] = make([]map[string]string, len(tracks))
// 	mTracks := m["tracks"].([]map[string]string)
// 	for i, track := range tracks {
// 		mTracks[i] = make(map[string]string)
// 		mTracks[i]["title"] = track.Title
// 		mTracks[i]["artist"] = track.Artist
// 	}

// 	return json.Marshal(m)
// }
