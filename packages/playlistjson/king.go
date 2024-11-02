package playlistjson

import (
	"fmt"
	"net/http"
	"os"
)

func PasswordRight(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "authenticated",        // Name of the cookie
		Value:    "nLtRJ6AblhzlerzUBhKF", // Value of the cookie
		Path:     "/",                    // Path the cookie is valid for
		HttpOnly: true,                   // Prevent JavaScript access
		Secure:   true,                   // Set to true if using HTTPS
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)

	// Response data to send back
	http.Redirect(w, r, "/protected", http.StatusSeeOther)
}
func ProtectionHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("authenticated")

	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// Some other error occurred
		http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
		return
	}
	//	fmt.Println(cookie.Value)
	if cookie.Value != os.Getenv("cookie") {
		// If the cookie is missing or invalid, you can respond with an error or redirect
		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to a login page
		return
	} else {
		data, err := os.ReadFile("html/video/pickafterauth.html")
		if err != nil {
			fmt.Println("reading html/video/pickafterauth.html", err)
			http.Error(w, "Error readinghtml/video/pickafterauth.html", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, string(data))
	}

}
