package auth

import (
	"log"
	"net/http"
	"time"
)

func SignOutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				log.Fatalf("Error occured during Sign Out: %s", err)
			}
			cookie.Expires = time.Now()
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
