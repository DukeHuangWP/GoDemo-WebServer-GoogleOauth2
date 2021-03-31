package googleOauth2

import (
	"golang.org/x/oauth2"
	googleInc "golang.org/x/oauth2/google"
)

var (
	OauthConfig *oauth2.Config
)

//Google帳號資訊
type GoogleAcc struct {
	ID            int
	Email         string
	VerifiedEmail bool
	PictureUrl    string
}

func InitConfig(redirectURL, clientID, clientSecret string) {
	//https://console.cloud.google.com/apis/credentials
	OauthConfig = &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret, //set from google credentials
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     googleInc.Endpoint,
	}

}
