package auth

import (
	"database/sql"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/handler"
	"mentoref-webapp/internal/middleware"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SignInHandler(dbClient *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := handler.SignIn.Execute(w, nil)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<div>Failed to parse form</div>`))
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")

			var user db.User
			query := `SELECT id, email, password FROM users WHERE email = $1`
			err = dbClient.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Printf("No user found with email: %s", email)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`<div>Invalid email or password.</div>`))
					return
				}
				log.Printf("Database error: %v", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				log.Printf("Password mismatch for user %s: %v", email, err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`<div>Invalid email or password.</div>`))
				return
			}

			token, err := middleware.GenerateToken(&user)
			if err != nil {
				log.Printf("Failed to generate token: %v", err)
				http.Error(w, "Failed to generate token", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    token,
				HttpOnly: true,
				Expires:  time.Now().Add(30 * time.Minute),
				// Secure:   true, // for production
				SameSite: http.SameSiteStrictMode,
				// Path:     "/",
			})
			w.Header().Set("HX-Redirect", "/dashboard")
			return
		}
	}
}
