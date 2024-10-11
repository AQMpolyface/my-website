package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	//"os/exec"
	"os"
	"regexp"
	"strings"
	//"time"
	"website/packages/playlistjson"
)

const visits string = "visits.txt"

func main() {

	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/privacy_policy", privacyPolicyHandler)
	http.HandleFunc("/blahaj", blahajHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/projects/playlistjson", playlistjsonHandler)
	//http.HandleFunc("/submit-playlist-json", playlistjson.PlaylistJson)
	http.HandleFunc("/submit-playlist-json", playlistjsonHandlerPost)
	http.HandleFunc("/submit", formHandler)
//	http.HandleFunc("/temp", downloadHandler)
	http.Handle("/temp/", http.StripPrefix("/temp/", http.FileServer(http.Dir("temp"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("Server started at :8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/temp/"):] // Get the file name from the URL
	fmt.Println(fileName)

	filePath := "temp/" + fileName         // Construct the full file path
	fmt.Println(filePath)

	/*file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()*/

	//force download
/*w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
w.Header().Set("Content-Type", "application/octet-stream")*/
	// Serve the file
	body, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("error reading" + filePath + ":", err)
				return
		}
	fmt.Fprint(w, string(body))
//http.ServeFile(w, r, filePath)
}

func playlistjsonHandlerPost(w http.ResponseWriter, r *http.Request) {

	//allowed Header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//TEAPOT LETS GOO :3
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

	//Parsing values from what htmx sent (url format)
	values, err := url.ParseQuery(decodedMessage)
	if err != nil {
		errorMessage := fmt.Sprintf("Error parsing query string:%s", err)
		fmt.Fprint(w, "<p>"+errorMessage+"</p>")
		return
	}

	token := values.Get("token")
	verbose := values.Get("verbose")
	if verbose != "" {
		playlistjson.PlaylistJsonSocket(w, r, token)

	} else {
		playlistjson.PlaylistJson(w, r, token)
	}
}

func playlistjsonHandler(w http.ResponseWriter, r *http.Request) {
	//serve the page for playlistjson
	playlistjsonHtml, err := os.ReadFile("html/projects/playlistjson.html")
	if err != nil {
		fmt.Println("error reading playlistjson.html:", err)
		return
	}
	fmt.Fprint(w, string(playlistjsonHtml))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	contactData, err := os.ReadFile("html/contact.html")
	if err != nil {
		fmt.Println("error reading contact.html:", err)
		return
	}
	fmt.Fprint(w, string(contactData))

}

func formHandler(w http.ResponseWriter, r *http.Request) {
	//contact form  todo: ( add an automail soon)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusTeapot)
		fmt.Println("teapot party")
		fmt.Fprintln(w, "I'm a teapot!")
	}

	if r.Method == http.MethodPost {
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

		fmt.Println(string(decodedMessage))
		values, err := url.ParseQuery(decodedMessage)
		if err != nil {
			fmt.Println("Error parsing query string:", err)
			return
		}
		//regex to check if its an email, idk if it works xd
		emailRegex := `(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\$$\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)$$)`
		emailTrue, err := regexp.MatchString(emailRegex, values.Get("email"))
		if err != nil {
			fmt.Println("error parsing email:", err)
			return
		}
		if !emailTrue {
			errorMessage := fmt.Sprintf("<p>Error: %s isnt an email! please enter a valid email.</p>", values.Get("email"))
			fmt.Fprint(w, errorMessage)
		}

		messageFileHandler, err := os.OpenFile("messages.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("error opening message file:", err)
		}
		defer messageFileHandler.Close()
		decodedMessage += "\n"
		messageFileHandler.WriteString(decodedMessage)

		//htmx response, with js clock
		responseMessage := `<h3>Thanks for your submission. it is now <span id="time">uwu</span> in my timezone, so i will see when i can get back at you!</h3>
		<script>	function sleep(ms) {
		return new Promise(resolve => setTimeout(resolve, ms));
	}

	async function time() {
	timeNow = document.getElementById("time");

	const options = {
	timeZone: 'Europe/Zurich',
	 dateStyle: 'full',
	 timeStyle: 'long',
	 /*hour: '2-digit',
	minute: '2-digit',
	second: '2-digit',
	hour12: false*/
	};

	const formatter = new Intl.DateTimeFormat('en-US', options);
	while (true) {

		let date = new Date();
		let formattedDate = formatter.format(date);
		timeNow.innerHTML = formattedDate;
		await sleep(1000);

	}

	}

	time()
	</script>
		`

		fmt.Fprint(w, responseMessage)
	}
}

func privacyPolicyHandler(w http.ResponseWriter, r *http.Request) {

	//privacy policy
	privacyPolicyData, err := os.ReadFile("html/privacy_policy.html")
	if err != nil {
		fmt.Println("error reading privacy_policy.html", err)
		return
	}
	fmt.Fprint(w, string(privacyPolicyData))
}

func blahajHandler(w http.ResponseWriter, r *http.Request) {

	blahajData, err := os.ReadFile("html/blahaj.html")
	if err != nil {
		fmt.Println("error reading blahaj.html", err)
		return
	}

	fmt.Fprint(w, string(blahajData))

}
func projectsHandler(w http.ResponseWriter, r *http.Request) {

	//project page
	projectsData, err := os.ReadFile("html/projects.html")
	if err != nil {
		fmt.Println("error reading projects.html", err)
		return
	}

	fmt.Fprint(w, string(projectsData))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	//indedx page
	indexData, err := os.ReadFile("html/index.html")
	if err != nil {
		fmt.Println("error reading index.html", err)
		return
	}
	visitsFile, err := os.Open(visits)
	if err != nil {
		log.Fatal("error opening message visits.txt:", err)
	}
	defer visitsFile.Close()

	numberOfVisitsByte, err := io.ReadAll(visitsFile)
	if err != nil {
		log.Fatal("error reading visits file:", err)
	}
	str := string(numberOfVisitsByte)
	str = strings.TrimSpace(str)

	numberOfVisits, err := strconv.Atoi(string(str))
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}
	numberOfVisits += 1
	fmt.Println(numberOfVisits)

	indec, err := os.OpenFile(visits, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("error opening visits file:", err)
	}
	defer indec.Close()
	_, err = indec.WriteString(strconv.Itoa(numberOfVisits))
	if err != nil {
		log.Fatal("error writing to file:", err)
	}
	fmt.Fprint(w, string(indexData))
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {

	//about page
	aboutData, err := os.ReadFile("html/about.html")
	if err != nil {
		fmt.Println("error reading about.html", err)
		return
	}
	fmt.Fprint(w, string(aboutData))

}
