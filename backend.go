package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
  http.HandleFunc("/privacy_policy", privacyPolicyHandler)
  http.HandleFunc("/blahaj", blahajHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
  http.HandleFunc("/projects", projectsHandler)
 
  http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("Server started at :8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}
func privacyPolicyHandler(w http.ResponseWriter, r *http.Request) {
  privacyPolicyData, err := os.ReadFile("html/blahaj.html")
if err != nil {
    fmt.Println("error reading blahaj.html", err)
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
  projectsData, err := os.ReadFile("html/projects.html")
      if err != nil {
      fmt.Println("error reading projects.html", err)
        return
  }

  fmt.Fprint(w, string(projectsData))
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


