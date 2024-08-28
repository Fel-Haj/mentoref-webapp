package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed static
var StaticFiles embed.FS

//go:embed template
var HTMLTemplates embed.FS

func GetHTMLTemplateFS() fs.FS {
	HTMLTemplatesFileSystem, err := fs.Sub(HTMLTemplates, "template")
	if err != nil {
		log.Fatalf("Error getting HTML templates FS: %s", err)
	}
	return HTMLTemplatesFileSystem
}
