package pages

import (
	"database/sql"
	"fmt"
	"log"
	"mentoref-webapp/db"
	"mentoref-webapp/internal/middleware"
	"mentoref-webapp/web"
	"net/http"
	"time"
)

type CompanyDashboardPageData struct {
	Title         string
	ContactMail   string
	Name          string
	Address       string
	Postcode      string
	Country       string
	Authenticated bool
}

type UserDashboardPageData struct {
	Title         string
	FirstName     string
	LastName      string
	Mail          string
	Phone         string
	Experience    []WorkExperience
	BlankShots    []BlankShot
	UserType      string
	Authenticated bool
}

type WorkExperience struct {
	ID              uint32
	StartDate       string
	EndDate         string
	JobFunction     string
	SeniorityLevel  string
	ManagementLevel string
}

type BlankShot struct {
	Date            string
	RemainingDays   uint8
	Continent       string
	Country         string
	State           string
	JobFunction     string
	Seniority       string
	Remote          bool
	VisaSponsorship bool
	UserID          uint32
}

func DashboardHandler(dbClient *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			cookie, err := r.Cookie("session")
			claims, err := middleware.GetClaims(cookie)
			if err != nil {
				fmt.Println(err)
			}
			accType := r.URL.Query().Get("type")
			if accType == "user" {
				userId := fmt.Sprintf("%.0f", claims["userId"])
				var user db.User
				err = dbClient.QueryRow("SELECT first_name, last_name, email, phone FROM users WHERE id = $1", userId).Scan(
					&user.FirstName,
					&user.LastName,
					&user.Email,
					&user.Phone,
				)
				if err != nil {
					fmt.Printf("Error getting user data: %v", err)
					return
				}

				var experience []WorkExperience
				rows, err := dbClient.Query("SELECT id, start_date, end_date, job_function, seniority_level, management_level FROM work_experience WHERE user_id = $1", userId)
				if err != nil {
					log.Printf("Error querying work experience: %v", err)
					return
				}
				defer rows.Close()

				for rows.Next() {
					var we WorkExperience
					err := rows.Scan(&we.ID, &we.StartDate, &we.EndDate, &we.JobFunction, &we.SeniorityLevel, &we.ManagementLevel)
					if err != nil {
						log.Printf("Error scanning work experience row: %v", err)
						continue
					}
					experience = append(experience, we)
				}

				var shots []BlankShot
				shotRows, err := dbClient.Query("SELECT DATE(date), continent, country, state, job_function, seniority, remote, visa_sponsorship FROM blank_shot WHERE user_id = $1", userId)
				if err != nil {
					log.Printf("Error querying blank shots: %v", err)
					return
				}
				defer shotRows.Close()

				for shotRows.Next() {
					var shot BlankShot
					var queryDate time.Time
					err := shotRows.Scan(&queryDate, &shot.Continent, &shot.Country, &shot.State, &shot.JobFunction, &shot.Seniority, &shot.Remote, &shot.VisaSponsorship)
					if err != nil {
						log.Printf("Error scanning blank shot row: %v", err)
						continue
					}
					expiryDate := queryDate.Add(14 * 24 * time.Hour)
					remaining := expiryDate.Sub(time.Now().Truncate(24 * time.Hour))
					shot.RemainingDays = uint8(remaining.Hours() / 24)
					shot.Date = queryDate.Format("02.01.2006")
					shots = append(shots, shot)
				}

				data := UserDashboardPageData{
					Title:         fmt.Sprintf("%s %s - Profile", user.FirstName, user.LastName),
					FirstName:     user.FirstName,
					LastName:      user.LastName,
					Mail:          user.Email,
					Phone:         user.Phone,
					Experience:    experience,
					BlankShots:    shots,
					UserType:      accType,
					Authenticated: true,
				}
				err = web.UserDashboard.Execute(w, data)
				if err != nil {
					log.Printf("Template execution error: %v", err)
					http.Error(w, "Error rendering template", http.StatusInternalServerError)
					return
				}
			} else if accType == "company" {
				userId := fmt.Sprintf("%.0f", claims["companyId"])
				var company db.Company
				err = dbClient.QueryRow("SELECT contact_email, company_name, company_address, company_postcode, company_country FROM companies WHERE id = $1", userId).Scan(
					&company.ContactEmail,
					&company.Name,
					&company.Address,
					&company.Postcode,
					&company.Country,
				)
				if err != nil {
					fmt.Printf("Error getting user data: %v", err)
					return
				}

				data := CompanyDashboardPageData{
					Title:         fmt.Sprintf("%s %s - Profile", company.Name, company.Country),
					ContactMail:   company.ContactEmail,
					Name:          company.Name,
					Address:       company.Address,
					Postcode:      company.Postcode,
					Country:       company.Country,
					Authenticated: true,
				}
				err = web.CompanyDashboard.Execute(w, data)
				if err != nil {
					log.Printf("Template execution error: %v", err)
					http.Error(w, "Error rendering template", http.StatusInternalServerError)
					return
				}
			}
		}
	}
}
