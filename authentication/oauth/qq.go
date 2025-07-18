package oauth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type QQUserInfo struct {
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`

	FigureURL  string `json:"figureurl"`
	FigureURLQ string `json:"figureurl_qq_1"` // 50x50
	FigureURL2 string `json:"figureurl_qq_2"` // 100x100
}

type QQAccessToken struct {
	AccessToken  string `json:"access_token"`  // Access token
	ExpiresIn    int    `json:"expires_in"`    // Token expiration time in seconds
	RefreshToken string `json:"refresh_token"` // Refresh token
	OpenID       string `json:"openid"`        // User's OpenID
	Scope        string `json:"scope"`         // Scope of access
}

type QQProvider struct {
	Client *http.Client
	Config *oauth2.Config
}

/*
整体流程：

第三方发起QQ授权登录请求，QQ用户允许授权第三方应用后，QQ会拉起应用或重定向到第三方网站，并且带上授权临时票据code参数，
1. 获取Authorization Code；
2. 通过获取Authorization Code，加上参数client_id等，通过API获取Access Token；
3. 通过Access Token进行接口调用，获取用户基本数据资源或帮助用户实现基本操作。

QQ二维码链接
配置链接参数，生成完整的QQ二维码链接；
 https://graph.qq.com/oauth2.0/authorize[?](https://open.weixin.qq.com/connect/qrconnect?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect)which=Login&response_type=code&redirect_uri=redirect_uri&client_id=client_id&state=state


*/
// NewQQProvider creates a new QQ OAuth provider instance.
// NewQQProvider 创建一个新的QQ OAuth提供者实例。
// clientId 和 clientSecret 是在QQ开放平台注册应用时获取的应用ID和密钥。
// RedirectURL 是授权回调地址，必须与在QQ开放平台注册应用时设置的回调地址一致。
// 注意：QQ的OAuth2.0流程与其他OAuth2.0提供者略有不同，特别是在获取Access Token和用户信息的API上。
func NewQQProvider(clientId string, clientSecret string) *QQProvider {
	auth := &QQProvider{
		Config: &oauth2.Config{
			Scopes: []string{"get_user_info"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://graph.qq.com/oauth2.0/authorize",
				TokenURL: "https://graph.qq.com/oauth2.0/token",
			},
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  "https://www.inkgo.com/callback", // Set your redirect URL
		},
		Client: &http.Client{
			Transport: &http.Transport{

				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second, // Set a timeout for the dial
					KeepAlive: 30 * time.Second, // Set keep-alive for the connection
				}).DialContext,
			},
		},
	}
	return auth
}

// GetAccessToken retrieves the access token using the provided code.
// 使用提供的code获取访问令牌。 code 是从QQ授权回调中获取的。
func (auth *QQProvider) GetAccessToken(code string) (*oauth2.Token, error) {
	params := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     auth.Config.ClientID,
		"client_secret": auth.Config.ClientSecret,
		"code":          code,
		"redirect_uri":  auth.Config.RedirectURL,
	}
	data, err := auth.postWithBody(params, auth.Config.Endpoint.TokenURL)
	if err != nil {
		return nil, err
	}

	qqToken := &QQAccessToken{}
	if err = json.Unmarshal(data, &qqToken); err != nil {
		return nil, err
	}

	if qqToken.AccessToken == "" {
		return nil, fmt.Errorf("failed to get access token: %s", string(data))
	}

	if qqToken.OpenID == "" {
		return nil, fmt.Errorf("failed to get OpenID from QQ access token response: %s", string(data))
	}

	raw := map[string]interface{}{
		"openid": qqToken.OpenID,
	}

	token := oauth2.Token{
		AccessToken:  qqToken.AccessToken,
		TokenType:    "QQToken",
		RefreshToken: qqToken.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(qqToken.ExpiresIn) * time.Second),
	}

	token.WithExtra(raw)

	return &token, nil

}

// GetUserInfo retrieves user information from QQ using the provided access token.
// 使用提供的访问令牌从QQ获取用户信息。
func (auth *QQProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	if token == nil || token.AccessToken == "" {
		return nil, fmt.Errorf("invalid access token")
	}

	params := map[string]string{
		"access_token":       token.AccessToken,
		"oauth_consumer_key": auth.Config.ClientID,
		"openid":             token.Extra("openid").(string),
	}

	data, err := auth.getWithParams(params, "https://graph.qq.com/user/get_user_info")
	if err != nil {
		return nil, err
	}

	var qqUserInfo QQUserInfo
	if err = json.Unmarshal(data, &qqUserInfo); err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		UnionID:     qqUserInfo.OpenID,
		URL:         "https://graph.qq.com/user/get_user_info",
		Provider:    auth.Name(),
		Username:    qqUserInfo.Nickname,
		DisplayName: qqUserInfo.Nickname,
		Email:       "", // QQ does not provide email
		Avatar:      qqUserInfo.FigureURL2,
	}

	return userInfo, nil

}

func (auth *QQProvider) Name() string {
	return "qq"
}

func (auth *QQProvider) postWithBody(params map[string]string, rawURL string) ([]byte, error) {
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}
	resp, err := auth.Client.PostForm(rawURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to post to %s: %s", rawURL, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (auth *QQProvider) getWithParams(params map[string]string, rawURL string) ([]byte, error) {
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	resp, err := auth.Client.Get(rawURL + "?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get from %s: %s", rawURL, resp.Status)
	}

	return io.ReadAll(resp.Body)
}
