package pages

import (
	"fmt"
	"mentoref-webapp/internal/middleware"
	"mentoref-webapp/web"
	"net/http"
)

type IndexPageData struct {
	Title         string
	UserType      string
	Authenticated bool
}

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session")
		var userType string
		var authenticated bool

		if err == nil {
			claims, err := middleware.GetClaims(session)
			if err == nil {
				if claims["userId"] != nil {
					userType = "user"
				} else if claims["companyId"] != nil {
					userType = "company"
				}
				authenticated = true
			}
		}

		data := IndexPageData{
			Title:         "MentoRef",
			UserType:      userType,
			Authenticated: authenticated,
		}

		err = web.Index.Execute(w, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
