package auth

import (
	"context"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/types"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := types.SignUp.Execute(w, nil)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println("invalid form data: %w", err)
				return
			}

			var newUser types.User

			newUser.ID = primitive.NewObjectID()
			newUser.Email = r.FormValue("email")
			newUser.Password = r.FormValue("password")
			newUser.FirstName = r.FormValue("firstname")
			newUser.Surname = r.FormValue("surname")
			newUser.Phone = r.FormValue("phone")
			newUser.CompanyName = r.FormValue("company")

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
			if err != nil {
				http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			}
			newUser.Password = string(hashedPassword)

			usersCollection := db.UserCollection(client)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			_, err = usersCollection.InsertOne(ctx, newUser)
			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					log.Println("request timed out, please try again: %w", err)
				}
				http.Error(w, "Failed to save user", http.StatusInternalServerError)
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
