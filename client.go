package goplurk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kinabcd/goplurk/oauth"
)

var baseURL = "https://www.plurk.com"

type Engine interface {
	CallAPI(_url string, opt map[string]string) ([]byte, error)
	CallAPIUnmarshal(_url string, opt map[string]string, v any) error
}
type EngineImpl struct {
	oauthClient *oauth.Client
	Credentials *oauth.Credentials
}
type Client struct {
	Users     *APIUsers
	Timeline  *APITimeline
	Responses *APIResponses
	Profile   *APIProfile
	Polling   *APIPolling
	Realtime  *APIRealtime
	Alerts    *APIAlerts
	Engine    Engine
}

func NewClient(consumerToken string, consumerSecret string, accessToken string, accessSecret string) (*Client, error) {
	if consumerToken == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return nil, fmt.Errorf("can not be authorized")
	}

	return newClient(newOAuthClient(consumerToken, consumerSecret), &oauth.Credentials{
		Token:  accessToken,
		Secret: accessSecret,
	}), nil
}

func newClient(oauthClient *oauth.Client, credentials *oauth.Credentials) *Client {
	client := &Client{}
	client.Engine = &EngineImpl{
		oauthClient: oauthClient,
		Credentials: credentials,
	}
	client.Users = &APIUsers{client: client}
	client.Timeline = &APITimeline{client: client}
	client.Responses = &APIResponses{client: client}
	client.Profile = &APIProfile{client: client}
	client.Polling = &APIPolling{client: client}
	client.Realtime = &APIRealtime{client: client}
	client.Alerts = &APIAlerts{client: client}
	return client
}

func newOAuthClient(consumerToken string, consumerSecret string) *oauth.Client {
	return &oauth.Client{
		TemporaryCredentialRequestURI: "https://www.plurk.com/OAuth/request_token",
		ResourceOwnerAuthorizationURI: "https://www.plurk.com/OAuth/authorize",
		TokenRequestURI:               "https://www.plurk.com/OAuth/access_token",
		Credentials: oauth.Credentials{
			Token:  consumerToken,
			Secret: consumerSecret,
		},
	}
}

func (c *EngineImpl) CallAPI(_url string, opt map[string]string) ([]byte, error) {
	var apiURL = baseURL + _url
	param := make(url.Values)
	for k, v := range opt {
		param.Set(k, v)
	}
	c.oauthClient.SignForm(c.Credentials, "POST", apiURL, param)
	res, err := http.PostForm(apiURL, url.Values(param))
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %s, %s, %v", apiURL, fmt.Sprint(param), err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body))
	}
	return body, nil
}
func (c *EngineImpl) CallAPIUnmarshal(_url string, opt map[string]string, v any) error {
	if bytes, err := c.CallAPI(_url, opt); err != nil {
		return err
	} else if err := json.Unmarshal(bytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal: %v, %s", err, string(bytes))
	} else {
		return nil
	}
}

type OAuthRequest struct {
	client               *oauth.Client
	Url                  string
	temporaryCredentials *oauth.Credentials
}

func NewOAuthRequest(consumerToken string, consumerSecret string) (*OAuthRequest, error) {
	oauthClient := newOAuthClient(consumerToken, consumerSecret)
	requestToken, err := oauthClient.RequestTemporaryCredentials(http.DefaultClient, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request temporary credentials: %v", err)
	}
	url := oauthClient.AuthorizationURL(requestToken, nil)

	return &OAuthRequest{
		client:               oauthClient,
		Url:                  url,
		temporaryCredentials: requestToken,
	}, nil
}

func (c *OAuthRequest) SendPin(pin string) (*Client, string, string, error) {
	credentials, _, err := c.client.RequestToken(http.DefaultClient, c.temporaryCredentials, pin)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to request accessToken: %v", err)
	}
	return newClient(c.client, credentials), credentials.Token, credentials.Secret, nil
}
