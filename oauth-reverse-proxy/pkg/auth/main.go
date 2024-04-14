package auth

import (
	"golang.org/x/oauth2"
)

type OAuth interface {
	GetAuthCodeURL(state string) string
	GetAccessToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
}

type UserInfo struct {
	Id    string
	Email string
}
