package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed static
var StaticFiles embed.FS

//go:embed templates
var HTMLTemplates embed.FS

func GetHTMLTemplateFS() fs.FS {
	HTMLTemplatesFileSystem, err := fs.Sub(HTMLTemplates, "templates")
	if err != nil {
		log.Fatalf("Error getting HTML templates FS: %s", err)
	}
	return HTMLTemplatesFileSystem
}
