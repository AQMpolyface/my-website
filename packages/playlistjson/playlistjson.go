package playlistjson

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var playlistfile string
var playlistfilev2 string

func PlaylistJson(w http.ResponseWriter, r *http.Request, token string) (string, string) {

	req, err := http.NewRequest("GET", endpoint, nil)

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("error reading the body: %s", err)
		fmt.Fprint(w, "<p>"+errorMessage+"</p>")
		return "", ""

	}

	if resp.Header.Get("Retry-After") != "" {
		errorMessage := fmt.Sprintf("rate limited by the spotify api, you ran the code too much, retry in %s:\n %s", resp.Header.Get("Retry-After"), string(body))
		fmt.Fprint(w, "<p>"+errorMessage+"</p>")
		return "", ""

	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprint(w, "token probably need to be refreshed:"+string(body))
		return "", ""

	}

	var data PlaylistInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error unmarshaling data", err)
		return "", ""

	}

	if err := os.MkdirAll("temp", os.ModePerm); err != nil {
		fmt.Println("Error creating temp directory:", err)
		return "", ""

	}
	randData := strconv.FormatInt(time.Now().UnixNano(), 10) + ".json"
	link := "temp/" + randData
	playlistFile := "projects/" + link

	client = &http.Client{}
	for _, playlist := range data.Items {
		fmt.Printf("Playlist Name: %s, ID: %s\n", playlist.Name, playlist.Id)

		fields := "items.track(name,id)"
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?fields=%s", playlist.Id, fields)
		playlistReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
			return "", ""

		}
		playlistReq.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		playlistResp, err := client.Do(playlistReq)
		if err != nil {
			log.Fatal(err)
			return "", ""

		}
		defer playlistResp.Body.Close()

		playlistBody, err := io.ReadAll(playlistResp.Body)
		if err != nil {
			fmt.Println("error reading body", err)
			return "", ""

		}

		var musicData PlaylistResponse
		err = json.Unmarshal(playlistBody, &musicData)
		if err != nil {
			fmt.Println("error unmarshaling playlist content:", err)
			return "", ""

		}

		fileWriter, err := os.OpenFile(playlistFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("error opening playlist.json file:", err)
			fmt.Println("uwu")
			return "", ""

		}
		defer fileWriter.Close()

		jsonData, err := json.MarshalIndent(musicData, "", "  ")
		if err != nil {
			log.Fatal("error marshaling data to JSON:", err)
			return "", ""

		}

		_, err = fileWriter.Write(jsonData)
		if err != nil {
			log.Fatal("error writing to the .json file:", err)
			return "", ""

		}
	}
	data2, err := os.ReadFile(playlistFile)
	if err != nil {
		fmt.Println("error uwu")
		return "", ""
	}
	fmt.Println(string(data2))
	fmt.Println("link is equal to ", link)
	responseMessage := `<h4> You can download your json file <a href="temp/tempfile">here</a></h5>
		<h5 style="color:red;">Warning: the file will be deleted after downloading it, you will have to redo the process</h5>`
	fmt.Println("uwu world")
	fmt.Fprint(w, responseMessage)
	return "projects/" + link, randData

}
