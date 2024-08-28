package index

import (
	"fmt"
	"mentoref-webapp/internal/types"
	"net/http"
)

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, ok := r.Context().Value(types.AuthContextKey).(bool)
		if !ok {
			authenticated = false
		}

		data := types.PageData{
			Title:         "MentoRef",
			Authenticated: authenticated,
		}

		err := types.Index.Execute(w, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
