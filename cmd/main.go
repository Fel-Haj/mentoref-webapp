package main

import (
	"context"
	"io/fs"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/handler/auth"
	"mentoref-webapp/internal/handler/page/dashboard/user"
	"mentoref-webapp/internal/handler/page/index"
	"mentoref-webapp/internal/handler/page/index/index_component"
	"mentoref-webapp/internal/middleware"
	"mentoref-webapp/web"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	client := db.ConnectDB(os.Getenv("MONGODB_URI"))
	defer client.Disconnect(context.Background())

	var StaticFiles = web.StaticFiles

	staticFileSystem, err := fs.Sub(StaticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	staticFileServer := http.FileServer(http.FS(staticFileSystem))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	// page
	mux.HandleFunc("/", middleware.JWTAuthMiddleware(index.IndexHandler()))
	mux.HandleFunc("/dashboard", middleware.JWTAuthMiddleware(user.DashboardHandler(client)))

	// auth
	mux.HandleFunc("/signin", auth.SignInHandler(client))
	mux.HandleFunc("/signup", auth.SignUpHandler(client))
	mux.HandleFunc("/signout", auth.SignOutHandler())

	// page/index/component
	mux.HandleFunc("/blank-shot", index_component.BlankShotHandler())
	mux.HandleFunc("/referral", index_component.ReferralHandler())
	mux.HandleFunc("/mentorship", index_component.MentorshipHandler())

	log.Fatal(http.ListenAndServeTLS(":443", os.Getenv("CERTIFICATE"), os.Getenv("PRIVATE_KEY"), mux))
}
