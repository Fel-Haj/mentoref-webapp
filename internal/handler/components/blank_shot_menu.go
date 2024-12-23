package components

import (
	"mentoref-webapp/internal/handler"
	"net/http"
)

func BlankShotMenuHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := handler.BlankShotMenu.Execute(w, nil)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		}
	}
}
