package pages

import (
	"database/sql"
	"fmt"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/handler"
	"mentoref-webapp/internal/middleware"
	"net/http"
)

func DashboardHandler(dbClient *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		claims, err := middleware.GetClaims(cookie)
		if err != nil {
			fmt.Println(err)
		}

		userId := fmt.Sprintf("%.0f", claims["userId"])
		var user db.User
		err = dbClient.QueryRow("SELECT first_name, last_name, email, phone FROM users WHERE id = $1", userId).Scan(
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return
			}
			fmt.Printf("Error getting user data: %v", err)
			return
		}

		data := handler.PageData{
			Title:         fmt.Sprintf("%s %s - Profile", user.FirstName, user.LastName),
			UserFirstName: user.FirstName,
			UserLastName:  user.LastName,
			UserMail:      user.Email,
			UserPhone:     user.Phone,
			Authenticated: true,
		}

		err = handler.Dashboard.Execute(w, data)
		if err != nil {
			log.Printf("Template execution error: %v", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}
}
