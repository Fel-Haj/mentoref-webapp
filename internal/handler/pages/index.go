package pages

import (
	"fmt"
	"mentoref-webapp/internal/middleware"
	"mentoref-webapp/web"
	"net/http"
)

type IndexPageData struct {
	Title         string
	AccType       string
	Authenticated bool
}

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accType string
		var authenticated bool
		accType, err := middleware.GetAccountType(r)
		if err == nil {
			authenticated = true
		}

		data := IndexPageData{
			Title:         "MentoRef",
			AccType:       accType,
			Authenticated: authenticated,
		}

		err = web.Index.Execute(w, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
