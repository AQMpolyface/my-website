package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"website/packages/database"
	"website/packages/htmx"
	"website/packages/playlistjson"
	"website/packages/video"
)

const visits string = "visits.txt"

var fileNames string
var playlistFile string

func main() {
	http.HandleFunc("/", mainHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/vid/", serveProtectedFiles)
	fmt.Println("Server started at :8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func serveProtectedFiles(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("uuid")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
		return
	}
	db, err := database.ConnectToDB()
	if err != nil {
		http.Error(w, "error connecting to db", http.StatusInternalServerError)
		fmt.Println("error connecting to db", err)
		return
	}
	defer db.Close()
	valid, err := database.CheckUuid(db, cookie.Value)
	if err != nil {
		fmt.Println("error retrieving uuid from db:", err)
		return
	}
	// fmt.Println(cookie.Value)
	if !valid {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	} else {
		http.StripPrefix("/vid", http.FileServer(http.Dir("vid"))).ServeHTTP(w, r)
		return
	}
}
func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		indexHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/privacy_policy":
		privacyPolicyHandler(w, r)
	case "/blahaj":
		blahajHandler(w, r)
	case "/protected":
		database.ProtectionHandler(w, r)
	case "/about":
		aboutHandler(w, r)
	case "/projects":
		projectsHandler(w, r)
	case "/projects/playlistjson":
		playlistjsonHandler(w, r)
	case "/api/submit-playlist-json":
		playlistjsonHandlerPost(w, r)
	case "/api/submit-password":
		passwordPost(w, r)
	case "/api/submit-registration":
		database.RegisterPost(w, r)
	case "/api/submit":
		formHandler(w, r)
	case "/uwu":
		uwuHandler(w, r)
	case "/api/uwunumber":
		uwuNumberHandler(w, r)
	case "/projects/temp/":
		serveFileHandler(w, r)
	case "/auth":
		serveAuthHandler(w, r)
	case "/api/register":
		registerHandler(w, r)
	case "/api/relogin":
		reloginHandler(w, r)
	case "/protected/king":
		video.KingHandler(w, r)
	case "/protected/theoden":
		video.TheodenHandler(w, r)
	case "/protected/fellowship":
		video.FellowshipHandler(w, r)
	case "/protected/towers":
		video.TowerHandler(w, r)
	default:
		notFoundHandler(w, r)
	}
}

func reloginHandler(w http.ResponseWriter, r *http.Request) {
	stringuwu := htmx.ReturnReloginString()
	fmt.Fprintf(w, stringuwu)

}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	stringuwu := htmx.ReturnRegisterString()
	fmt.Fprintf(w, stringuwu)

}
func passwordPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//TEAPOT  nbr 3 LETS GOO :3
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusTeapot)
		fmt.Println("teapot party")
		fmt.Fprintln(w, "I'm a teapot!")
		return
	}

	r.ParseForm()
	password := r.FormValue("password")
	username := r.FormValue("username")
	db, err := database.ConnectToDB()
	if err != nil {
		fmt.Println("error connecting to the db", err)
		return
	}
	defer db.Close()
	valid, err := database.CheckUserCredentials(db, username, password)
	if err != nil {
		fmt.Println("error checking credentials", err)
		return
	}
	if valid {
		fmt.Println("Password is valid!")
		db, err := database.ConnectToDB()
		if err != nil {
			fmt.Println("error connecting to db", err)
			return
		}
		defer db.Close()
		var uuid string
		username = strings.TrimSpace(username)
		//get the uuid of the logged in person
		err = db.QueryRow("SELECT uuid FROM authentification WHERE username = ?", username).Scan(&uuid)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No rows returned for username:", username)
				return
			}
			fmt.Println("Error executing query:", err)
			return
		}
		cookie := http.Cookie{
			Name:     strings.TrimSpace("uuid"),
			Value:    strings.TrimSpace(uuid),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		}
		http.SetCookie(w, &cookie)
		fmt.Fprintf(w, `<h4 style="color:green;">You are logged in. you can you go to <a href="https://polyface.ch/protected">https://polyface.ch/protected</a></h4>`)
		http.Redirect(w, r, "/protected", http.StatusSeeOther)
		//err := database.AddUser(db, username, password)
		return
	} else {
		fmt.Println("Invalid password.")
		errorMessage := htmx.BadPassword()
		fmt.Fprintf(w, errorMessage)
		return
	}
}

