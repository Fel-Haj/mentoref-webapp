package auth

import (
	"net/http"
)

func SignOutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			cookie := &http.Cookie{
				Name:     "access_token",
				Value:    "",
				MaxAge:   -1,
				HttpOnly: true,
				// Secure:   true, // for production
				SameSite: http.SameSiteStrictMode,
				// Path:     "/",
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}
