package auth

import (
	"database/sql"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/middleware"
	"mentoref-webapp/web"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func checkPassword(w http.ResponseWriter, storedPassword string, email string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		log.Printf("Password mismatch for email %s.\n%v", email, err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<div>Invalid email or password.</div>`))
	}
	return err
}

func SignIn(w http.ResponseWriter, dbClient *sql.DB, email string, password string) {
	var token string

	var user db.User
	err := dbClient.QueryRow(`SELECT id, email, password FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err == nil {
		err = checkPassword(w, user.Password, user.Email, password)
		if err != nil {
			return
		}
		token = middleware.GenerateToken(&user)
	} else {
		var company db.Company
		err = dbClient.QueryRow(`SELECT id, contact_email, password FROM companies WHERE contact_email = $1`, email).
			Scan(&company.ID, &company.ContactEmail, &company.Password)
		if err == sql.ErrNoRows {
			log.Printf("No user or company found with email: %s", email)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`<div>Invalid email or password.</div>`))
			return
		} else if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		err = checkPassword(w, company.Password, company.ContactEmail, password)
		if err != nil {
			return
		}
		token = middleware.GenerateToken(&company)
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
	// w.Header().Set("HX-Redirect", "/dashboard")
	return
}

func SignInHandler(dbClient *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := web.SignIn.Execute(w, nil)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<div>Failed to parse form</div>`))
			}

			email := r.FormValue("email")
			password := r.FormValue("password")

			SignIn(w, dbClient, email, password)
		}
	}
}
