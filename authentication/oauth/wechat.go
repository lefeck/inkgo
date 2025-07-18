package oauth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"time"
)

// WechatAccessToken 是微信授权登录返回的access_token结构体
type WechatAccessToken struct {
	AccessToken  string `json:"access_token"`  //Interface call credentials
	ExpiresIn    int64  `json:"expires_in"`    //access_token interface call credential timeout time, unit (seconds)
	RefreshToken string `json:"refresh_token"` //User refresh access_token
	Openid       string `json:"openid"`        //Unique ID of authorized user
	Scope        string `json:"scope"`         //The scope of user authorization, separated by commas. (,)
	Unionid      string `json:"unionid"`       //This field will appear if and only if the website application has been authorized by the user's UserInfo.
}

// WechatUserInfo 是微信授权登录返回的用户信息结构体
type WechatUserInfo struct {
	Openid     string `json:"openid"`     // Unique ID of authorized user
	Nickname   string `json:"nickname"`   // User nickname
	Headimgurl string `json:"headimgurl"` // User avatar URL
	Unionid    string `json:"unionid"`    // This field will appear if and only if the website application has been authorized by the user's UserInfo.
	Province   string `json:"province"`   // User province
	City       string `json:"city"`       // User city

}

// WechatAuth 是微信授权登录的结构体
type WechatProvider struct {
	Client *http.Client
	Config *oauth2.Config
}

// NewWeChatAuth 创建一个新的微信授权登录实例
func NewWeChatAuth(clientId string, clientSecret string) *WechatProvider {
	auth := &WechatProvider{
		Config: &oauth2.Config{
			Scopes: []string{"snsapi_login"},
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://graph.qq.com/oauth2.0/token",
			},
			ClientID:     clientId,
			ClientSecret: clientSecret,
		},
		Client: http.DefaultClient,
	}

	return auth
}

func (auth *WechatProvider) Name() string {
	return "wechat"
}

// GetToken 使用code获取access_token
func (auth *WechatProvider) GetAccessToken(code string) (*oauth2.Token, error) {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("appid", auth.Config.ClientID)
	params.Add("secret", auth.Config.ClientSecret)
	params.Add("code", code)

	accessTokenUrl := "https://api.weixin.qq.com/sns/oauth2/access_token?" + params.Encode()
	resp, err := auth.Client.Get(accessTokenUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get access token: %s", resp.Status)
	}

	var wechatToken WechatAccessToken
	if err := json.NewDecoder(resp.Body).Decode(&wechatToken); err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  wechatToken.AccessToken,
		TokenType:    "WechatToken",
		RefreshToken: wechatToken.RefreshToken,
		// expires 是 access_token 的过期时间
		Expiry: time.Now().Add(time.Duration(wechatToken.ExpiresIn) * time.Second),
	}
	if wechatToken.Openid == "" {
		return nil, fmt.Errorf("failed to get OpenID from WeChat access token response")
	}

	raw := map[string]interface{}{
		"openid": wechatToken.Openid,
	}

	token.WithExtra(raw)

	return token, nil
}

// GetUserInfo 使用access_token获取用户信息
func (auth *WechatProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	params := url.Values{}
	params.Add("access_token", token.AccessToken)
	params.Add("openid", token.Extra("openid").(string))
	params.Add("lang", "zh_CN")

	userInfoUrl := "https://api.weixin.qq.com/sns/userinfo?" + params.Encode()
	resp, err := auth.Client.Get(userInfoUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	var wechatUserInfo WechatUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&wechatUserInfo); err != nil {
		return nil, err
	}

	userInfo := UserInfo{

		UnionID:     wechatUserInfo.Unionid,
		Provider:    auth.Name(),
		Username:    wechatUserInfo.Nickname,
		DisplayName: wechatUserInfo.Nickname,
		Avatar:      wechatUserInfo.Headimgurl,
	}

	return &userInfo, nil
}