func serveAuthHandler(w http.ResponseWriter, r *http.Request) {
	var hasValidCookie bool
	cookie, err := r.Cookie("uuid")
	if err != nil && err != http.ErrNoCookie {
		log.Fatal("error checking cookie")
	}
	if cookie != nil {
		db, err := database.ConnectToDB()
		if err != nil {
			http.Error(w, "error connecting to database", http.StatusInternalServerError)
			return
		}
		defer db.Close()
		hasValidCookie, err = database.CheckUuid(db, cookie.Value)
		if err != nil {
			http.Error(w, "error retrieving cookie", http.StatusInternalServerError)
			return
		}

		if hasValidCookie {
			http.Redirect(w, r, "/protected", http.StatusSeeOther)
			return
		}
	}

	kingData, err := os.ReadFile("html/video/pickvid.html")
	if err != nil {
		fmt.Printf("error readinf %s: %s", "html/video/pickvid.html", err)
		http.Error(w, "Error reading html/video/pickvid.html", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(kingData))
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	filename2 := r.URL.Path[len("/projects/temp/"):]
	fmt.Println("filename2 = ", filename2)
	fmt.Println("playlistFile = ", playlistFile)

	//fixed directory traversal uwu
	filename2 = filepath.Clean(filename2)
	playlistFile = filepath.Clean(playlistFile)
	osOpenFile := playlistFile
	filejsonData, err := os.ReadFile(osOpenFile)
	fmt.Println("filejsonData = ", string(filejsonData))
	if err != nil {
		fmt.Printf("error readinf %s: %s", filename2, err)
		http.Error(w, "Error reading "+filename2, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filename2)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(filejsonData)))

	// Write to response
	_, err = w.Write(filejsonData)
	if err != nil {
		fmt.Println("error writing response", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(filejsonData))
	// remove the file after download :3
	defer os.Remove(osOpenFile)
}

func uwuNumberHandler(w http.ResponseWriter, r *http.Request) {
	uwuNumberData, err := os.ReadFile("uwunumber.txt")
	if err != nil {
		fmt.Println("reading uwu number", err)
		http.Error(w, "Error reading uwu number", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(uwuNumberData))
}

func uwuHandler(w http.ResponseWriter, r *http.Request) {
	//read uwu.html
	uwuData, err := os.ReadFile("html/uwu.html")
	if err != nil {
		fmt.Println("reading uwu.html", err)
		http.Error(w, "Error reading uwu.html", http.StatusInternalServerError)
		return
	}
	uwunumber, err := os.ReadFile("uwunumber.txt")
	if err != nil {
		fmt.Println("reading uwu.txt", err)
		http.Error(w, "Error reading uwu.txt", http.StatusInternalServerError)
	}
	//convert number of visits on /uwu to string
	numberOfUwu, err := strconv.Atoi(string(uwunumber))
	if err != nil {
		fmt.Println("reading uwu.txt", err)
		http.Error(w, "Error reading uwu.txt", http.StatusInternalServerError)
	}
	numberOfUwu++
	fmt.Println(numberOfUwu)
	//tracking number of person who pressed on the yes button on the index page
	os.WriteFile("uwunumber.txt", []byte(fmt.Sprintf("%d", numberOfUwu)), 0644)

	fmt.Fprint(w, string(uwuData))
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
		fmt.Fprintf(w, "I'm a teapot!")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("unable to read request body")
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// get the token
	decodedMessage, err := url.QueryUnescape(string(body))
	if err != nil {
		fmt.Println("decoding message", err)
		http.Error(w, "decoding message", http.StatusInternalServerError)
		return
	}
	fmt.Println(decodedMessage)

	//Parsing values from what htmx sent (url format)
	values, err := url.ParseQuery(decodedMessage)
	if err != nil {
		fmt.Println("parsing url values", err)
		http.Error(w, "parsing url values", http.StatusInternalServerError)
		return
	}

	token := values.Get("token")
	//call a function from packages/playlistjson that returns the exact url and just the filename
	playlistjson.PlaylistJson(w, r, token)
	// debug	fmt.Println("playlistfile at the end of playlistjsonpost:", playlistFile)

}

func playlistjsonHandler(w http.ResponseWriter, r *http.Request) {
	//serve the page for playlistjson
	playlistjsonHtml, err := os.ReadFile("html/projects/playlistjson.html")
	if err != nil {
		fmt.Println("reading playlistjson.html", err)
		http.Error(w, "Error reading playlistjson.html", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(playlistjsonHtml))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	//serve the contact page
	contactData, err := os.ReadFile("html/contact.html")
	if err != nil {
		fmt.Println("error reading contact.html", err)
		http.Error(w, "Error reading contact.html", http.StatusInternalServerError)
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
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error getting request body:", err)
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		//decode the message to write it to messages.txt
		decodedMessage, err := url.QueryUnescape(string(body))
		if err != nil {
			fmt.Println("error decoding message", err)
			http.Error(w, "Error decoding message", http.StatusInternalServerError)
			return
		}

		fmt.Println(string(decodedMessage))
		values, err := url.ParseQuery(decodedMessage)
		if err != nil {
			fmt.Println("error parsing query string", err)
			http.Error(w, "Error parsing query string", http.StatusInternalServerError)
			return
		}
		//regex to check if its an email, no clue if it works
		emailRegex := `(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\$$\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)$$)`
		emailTrue, err := regexp.MatchString(emailRegex, values.Get("email"))
		if err != nil {
			fmt.Println("error parsing email", err)
			http.Error(w, "Error parsing email", http.StatusInternalServerError)
			return
		}
		if !emailTrue {
			errorMessage := fmt.Sprintf("<p>Error: %s isnt an email! please enter a valid email.</p>", values.Get("email"))
			fmt.Fprint(w, errorMessage)
			return
		}

		//write the message and the email
		messageFileHandler, err := os.OpenFile("messages.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("error opening message file:", err)
		}
		defer messageFileHandler.Close()
		decodedMessage += "\n"
		messageFileHandler.WriteString(decodedMessage)

		//htmx response, with js clock
		responseMessage := htmx.GetSubmissionSuccess()
		fmt.Fprint(w, responseMessage)
	}
}

func privacyPolicyHandler(w http.ResponseWriter, r *http.Request) {

	//privacy policy
	privacyPolicyData, err := os.ReadFile("html/privacy_policy.html")
	if err != nil {
		fmt.Println("error reading privacy_policy.html", err)
		http.Error(w, "Error reading privacy_policy.html", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(privacyPolicyData))
}

func blahajHandler(w http.ResponseWriter, r *http.Request) {
	//super duper secret page
	blahajData, err := os.ReadFile("html/blahaj.html")
	if err != nil {
		fmt.Println("error reading blahaj.html", err)
		http.Error(w, "Error reading blahaj.html", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(blahajData))

}
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	//project page
	projectsData, err := os.ReadFile("html/projects.html")
	if err != nil {
		fmt.Println("error reading projects.html", err)
		http.Error(w, "Error reading projects.html", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(projectsData))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	notFoundData, err := os.ReadFile("html/404.html")
	if err != nil {
		fmt.Println("error reading 404.html", err)
		http.Error(w, "Error reading 404.html", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(notFoundData)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	//indedx page
	indexData, err := os.ReadFile("html/index.html")
	if err != nil {
		fmt.Println("error reading index.html", err)
		http.Error(w, "Error reading index.html", http.StatusInternalServerError)
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
		fmt.Println("error converting string to int", err)
		http.Error(w, "Error converting string to int", http.StatusInternalServerError)
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
		http.Error(w, "Error reading about.html", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(aboutData))
}
