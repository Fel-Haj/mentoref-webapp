package features

import (
	"mentoref-webapp/web"
	"net/http"
)

type FeatureData struct {
	MainHeader   string
	SubHeader    string
	Description  string
	ButtonAction string
	VideoSrc     string
	HxEndpoint   string
}

var featureConfigs = map[string]FeatureData{
	"blank-shot": {
		SubHeader:    "Boost your career with one click!",
		MainHeader:   "Blank Shot\nThrow yourself out there",
		Description:  "Take a leap into new opportunities. Our Blank Shot feature connects you with potential employers and collaborators, showcasing your unique talents and aspirations.",
		ButtonAction: "Shoot!",
		VideoSrc:     "/static/image/animation/MentoRef_Anim.mov",
		HxEndpoint:   "/menu?type=blank-shot",
	},
	"mentorship": {
		SubHeader:    "Navigate the job market with a seasoned mentor!",
		MainHeader:   "Mentorship\nGaining insights and advice",
		Description:  "This is a placeholder for our mentorship feature. It offers a range of benefits and functionalities that can be customized to meet your specific needs. Explore the possibilities and see how it can improve your workflow.",
		ButtonAction: "Try Now",
		VideoSrc:     "/static/image/animation/MentoRef_Anim.mov",
		HxEndpoint:   "/mentorship-menu",
	},
	"referral": {
		SubHeader:    "Leverage the power of connections for your next career move!",
		MainHeader:   "Referral\nNetworking pays off",
		Description:  "This is a placeholder for our referral feature. It offers a range of benefits and functionalities that can be customised to meet your specific needs. Explore the possibilities and see how it can improve your workflow.",
		ButtonAction: "Try Now",
		VideoSrc:     "/static/image/animation/MentoRef_Anim.mov",
		HxEndpoint:   "/referral-menu",
	},
}

func FeatureBlockHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			featureType := r.URL.Query().Get("type")

			data, exists := featureConfigs[featureType]
			if !exists {
				http.Error(w, "Feature not found", http.StatusNotFound)
				return
			}

			err := web.FeatureBlock.Execute(w, data)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		}
	}
}
