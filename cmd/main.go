package main

import (
	"context"
	"io/fs"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/handler/auth"
	"mentoref-webapp/internal/handler/components"
	"mentoref-webapp/internal/handler/pages"
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
	mux.HandleFunc("/", middleware.JWTAuthMiddleware(pages.IndexHandler()))
	mux.HandleFunc("/dashboard", middleware.JWTAuthMiddleware(pages.DashboardHandler(client)))

	// auth
	mux.HandleFunc("/signin", auth.SignInHandler(client))
	mux.HandleFunc("/signup", auth.SignUpHandler(client))
	mux.HandleFunc("/signout", auth.SignOutHandler())

	// components
	mux.HandleFunc("/blank-shot", components.FeatureBlockHandler())
	mux.HandleFunc("/referral", components.FeatureBlockHandler())
	mux.HandleFunc("/mentorship", components.FeatureBlockHandler())
	mux.HandleFunc("/blank-shot-menu", components.BlankShotMenuHandler())

	log.Fatal(http.ListenAndServeTLS(":443", os.Getenv("CERTIFICATE"), os.Getenv("PRIVATE_KEY"), mux))
}
