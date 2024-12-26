package pages

import (
	"fmt"
	"mentoref-webapp/internal/handler"
	"net/http"
)

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session")

		data := handler.PageData{
			Title:         "MentoRef",
			Authenticated: session != nil,
		}

		err = handler.Index.Execute(w, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
