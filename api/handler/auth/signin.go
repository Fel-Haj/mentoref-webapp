package auth

import (
	"context"
	"fmt"
	"log"
	"mentoref-webapp/api/handler"
	"mentoref-webapp/api/middleware"
	"mentoref-webapp/db"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignInHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := handler.SignIn.Execute(w, nil)
			if err != nil {
				return
			}
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")

			coll := db.UserCollection(client)
			user, err := AuthenticateUser(r.Context(), coll, email, password)
			if err != nil {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			}

			token, err := middleware.GenerateToken(user)
			if err != nil {
				log.Println("failed to generate token: %w", err)
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

func AuthenticateUser(ctx context.Context, coll *mongo.Collection, email string, password string) (*db.User, error) {
	var user db.User
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Error during PW comparison ", err)
		return nil, err
	}

	return &user, nil
}
