package video

import (
	"fmt"
	"net/http"
	"os"
	"website/packages/database"
)

func KingHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("error connecting to db", err)
		return
	}
	defer db.Close()
	valid, err := database.CheckUuid(db, cookie.Value)
	if err != nil {
		fmt.Println("error retrieving uuid from db:", err)
		return
	}
	if !valid {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	} else {
		data, err := os.ReadFile("html/video/king.html")
		if err != nil {
			http.Error(w, "error reading file king.html", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(data))
	}
}

func TowerHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("error connecting to db", err)
		return
	}
	defer db.Close()
	valid, err := database.CheckUuid(db, cookie.Value)
	if err != nil {
		fmt.Println("error retrieving uuid from db:", err)
		return
	}
	if !valid {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	} else {
		data, err := os.ReadFile("html/video/towers.html")
		if err != nil {
			http.Error(w, "reading file towers.html", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(data))
	}
}
func TheodenHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("error connecting to db", err)
		return
	}
	defer db.Close()
	valid, err := database.CheckUuid(db, cookie.Value)
	if err != nil {
		fmt.Println("error retrieving uuid from db:", err)
		return
	}
	if !valid {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	} else {
		data, err := os.ReadFile("html/video/theoden.html")
		if err != nil {
			http.Error(w, "error reading file theoden.html", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(data))
	}
}

func TowerHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("error connecting to db", err)
		return
	}
	defer db.Close()
	valid, err := database.CheckUuid(db, cookie.Value)
	if err != nil {
		fmt.Println("error retrieving uuid from db:", err)
		return
	}
	if !valid {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	} else {
		data, err := os.ReadFile("html/video/towers.html")
		if err != nil {
			http.Error(w, "reading file towers.html", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(data))
	}
}
func FellowshipHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("error connecting to db", err)
		return
	}
	defer db.Close()
	valid, err := database.CheckUuid(db, cookie.Value)
	if err != nil {
		fmt.Println("error retrieving uuid from db:", err)
		return
	}
	if !valid {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	} else {
		data, err := os.ReadFile("html/video/fellowship.html")
		if err != nil {
			http.Error(w, "error reading file followship.html", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(data))
	}
}
