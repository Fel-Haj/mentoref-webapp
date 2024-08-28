package types

import (
	"html/template"
	"mentoref-webapp/web"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// handler
type PageData struct {
	Title         string
	UserName      string
	UserSurname   string
	UserMail      string
	Authenticated bool
}

var (
	templateFS = web.GetHTMLTemplateFS()
	Index      = template.Must(template.ParseFS(templateFS, "base.html", "page/index.html"))
	Dashboard  = template.Must(template.ParseFS(templateFS, "base.html", "page/dashboard.html"))
	SignIn     = template.Must(template.ParseFS(templateFS, "component/auth/signin.html"))
	SignUp     = template.Must(template.ParseFS(templateFS, "component/auth/signup.html"))
	BlankShot  = template.Must(template.ParseFS(templateFS, "component/index/feature/blank-shot.html"))
	Referral   = template.Must(template.ParseFS(templateFS, "component/index/feature/referral.html"))
	Mentorship = template.Must(template.ParseFS(templateFS, "component/index/feature/mentorship.html"))
)

// middleware
type BoolContextKey bool
type UserMailContextKey string

var (
	AuthContextKey BoolContextKey
	UserContextKey UserMailContextKey
)

const (
	UserMailKey string = "userMail"
)

// db
type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	FirstName   string             `bson:"first_name"`
	Surname     string             `bson:"surname"`
	Phone       string             `bson:"phone,omitempty"`
	CompanyName string             `bson:"company_name,omitempty"`
}
