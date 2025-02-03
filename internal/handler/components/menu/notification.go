package menu

import (
	"mentoref-webapp/web"
	"net/http"
)

type NotificationComponentData struct {
	Legend     string
	Text       string
	HxEndpoint string
}

func NotificationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nType := r.URL.Query().Get("type")
		var data NotificationComponentData
		if r.Method == "GET" {
			if nType == "blank-shot" {
				data = NotificationComponentData{
					Legend:     "Blank-Shot fired.",
					Text:       "You have succesfully sent out a blank-shot.",
					HxEndpoint: "/",
				}
			}
			if nType == "blank-shot-error-signin" {
				data = NotificationComponentData{
					Legend:     "Not signed in.",
					Text:       "You have to be signed in in order to send out a blank-shot.<br>Sign-Up now if you haven't done so yet.",
					HxEndpoint: "/",
				}
			}
			err := web.Notification.Execute(w, data)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		}
	}
}
