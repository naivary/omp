package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func NewTokenForClient(ctx context.Context, issuer, realm, clientID, clientSecret string) (*oauth2.Token, error) {
	tokenURL, err := url.JoinPath(issuer, "realms", realm, "protocol", "openid-connect", "token")
	if err != nil {
		return nil, err
	}
	clientCredsConfig := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}
	return clientCredsConfig.Token(ctx)
}

func CreateUser[T any](ctx context.Context, issuer, token, realm string, user T) (int, error) {
	endpoint, err := url.JoinPath(issuer, "admin", "realms", realm, "users")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	data, err := json.Marshal(&user)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	fmt.Println(string(data))
	body := bytes.NewReader(data)
	r, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := http.DefaultClient.Do(r)
	fmt.Println(res.Header)
	return res.StatusCode, err
}
