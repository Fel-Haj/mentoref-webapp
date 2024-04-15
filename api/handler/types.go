package handler

import "html/template"

type PageData struct {
	Title         string
	UserName      string
	UserSurname   string
	UserMail      string
	Authenticated bool
}

var (
	Index     = template.Must(template.ParseFiles("web/template/base.html", "web/template/page/index.html"))
	Dashboard = template.Must(template.ParseFiles("web/template/base.html", "web/template/page/dashboard.html"))
	SignIn    = template.Must(template.ParseFiles("web/template/component/signin.html"))
	SignUp    = template.Must(template.ParseFiles("web/template/component/signup.html"))
)
