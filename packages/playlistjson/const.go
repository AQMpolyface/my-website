// constants.go
package playlistjson

type Playlist struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type PlaylistInfo struct {
	Items []Playlist `json:"items"`
}

type PlaylistResponse struct {
	Items []TrackItem `json:"items"`
}

type TrackItem struct {
	Track Track `json:"track"`
}

type Track struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// defining a few cowonst
const endpoint string = "https://api.spotify.com/v1/me/playlists"

const playlistUrl string = "https://api.spotify.com/v1/playlists/"
