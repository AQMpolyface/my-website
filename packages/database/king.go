// packages/database/king.go
package database

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"website/packages/htmx"
	//"website/packages/database"
)

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("got the post req")
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

	//fmt.Println("password:", password, "username:", username)
	db, err := ConnectToDB()
	if err != nil {
		fmt.Println("error connecting to database")
		return
	}

	defer db.Close()

	valid, err := CheckUsername(db, username)
	if err != nil {
		fmt.Println("Error fetching after checkUsername database", err)
		errorMessage := htmx.ErrorRegister()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, errorMessage, http.StatusInternalServerError)
		return
	}
	if valid {
		fmt.Println("adding user to db")
		err = AddUser(db, username, password)
		fmt.Fprintf(w, `<h4 style="color:green;"> You have registered succesfully. You can now login safely, press on the login button</h3>`)
		if err != nil {
			log.Fatal(err)
		}
		response := htmx.SuccessRegister()
		fmt.Fprintf(w, response)
		return
	} else {
		fmt.Println("unauthorized user:", username)
		//w.WriteHeader(http.StatusForbidden)
		errorMessage := htmx.UnauthorizedRegister()
		fmt.Fprintf(w, errorMessage)
		return

	}

}

func PasswordRight(w http.ResponseWriter, r *http.Request) {

	db, err := ConnectToDB()
	if err != nil {
		fmt.Println("error connecting to db", err)
		return
	}
	defer db.Close()
	uuid, err := MakeUuid(db)
	if err != nil || uuid == "" {
		fmt.Println("error making uuid", err)
		return
	}
	cookie := http.Cookie{
		Name:     strings.TrimSpace("uuid"),
		Value:    strings.TrimSpace(uuid),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	// Response data to send back
	http.Redirect(w, r, "/protected", http.StatusSeeOther)
}
func ProtectionHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("uuid")
	if err != nil {
		fmt.Println("error getting cookie", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	valid, err := CheckCookie(cookie.Value)

	if !valid {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data, err := os.ReadFile("html/video/pickafterauth.html")
	if err != nil {
		fmt.Println("reading html/video/pickafterauth.html", err)
		http.Error(w, "Error readinghtml/video/pickafterauth.html", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(data))
}

func CheckCookie(cookie string) (bool, error) {
	db, err := ConnectToDB()
	if err != nil {
		fmt.Println("error connecting to db:", err)
		return false, err
	}

	defer db.Close()
	valid, err := CheckUuid(db, cookie)
	if err != nil {
		fmt.Println("error retrieving uuid:", err)
		return false, err
	}
	if valid {
		return true, nil
	} else if !valid {
		return false, nil
	} else {
		fmt.Println("non true/false :(")
		return false, nil
	}

}
