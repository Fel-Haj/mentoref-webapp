package pages

import (
	"fmt"
	"mentoref-webapp/internal/handler"
	"mentoref-webapp/internal/middleware"
	"net/http"
)

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, ok := r.Context().Value(middleware.AuthContextKey).(bool)
		if !ok {
			authenticated = false
		}

		data := handler.PageData{
			Title:         "MentoRef",
			Authenticated: authenticated,
		}

		err := handler.Index.Execute(w, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
