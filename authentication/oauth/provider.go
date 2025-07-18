package oauth

import (
	"fmt"
	"golang.org/x/oauth2"
	"inkgo/config"
	"inkgo/model"
)

type OAuth2Provider interface {
	// GetName returns the name of the OAuth provider.
	Name() string
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
	// GetAccessToken retrieves the access token using the provided code.
	GetAccessToken(code string) (*oauth2.Token, error)
}

type OAuthInfo struct {
	OpenID   string // 微信/QQ 特有
	UnionID  string // 微信开放平台
	Email    string
	Name     string
	Avatar   string
	Provider string // google / wechat / qq 等
}

type UserInfo struct {
	UnionID     string
	URL         string // URL for the user's profile
	Provider    string // OAuth provider (e.g., google, wechat, qq)
	Username    string // Username or unique identifier for the user
	DisplayName string // Display name of the user
	Email       string // Email address of the user
	Avatar      string // URL to the user's avatar image
}

func (u *UserInfo) User() *model.User {
	return &model.User{
		UserName: u.Username,
		Email:    u.Email,
		Avatar:   u.Avatar,
		OAuthAccounts: []model.OAuthAccount{
			{
				Provider:   u.Provider,
				ProviderID: u.UnionID, // Assuming Username is unique for the provider
				URL:        u.URL,
			},
		},
	}
}

const (
	OAuthGoogle   = "google"
	OAuthWeChat   = "wechat"
	OAuthQQ       = "qq"
	OAuthGitHub   = "github"
	EmptyAuthType = "nil"
)

func IsEmptyAuthType(authType string) bool {
	return authType == "" || authType == EmptyAuthType
}

type OAuthManager struct {
	config map[string]config.OAuthConfig
}

func NewOAuthManager(conf map[string]config.OAuthConfig) *OAuthManager {
	return &OAuthManager{
		config: conf,
	}
}

func (m *OAuthManager) GetAuthProvider(authType string) (OAuth2Provider, error) {
	var provider OAuth2Provider
	conf, ok := m.config[authType]
	if !ok {
		return nil, fmt.Errorf("auth type not found in config: %s ", authType)
	}
	switch authType {
	case OAuthGitHub:
		provider = NewGithubProvider(conf.ClientID, conf.ClientSecret)
	case OAuthWeChat:
		provider = NewWeChatAuth(conf.ClientID, conf.ClientSecret)
	case OAuthQQ:
		provider = NewQQProvider(conf.ClientID, conf.ClientSecret)
	//case OAuthGoogle:
	//	provider = NewGoogleProvider(conf)
	default:
		return nil, fmt.Errorf("unknown auth type: %s", authType)
	}

	return provider, nil

}
