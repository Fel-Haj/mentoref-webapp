package menu

import (
	"net/http"
)

func OptionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		optionType := r.URL.Query().Get("type")
		optionChoice := r.URL.Query().Get("choice")
		if r.Method == "GET" {
			switch optionType {
			case "continent":
				switch optionChoice {
				case "europe":
					print("test")
				}
			case "country":
				switch optionChoice {
				case "spain":
					print("test")
				}
			}
		}
	}
}
