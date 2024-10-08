package main

import (
	"io"
  "fmt"
	"log"
	"net/http"
  "net/url"
  "os"
	"time"
)

func main() {
http.HandleFunc("/contact", contactHandler)
  http.HandleFunc("/privacy_policy", privacyPolicyHandler)
  http.HandleFunc("/blahaj", blahajHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
  http.HandleFunc("/projects", projectsHandler)
 http.HandleFunc("/submit", formHandler)
  http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("Server started at :8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
 contactData, err := os.ReadFile("html/contact.html")
if err != nil {
    fmt.Println("error reading contact.html", err)
    return
}
fmt.Fprint(w, string(contactData))

}

func formHandler(w http.ResponseWriter, r *http.Request) {

w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "POST")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

 if r.Method != http.MethodPost {
   w.WriteHeader(http.StatusTeapot) 
    fmt.Println("teapot party")
    fmt.Fprintln(w, "I'm a teapot")  
  }

  if r.Method == http.MethodPost {
        // Read the request body
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Unable to read request body", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()
    fmt.Println(string(body))
  //      email := r.FormValue("email")
    //    message := r.FormValue("message")

//dataToWrite := fmt.Sprintf("email: %s\nmessage: %s", email, message)
    decodedMessage, err := url.QueryUnescape(string(body))
      if err != nil {
        fmt.Println("error decoding message", err)
      return
    }
    os.WriteFile("messages.txt", []byte(decodedMessage), 0644)
        // Create a response messag
    timeString:= time.Now().Format("15:04:05")
        responseMessage := fmt.Sprintf("<h3>Thanks for your submission. it is now %s in my timezone, so i will see when i can get back at you!</h3>", timeString)



//        w.Write([]byte(responseMessage))

     fmt.Fprint(w, responseMessage) 
}
}

func privacyPolicyHandler(w http.ResponseWriter, r *http.Request) {
  privacyPolicyData, err := os.ReadFile("html/privacy_policy.html")
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


