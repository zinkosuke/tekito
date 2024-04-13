package auth

import (
	"context"
	"encoding/json"
	"log/slog"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	conf *oauth2.Config
}

func NewGoogleOAuth(redirectURL string, clientId string, clientSecret string) OAuth {
	conf := &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}
	return &GoogleOAuth{conf}
}

func (oa *GoogleOAuth) GetAuthCodeURL(state string) string {
	url := oa.conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	slog.Info("GetAuthCodeURL: Generated.", "url", url)
	return url
}

func (oa *GoogleOAuth) GetAccessToken(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	token, err := oa.conf.Exchange(ctx, code)
	if err != nil {
		slog.Error("GetAccessToken: Request failed.", err)
		return nil, err
	}
	return token, err
}

func (oa *GoogleOAuth) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	ctx := context.Background()
	client := oa.conf.Client(ctx, token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		slog.Error("GetUserInfo: Request failed.", err)
		return nil, err
	}
	rawUserInfo := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&rawUserInfo)
	if err != nil {
		slog.Error("GetUserInfo: Failed to decode.", err)
		return nil, err
	}
	userInfo := &UserInfo{
		Id:    rawUserInfo["id"].(string),
		Email: rawUserInfo["email"].(string),
	}
	return userInfo, nil
}
