package pages

import (
	"fmt"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/handler"
	"mentoref-webapp/internal/middleware"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DashboardHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userMail, ok := r.Context().Value(middleware.UserContextKey).(string)
		if !ok {
			http.Error(w, "Invalid userMail in context", http.StatusInternalServerError)
			return
		}

		coll := db.UserCollection(client)
		var user db.User
		err := coll.FindOne(r.Context(), bson.M{"email": userMail}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return
			}
			fmt.Printf("Error getting user data: %v", err)
			return
		}

		authenticated, ok := r.Context().Value(middleware.AuthContextKey).(bool)
		if !ok {
			authenticated = false
		}

		data := handler.PageData{
			Title:         fmt.Sprintf("%s %s - Profile", user.FirstName, user.Surname),
			UserName:      user.FirstName,
			UserSurname:   user.Surname,
			UserMail:      user.Email,
			UserPhone:     user.Phone,
			Authenticated: authenticated,
		}

		err = handler.Dashboard.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

	}
}
