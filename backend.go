package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// "strings"
)

const savefile string = "/home/polyface/Desktop/go/site2/data.json"
const jsonDb string = "/home/polyface/Desktop/go/site2/db.json"

func main() {

	http.HandleFunc("/", indexHandler)

	http.Handle("/uwu", http.FileServer(http.Dir("./")))

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	indexData, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Println("error reading index.html", err)
		return

	}
	fmt.Fprint(w, string(indexData))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	//var option string

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