package main

import (
	"io/fs"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/handler/components/auth"
	"mentoref-webapp/internal/handler/components/features"
	"mentoref-webapp/internal/handler/components/menu"
	"mentoref-webapp/internal/handler/pages"
	"mentoref-webapp/web"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	dbClient := db.ConnectDB()
	defer dbClient.Close()

	var StaticFiles = web.StaticFiles

	staticFileSystem, err := fs.Sub(StaticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	staticFileServer := http.FileServer(http.FS(staticFileSystem))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	// pages
	mux.HandleFunc("/", pages.IndexHandler())
	mux.HandleFunc("/dashboard", pages.DashboardHandler(dbClient))

	// auth
	mux.HandleFunc("/signin", auth.SignInHandler(dbClient))
	mux.HandleFunc("/signup", auth.SignUpHandler(dbClient))
	mux.HandleFunc("/signout", auth.SignOutHandler())

	// components
	mux.HandleFunc("/feature", features.FeatureBlockHandler())
	mux.HandleFunc("/menu", menu.MenuHandler(dbClient))
	mux.HandleFunc("/options", menu.OptionsHandler())
	mux.HandleFunc("/notification", menu.MenuHandler(dbClient))

	log.Fatal(http.ListenAndServeTLS(":443", os.Getenv("CERTIFICATE"), os.Getenv("PRIVATE_KEY"), mux))
}
