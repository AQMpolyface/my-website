package playlist-json

import (
	"net/http"
	"encoding/json"
	"io"
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
const endpoint string = "https://api.spotify.com/v1/me/playlists"

const playlistUrl string = "https://api.spotify.com/v1/playlists/"

func playlistJson(w http.ResponseWriter, r *http.Request) {

	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusTeapot)
		fmt.Println("teapot party")
		w.WriteHeader(418)
		fmt.Fprintln(w, "I'm a teapot!")
		}

	
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		decodedMessage, err := url.QueryUnescape(string(body))
		if err != nil {
			fmt.Println("error decoding message", err)
			return

    }
	fmt.Println("decodedMessage")

    values, err := url.ParseQuery(decodedMessage)
	if err != nil {
		errorMessage := "Error parsing query string:" + err
		fmt.Fprint("<p>" + errorMessage + "</p>")
		return
	}

	token := values.token
	if values.verbose != "" {
		var verbose bool = true
	}

	req, err := http.NewRequest("GET", endpoint, nil)

	fmt.Println(token)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("error reading the body: %s", err)
		fmt.Fprint(w, "<p>" + errorMessage + "</p>")
	}

	if resp.Header.Get("Retry-After") != "" {
		errorMessage := fmt.Sprintf("rate limited by the spotify api, you ran the code too much, retry in %s:\n %s", resp.Header.Get("Retry-After"), string(body))
		fmt.Fprint(w, "<p>" + errorMessage + "</p>")
		return

	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprint(w, "token probably need to be refreshed:" +  string(body))
		return
	}

	var data PlaylistInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error unmarshaling data", err)
		return
	}

	client = &http.Client{}
	for _, playlist := range data.Items {
		fmt.Printf("Playlist Name: %s, ID: %s\n", playlist.Name, playlist.Id)

		fields := "items.track(name,id)"
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?fields=%s", playlist.Id, fields)
		playlistReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		playlistReq.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		playlistResp, err := client.Do(playlistReq)
		if err != nil {
			log.Fatal(err)
		}
		defer playlistResp.Body.Close()

		playlistBody, err := io.ReadAll(playlistResp.Body)
		if err != nil {
			fmt.Println("error reading body", err)
			return
		}

		var musicData PlaylistResponse
		err = json.Unmarshal(playlistBody, &musicData)
		if err != nil {
			fmt.Println("error unmarshaling playlist content:", err)
			return
		}


    for _, item := range musicData.Items {
			fmt.Println(item.Track.Name)
			fmt.Println(item.Track.ID)
		}

		fileWriter, err := os.OpenFile(playlistFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("error opening playlist.json file:", err)
			fmt.Println("uwu")
		}
		defer fileWriter.Close()

		jsonData, err := json.MarshalIndent(musicData, "", "  ")
		if err != nil {
			log.Fatal("error marshaling data to JSON:", err)
		}

		_, err = fileWriter.Write(jsonData)
		if err != nil {
			log.Fatal("error writing to playlist.json:", err)
		}
	}
}
