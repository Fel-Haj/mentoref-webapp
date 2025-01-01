package web

import (
	"embed"
	"io/fs"
	"log"
	"text/template"
)

//go:embed static
var StaticFiles embed.FS

//go:embed templates
var HTMLTemplates embed.FS

var (
	templateFS       = GetHTMLTemplateFS()
	Index            = template.Must(template.ParseFS(templateFS, "base.html", "pages/index.html"))
	UserDashboard    = template.Must(template.ParseFS(templateFS, "base.html", "pages/dashboard_user.html"))
	CompanyDashboard = template.Must(template.ParseFS(templateFS, "base.html", "pages/dashboard_company.html"))
	SignIn           = template.Must(template.ParseFS(templateFS, "components/auth/signin.html"))
	SignUp           = template.Must(template.ParseFS(templateFS, "components/auth/signup.html"))
	FeatureBlock     = template.Must(template.ParseFS(templateFS, "components/features/feature_block.html"))
	Menu             = template.Must(template.ParseFS(templateFS, "components/features/menu.html"))
	Notification     = template.Must(template.ParseFS(templateFS, "components/features/notification.html"))
)

func GetHTMLTemplateFS() fs.FS {
	HTMLTemplatesFileSystem, err := fs.Sub(HTMLTemplates, "templates")
	if err != nil {
		log.Fatalf("Error getting HTML templates FS: %s", err)
	}
	return HTMLTemplatesFileSystem
}
