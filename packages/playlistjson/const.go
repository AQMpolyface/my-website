// constants.go
package playlistjson

const (
	endpoint   = "https://api.spotify.com/v1/me/playlists"
	playlistUrl = "https://api.spotify.com/v1/playlists/"
)

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

