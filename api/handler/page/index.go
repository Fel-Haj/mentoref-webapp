package page

import (
	"mentoref-webapp/api/handler"
	"mentoref-webapp/api/middleware"
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

		handler.Index.Execute(w, data)
	}
}
