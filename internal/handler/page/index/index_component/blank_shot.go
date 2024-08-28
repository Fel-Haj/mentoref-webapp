package index_component

import (
	"mentoref-webapp/internal/types"
	"net/http"
)

func BlankShotHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := types.BlankShot.Execute(w, nil)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		}
	}
}
