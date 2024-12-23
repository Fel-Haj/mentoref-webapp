package handler

import (
	"html/template"
	"mentoref-webapp/web"
)

type PageData struct {
	Title         string
	UserName      string
	UserSurname   string
	UserMail      string
	UserPhone     string
	Authenticated bool
}

var (
	templateFS    = web.GetHTMLTemplateFS()
	Index         = template.Must(template.ParseFS(templateFS, "base.html", "pages/index.html"))
	Dashboard     = template.Must(template.ParseFS(templateFS, "base.html", "pages/dashboard.html"))
	SignIn        = template.Must(template.ParseFS(templateFS, "components/auth/signin.html"))
	SignUp        = template.Must(template.ParseFS(templateFS, "components/auth/signup.html"))
	FeatureBlock  = template.Must(template.ParseFS(templateFS, "components/features/feature_block.html"))
	BlankShotMenu = template.Must(template.ParseFS(templateFS, "components/features/blank-shot-menu.html"))
)
