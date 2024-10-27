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
	link := "projects/temp/" + randData
	//	playlistFile := "projects/" + link

	client = &http.Client{}
	playlistFile2, err := os.OpenFile("playlist.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer playlistFile2.Close()
	startPlaylist := `{`
	if _, err := playlistFile2.WriteString(startPlaylist + "\n"); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	for _, playlist := range data.Items {
		//	debugging	fmt.Printf("Playlist Name: %s, ID: %s\n", playlist.Name, playlist.Id)

		playlistName := fmt.Sprintf(`"playlistname" : "%s",
    "playlistis" : "%s",
    "items" [`, playlist.Name, playlist.Id)
		if _, err := playlistFile2.WriteString(playlistName + "\n"); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
		fields := "items.track(name,id)"
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?fields=%s", playlist.Id, fields)
		playlistReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		playlistReq.Header.Set("Authorization", "Bearer "+token)
		//playlistReq.Header.Set("User-Agent", "curl/7.64.1") silly attempt 4
		client := &http.Client{}
		playlistResp, err := client.Do(playlistReq)
		if err != nil {
			log.Fatal(err)
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

		for _, item := range musicData.Items {
			//			time.Sleep(time.Second * 1) debugging 5
			fmt.Println(item.Track.Name)
			fmt.Println(item.Track.ID)
			songName := fmt.Sprintf(` {
  "song" : "%s",
  "id" : "%s"
        }`, item.Track.Name, item.Track.ID)
			if _, err := playlistFile2.WriteString(songName + "\n"); err != nil {
				log.Fatalf("Failed to write to file: %v", err)
			}
		}

		endSong := `]`
		if _, err := playlistFile2.WriteString(endSong + "\n"); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}

	}
	/*data2, err := os.ReadFile(playlistFile)
	if err != nil {
		fmt.Println("error uwu")
		return "", ""
	}
	fmt.Println(string(data2)) */
	fmt.Println("link is equal to ", link)
	responseMessage := fmt.Sprintf(`<h4> You can download your json file <a href="%s">here</a></h5>
		<h5 style="color:red;">Warning: the file will be deleted after downloading it, you will have to redo the process if you lose the file</h5>`, link)
	fmt.Println("uwu world")
	fmt.Fprint(w, responseMessage)
	//filename := "projects/" + link
	/* 	fileMap := map[string]string{
		"filename": filename,
		"url":      fmt.Sprintf("/files/%s", filename),
	}*/

	return link, randData

}
