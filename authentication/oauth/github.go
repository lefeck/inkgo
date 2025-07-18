package oauth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GithubAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
}

type GithubUserInfo struct {
	Login                   string      `json:"login"`
	ID                      int         `json:"id"`
	NodeId                  string      `json:"node_id"`
	AvatarUrl               string      `json:"avatar_url"`
	GravatarId              string      `json:"gravatar_id"`
	Url                     string      `json:"url"`
	HtmlUrl                 string      `json:"html_url"`
	FollowersUrl            string      `json:"followers_url"`
	FollowingUrl            string      `json:"following_url"`
	GistsUrl                string      `json:"gists_url"`
	StarredUrl              string      `json:"starred_url"`
	SubscriptionsUrl        string      `json:"subscriptions_url"`
	OrganizationsUrl        string      `json:"organizations_url"`
	ReposUrl                string      `json:"repos_url"`
	EventsUrl               string      `json:"events_url"`
	ReceivedEventsUrl       string      `json:"received_events_url"`
	Type                    string      `json:"type"`
	SiteAdmin               bool        `json:"site_admin"`
	Name                    string      `json:"name"`
	Company                 string      `json:"company"`
	Blog                    string      `json:"blog"`
	Location                string      `json:"location"`
	Email                   string      `json:"email"`
	Hireable                bool        `json:"hireable"`
	Bio                     string      `json:"bio"`
	TwitterUsername         interface{} `json:"twitter_username"`
	PublicRepos             int         `json:"public_repos"`
	PublicGists             int         `json:"public_gists"`
	Followers               int         `json:"followers"`
	Following               int         `json:"following"`
	CreatedAt               time.Time   `json:"created_at"`
	UpdatedAt               time.Time   `json:"updated_at"`
	PrivateGists            int         `json:"private_gists"`
	TotalPrivateRepos       int         `json:"total_private_repos"`
	OwnedPrivateRepos       int         `json:"owned_private_repos"`
	DiskUsage               int         `json:"disk_usage"`
	Collaborators           int         `json:"collaborators"`
	TwoFactorAuthentication bool        `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		Collaborators int    `json:"collaborators"`
		PrivateRepos  int    `json:"private_repos"`
	} `json:"plan"`
}

type GithubProvider struct {
	Config *oauth2.Config
	Client *http.Client
}

var (
	defaultHttpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
)

func NewGithubProvider(clientid string, clientSecret string) *GithubProvider {
	return &GithubProvider{
		Config: &oauth2.Config{
			Scopes: []string{"user:email", "read:user"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "http://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
			ClientID:     clientid,
			ClientSecret: clientSecret,
		},
		Client: defaultHttpClient,
	}
}

func (auth *GithubProvider) Name() string {
	return "github"
}

func (auth *GithubProvider) GetAccessToken(code string) (*oauth2.Token, error) {
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
	pToken := &GithubAccessToken{}
	if err = json.Unmarshal(data, pToken); err != nil {
		return nil, err
	}
	if pToken.Error != "" {
		return nil, fmt.Errorf("err: %s", pToken.Error)
	}
	token := &oauth2.Token{
		AccessToken: pToken.AccessToken,
		TokenType:   "Bearer",
	}

	return token, nil

}

func (auth *GithubProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {

	if token == nil || token.AccessToken == "" {
		return nil, nil
	}

	client := auth.Client
	if client == nil {
		client = defaultHttpClient
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
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

	var githubUserInfo GithubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&githubUserInfo); err != nil {
		return nil, err
	}

	return &UserInfo{
		UnionID:     strconv.Itoa(githubUserInfo.ID),
		URL:         githubUserInfo.HtmlUrl,
		Provider:    auth.Name(),
		Username:    githubUserInfo.Login,
		DisplayName: githubUserInfo.Name,
		Email:       githubUserInfo.Email,
		Avatar:      githubUserInfo.AvatarUrl,
	}, nil

}

func (auth *GithubProvider) postWithBody(body interface{}, url string) ([]byte, error) {
	bs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	r := strings.NewReader(string(bs))
	req, _ := http.NewRequest("POST", url, r)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := auth.Client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return data, nil
}
