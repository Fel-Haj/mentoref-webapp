package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"mentoref-webapp/api/handler/auth"
	"mentoref-webapp/api/handler/page"
	"mentoref-webapp/api/middleware"
	"mentoref-webapp/db"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

//go:embed web/static
var embeddedFiles embed.FS

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	mux := http.NewServeMux()
	client := db.ConnectDB(os.Getenv("MONGODB_URI"))
	defer client.Disconnect(context.Background())

	staticFiles, err := fs.Sub(embeddedFiles, "web/static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(staticFiles))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", middleware.JWTAuthMiddleware(page.IndexHandler()))
	mux.HandleFunc("/dashboard", middleware.JWTAuthMiddleware(page.DashboardHandler(client)))
	mux.HandleFunc("/signin", auth.SignInHandler(client))
	mux.HandleFunc("/signup", auth.SignUpHandler(client))
	mux.HandleFunc("/signout", auth.SignOutHandler())

	log.Fatal(http.ListenAndServe(":3000", mux))
}
