package auth

import (
	"context"
	"database/sql"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/handler"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(dbClient *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := handler.SignUp.Execute(w, nil)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		}

		if r.Method == "POST" {
			w.Header().Set("Content-Type", "text/html")

			err := r.ParseForm()
			if err != nil {
				log.Printf("Invalid form data: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<div>Invalid form data</div>`))
				return
			}

			var newUser db.User
			newUser.Email = r.FormValue("email")
			newUser.Password = r.FormValue("password")
			newUser.FirstName = r.FormValue("firstname")
			newUser.LastName = r.FormValue("lastname")
			newUser.Phone = r.FormValue("phone")

			if newUser.Email == "" || newUser.Password == "" || newUser.FirstName == "" || newUser.LastName == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<div>All required fields must be filled</div>`))
				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
			if err != nil {
				log.Printf("Password hashing error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`<div>Server error occurred</div>`))
				return
			}
			newUser.Password = string(hashedPassword)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			_, err = dbClient.ExecContext(ctx,
				`INSERT INTO users (email, password, first_name, last_name, phone)
                VALUES ($1, $2, $3, $4, $5)`,
				newUser.Email, newUser.Password, newUser.FirstName, newUser.LastName, newUser.Phone)

			if err != nil {
				log.Printf("Database error: %v", err)

				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusRequestTimeout)
					w.Write([]byte(`<div>Request timed out. Please try again</div>`))
					return
				}

				if strings.Contains(err.Error(), "unique constraint") ||
					strings.Contains(err.Error(), "duplicate key") {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(`<div>Email already exists</div>`))
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`<div>Failed to create account. Please try again</div>`))
				return
			}

			// Successful signup
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
