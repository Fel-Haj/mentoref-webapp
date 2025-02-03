package menu

import (
	"context"
	"database/sql"
	"fmt"
	"mentoref-webapp/internal/middleware"
	"mentoref-webapp/web"
	"net/http"
	"time"
)

type Option struct {
	Value string
	Name  string
}

type Label struct {
	Title   string
	Prop    string
	Options []Option
}

type MenuPageData struct {
	HxEndpointType string
	Legend         string
	Labels         []Label
	SubmitText     string
}

func MenuHandler(dbClient *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		featType := r.URL.Query().Get("type")
		var data MenuPageData
		if r.Method == "GET" {
			if featType == "blank-shot" {
				data = MenuPageData{
					HxEndpointType: "blank-shot",
					Legend:         "Select Filters for your Blank-Shot",
					Labels: []Label{
						{
							Prop: "continent",
							Options: []Option{
								{
									Value: "",
									Name:  "Select Continent",
								},
								{
									Value: "europe",
									Name:  "Europe",
								},
								{
									Value: "north_america",
									Name:  "North America",
								},
							},
						},
						{
							Prop: "country",
							Options: []Option{
								{
									Value: "",
									Name:  "Select Country",
								},
								{
									Value: "spain",
									Name:  "Spain",
								},
								{
									Value: "usa",
									Name:  "USA",
								},
								{
									Value: "canada",
									Name:  "Canada",
								},
							},
						},
						{
							Prop: "state",
							Options: []Option{
								{
									Value: "",
									Name:  "Select State",
								},
								{
									Value: "andalusia",
									Name:  "Andalusia",
								},
								{
									Value: "catalonia",
									Name:  "Catalonia",
								},
								{
									Value: "madrid",
									Name:  "Madrid",
								},
								{
									Value: "valencia",
									Name:  "Valencia",
								},
								{
									Value: "california",
									Name:  "California",
								},
								{
									Value: "texas",
									Name:  "Texas",
								},
							},
						},
						{
							Prop: "job-function",
							Options: []Option{
								{
									Value: "",
									Name:  "Select Job Function",
								},
								{
									Value: "consulting",
									Name:  "Consulting",
								},
								{
									Value: "marketing",
									Name:  "Marketing",
								},
								{
									Value: "sales",
									Name:  "Sales",
								},
								{
									Value: "engineering",
									Name:  "Engineering",
								},
							},
						},
						{
							Prop: "seniority",
							Options: []Option{
								{
									Value: "",
									Name:  "Select Seniority",
								},
								{
									Value: "junior",
									Name:  "Junior",
								},
								{
									Value: "mid",
									Name:  "Mid",
								},
								{
									Value: "senior",
									Name:  "Senior",
								},
								{
									Value: "lead",
									Name:  "Lead",
								},
								{
									Value: "director",
									Name:  "Director",
								},
							},
						},
						{
							Prop: "industry",
							Options: []Option{
								{
									Value: "",
									Name:  "Select Industry",
								},
								{
									Value: "technology",
									Name:  "Technology",
								},
								{
									Value: "finance",
									Name:  "Finance",
								},
								{
									Value: "healthcare",
									Name:  "Healthcare",
								},
								{
									Value: "education",
									Name:  "Education",
								},
							},
						},
						{
							Title: "Remote?",
							Prop:  "remote",
						},
						{
							Title: "Visa Sponsorship?",
							Prop:  "visa-sponsorship",
						},
					},
					SubmitText: "Apply Filters & Shoot!",
				}
			}

			err := web.Menu.Execute(w, data)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		}
		if r.Method == "POST" {
			userId, err := middleware.GetAccountID(r)
			if err != nil {
				fmt.Println("Error getting Account ID")
			}

			continent := r.FormValue("continent")
			continentMap := map[string]string{
				"europe":        "Europe",
				"north_america": "North America",
			}
			continentDisplayName, exists := continentMap[continent]
			if !exists {
				continentDisplayName = "Unknown Continent"
			}

			country := r.FormValue("country")
			countryMap := map[string]string{
				"spain":  "Spain",
				"usa":    "USA",
				"canada": "Canada",
			}
			countryDisplayName, exists := countryMap[country]
			if !exists {
				countryDisplayName = "Unknown Country"
			}

			state := r.FormValue("state")
			stateMap := map[string]string{
				"andalusia":  "Andalusia",
				"catalonia":  "Catalonia",
				"madrid":     "Madrid",
				"valencia":   "Valencia",
				"california": "California",
				"texas":      "Texas",
			}
			stateDisplayName, exists := stateMap[state]
			if !exists {
				stateDisplayName = "Unknown State"
			}

			jobFunction := r.FormValue("job-function")
			jobFunctionMap := map[string]string{
				"consulting":  "Consulting",
				"marketing":   "Marketing",
				"sales":       "Sales",
				"engineering": "Engineering",
			}
			jobFunctionDisplayName, exists := jobFunctionMap[jobFunction]
			if !exists {
				jobFunctionDisplayName = "Unknown Job Function"
			}

			seniority := r.FormValue("seniority")
			seniorityMap := map[string]string{
				"junior":   "Junior",
				"mid":      "Mid",
				"senior":   "Senior",
				"lead":     "Lead",
				"director": "Director",
			}
			seniorityDisplayName, exists := seniorityMap[seniority]
			if !exists {
				seniorityDisplayName = "Unknown Seniority"
			}

			industry := r.FormValue("industry")
			industryMap := map[string]string{
				"technology": "Technology",
				"finance":    "Finance",
				"healthcare": "Healthcare",
				"education":  "Education",
			}
			industryDisplayName, exists := industryMap[industry]
			if !exists {
				industryDisplayName = "Unknown Industry"
			}
			remote := r.FormValue("remote")
			var remoteSQL bool
			if remote == "true" {
				remoteSQL = true
			} else {
				remoteSQL = false
			}
			visaSponsorship := r.FormValue("visa-sponsorship")
			var visaSponsorshipSQL bool
			if visaSponsorship == "true" {
				visaSponsorshipSQL = true
			} else {
				visaSponsorshipSQL = false
			}

			todaysDate := time.Now().Format("2006-01-02")

			fmt.Printf("Inserting values: date=%s, continent=%s, country=%s, state=%s, job=%s, seniority=%s, industry=%s, remote=%v, visa=%v, userId=%s\n",
				todaysDate, continentDisplayName, countryDisplayName, stateDisplayName, jobFunctionDisplayName, seniorityDisplayName, industryDisplayName, remoteSQL, visaSponsorshipSQL, userId)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			_, err = dbClient.ExecContext(ctx,
				`INSERT INTO blank_shot (date, continent, country, state, job_function, seniority, industry, remote, visa_sponsorship, user_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
				todaysDate, continentDisplayName, countryDisplayName, stateDisplayName, jobFunctionDisplayName, seniorityDisplayName, industryDisplayName, remoteSQL, visaSponsorshipSQL, userId)
			if err != nil {
				fmt.Println("Database error:", err) // Add this line
				http.Error(w, "Error inserting data", http.StatusInternalServerError)
				return
			}
		}
	}
}
