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
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Printf("Form parse error: %v", err)
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")

			log.Printf("Attempting login for email: %s", email)

			if dbClient == nil {
				log.Printf("Database client is nil")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			var user db.User
			query := `SELECT id, email, password FROM users WHERE email = $1`
			err = dbClient.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Printf("No user found with email: %s", email)
					http.Error(w, "Invalid email or password", http.StatusUnauthorized)
					return
				}
				log.Printf("Database error: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				log.Printf("Password mismatch for user %s: %v", email, err)
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}

			token, err := middleware.GenerateToken(&user)
			if err != nil {
				log.Printf("Failed to generate token: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    token,
				HttpOnly: true,
				Expires:  time.Now().Add(30 * time.Minute),
				// Secure:   true, // for production
				SameSite: http.SameSiteStrictMode,
				// Path:     "/",
			})
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
}
