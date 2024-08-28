package index_component

import (
	"mentoref-webapp/internal/types"
	"net/http"
)

func ReferralHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := types.Referral.Execute(w, nil)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		}
	}
}
