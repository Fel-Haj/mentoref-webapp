package user

import (
	"fmt"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/types"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DashboardHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userMail, ok := r.Context().Value(types.UserContextKey).(string)
		if !ok {
			http.Error(w, "Invalid userMail in context", http.StatusInternalServerError)
			return
		}

		coll := db.UserCollection(client)
		var user types.User
		err := coll.FindOne(r.Context(), bson.M{"email": userMail}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return
			}
			panic(err)
		}

		authenticated, ok := r.Context().Value(types.AuthContextKey).(bool)
		if !ok {
			authenticated = false
		}

		data := types.PageData{
			Title:         fmt.Sprintf("%s %s - Profile", user.FirstName, user.Surname),
			UserName:      user.FirstName,
			UserSurname:   user.Surname,
			UserMail:      user.Email,
			Authenticated: authenticated,
		}

		err = types.Dashboard.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

	}
}
