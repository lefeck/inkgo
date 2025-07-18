package oauth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"inkgo/config"
	"net/http"
)

type GoogleProvider struct {
	Client *http.Client
	config *config.OAuthConfig
}

type GoogleAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

type GoogleUserInfo struct {
	UnionID     string `json:"unionid"`
	Url         string `json:"url"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
}

func (google *GoogleProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	client := google.Client
	if client == nil {
		client = &http.Client{}
	}

	userInfoURL := "https://www.googleapis.com/oauth2/v3/userinfo"
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	var googleUserInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUserInfo); err != nil {
		return nil, err
	}
	userInfo := UserInfo{
		UnionID:     googleUserInfo.UnionID,
		URL:         googleUserInfo.URL,
		Provider:    google.Name(),
		Username:    googleUserInfo.Username,
		DisplayName: googleUserInfo.DisplayName,
		Email:       googleUserInfo.Email,
		Avatar:      googleUserInfo.Avatar,
	}

	return &userInfo, nil

}

//
//func (g *GoogleProvider) GetAccessToken(ctx context.Context, code string) (*oauth2.Token, error) {
//	token, err := g.config.Config.Exchange(ctx, code)
//	if err != nil {
//		return nil, err
//	}
//
//}
//
//func NewGoogleProvider(config *config.OAuthConfig) OAuth2Provider {
//	return &GoogleProvider{
//		config: config,
//	}
//}
//
//func (g *GoogleProvider) GetAuthURL() string {
//	return "https://accounts.google.com/o/oauth2/auth" +
//		"?client_id=" + g.config.ClientID +
//		"&redirect_uri=" + g.config.RedirectURI +
//		"&response_type=code" +
//		"&scope=email%20profile" +
//		"&access_type=offline"
//}

func (g *GoogleProvider) GetAccessTokenURL() string {
	return "https://oauth2.googleapis.com/token"
}

func (g *GoogleProvider) Name() string {
	return "google"
}
