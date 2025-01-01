package auth

import (
	"context"
	"database/sql"
	"html/template"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/web"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type SignUpPageData struct {
	HxEndpoint     string
	Fields         template.HTML
	ButtonText     string
	LinkText       string
	HxLinkEndpoint string
}

func SignUpHandler(dbClient *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.RequestURI()

		if r.Method == "GET" {
			var data SignUpPageData
			var err error
			if strings.Contains(path, "company") {
				data = SignUpPageData{
					HxEndpoint: path,
					Fields: template.HTML(`<input type="email" name="email" id="email" required placeholder="Contact Email" tabindex="1"
						class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />
						<input type="password" name="password" id="password" required placeholder="Password" tabindex="2"
						class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />
						<input type="text" name="company_name" id="company_name" required placeholder="Company Name" tabindex="3"
						class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />
						<input type="text" name="company_address" id="company_address" required placeholder="Company Address" tabindex="4"
						class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />
						<input type="text" name="company_postcode" id="company_postcode" required placeholder="Company Postcode" tabindex="5"
						class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />
						<input type="text" name="company_country" id="company_country" required placeholder="Company Country" tabindex="6"
						class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />`),
					ButtonText:     "Sign Up as Company",
					LinkText:       "Sign Up as User",
					HxLinkEndpoint: "/signup?type=user",
				}
			} else {
				data = SignUpPageData{
					HxEndpoint: path,
					Fields: template.HTML(`<input type="email" name="email" id="email" required placeholder="E-Mail" tabindex="1"
							class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />
							<input type="password" name="password" id="password" required placeholder="Password" tabindex="2"
							class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />
							<input type="text" name="first_name" id="first_name" required placeholder="First Name" tabindex="3"
							class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />
							<input type="text" name="last_name" id="last_name" required placeholder="Last Name" tabindex="4"
							class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />
							<input type="tel" name="phone" id="phone" placeholder="Phonenumber" tabindex="5"
							class="input-text focus:outline-none focus:ring-2 focus:ring-blue-500" />`),
					ButtonText:     "Sign Up as User",
					LinkText:       "Sign Up as Company",
					HxLinkEndpoint: "/signup?type=company",
				}
			}
			err = web.SignUp.Execute(w, data)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
		}

		if r.Method == "POST" {
			w.Header().Set("Content-Type", "text/html")

			err := r.ParseForm()
			if err != nil {
				log.Printf("Invalid form data: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<div>Invalid form data</div>`))
			}

			if strings.Contains(path, "company") {
				var newCompany db.Company
				newCompany.ContactEmail = r.FormValue("email")
				newCompany.Password = r.FormValue("password")
				newCompany.Name = r.FormValue("company_name")
				newCompany.Address = r.FormValue("company_address")
				newCompany.Postcode = r.FormValue("company_postcode")
				newCompany.Country = r.FormValue("company_country")

				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newCompany.Password), 14)
				if err != nil {
					log.Printf("Password hashing error: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`<div>Server error occurred</div>`))
				}
				newCompany.Password = string(hashedPassword)

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				_, err = dbClient.ExecContext(ctx,
					`INSERT INTO companies (contact_email, password, company_name, company_address, company_postcode, company_country)
					VALUES ($1, $2, $3, $4, $5, $6)`,
					newCompany.ContactEmail, newCompany.Password, newCompany.Name, newCompany.Address, newCompany.Postcode, newCompany.Country)

				if err != nil {
					log.Printf("Database error: %v", err)

					if ctx.Err() == context.DeadlineExceeded {
						w.WriteHeader(http.StatusRequestTimeout)
						w.Write([]byte(`<div>Request timed out. Please try again</div>`))
					}

					if strings.Contains(err.Error(), "unique constraint") ||
						strings.Contains(err.Error(), "duplicate key") {
						w.WriteHeader(http.StatusConflict)
						w.Write([]byte(`<div>Email already exists</div>`))
					}

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`<div>Failed to create account. Please try again</div>`))
				}
				SignIn(w, dbClient, newCompany.ContactEmail, newCompany.Password)
			} else {
				var newUser db.User
				newUser.Email = r.FormValue("email")
				newUser.Password = r.FormValue("password")
				newUser.FirstName = r.FormValue("first_name")
				newUser.LastName = r.FormValue("last_name")
				newUser.Phone = r.FormValue("phone")

				if newUser.Email == "" || newUser.Password == "" || newUser.FirstName == "" || newUser.LastName == "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`<div>All required fields must be filled</div>`))
				}

				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
				if err != nil {
					log.Printf("Password hashing error: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`<div>Server error occurred</div>`))
				}
				hashedUserPassword := string(hashedPassword)

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				_, err = dbClient.ExecContext(ctx,
					`INSERT INTO users (email, password, first_name, last_name, phone)
					VALUES ($1, $2, $3, $4, $5)`,
					newUser.Email, hashedUserPassword, newUser.FirstName, newUser.LastName, newUser.Phone)

				if err != nil {
					log.Printf("Database error: %v", err)

					if ctx.Err() == context.DeadlineExceeded {
						w.WriteHeader(http.StatusRequestTimeout)
						w.Write([]byte(`<div>Request timed out.<br>Please try again</div>`))
					}

					if strings.Contains(err.Error(), "unique constraint") {
						w.WriteHeader(http.StatusConflict)
						w.Write([]byte(`<div>Email already exists.</div>`))
					}
				}
				SignIn(w, dbClient, newUser.Email, newUser.Password)
			}
		}
	}
}
