package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const savefile string = "/home/polyface/Desktop/go/site2/data.json"
const jsonDb string = "/home/polyface/Desktop/go/site2/db.json"

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)

	http.Handle("/uwu", http.FileServer(http.Dir("./")))

	fmt.Println("Server started at :80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	indexData, err := os.ReadFile("html/index.html")
	if err != nil {
		fmt.Println("error reading index.html", err)
		return

	}
	fmt.Fprint(w, string(indexData))
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
  aboutData, err := os.ReadFile("html/about.html")
  if err != nil {
    fmt.Println("error reading about.html", err)
    return
  }
fmt.Fprint(w, string(aboutData))

}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		fmt.Println(string(body))

		response := "uwu"
		fmt.Fprint(w, response)
	}
}

// C
