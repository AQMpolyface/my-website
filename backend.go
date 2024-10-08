package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("Server started at :8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
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


