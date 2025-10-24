package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	keycloakv1 "github.com/naivary/omp/api/keycloak/v1"
)

type Keycloak interface {
	CreateUser(user *keycloakv1.User) error
	EnableUser(email string) error
	GetUserID(email string) (string, error)
}

var _ Keycloak = (*keycloak)(nil)

type keycloak struct {
	// ctx the instance was created from. When this context is cancelled no
	// further request should be accepted
	ctx context.Context

	// url of the keycloak server
	url string
	// realm managed by the provided client credentials
	realm        string
	clientID     string
	clientSecret string

	credsConfig *clientcredentials.Config
}

func New(ctx context.Context, urll, realm, clientID, clientSecret string) (Keycloak, error) {
	tokenURL, err := url.JoinPath(urll, "realms", realm, "protocol", "openid-connect", "token")
	if err != nil {
		return nil, err
	}
	ccConfig := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}
	// check if you can get a token to verify that the provided credentials are
	// valid.
	_, err = ccConfig.Token(ctx)
	k := &keycloak{
		ctx:          ctx,
		url:          urll,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		credsConfig:  ccConfig,
	}
	return k, err
}

func (k *keycloak) newToken() (*oauth2.Token, error) {
	return k.credsConfig.Token(k.ctx)
}

func (k *keycloak) CreateUser(user *keycloakv1.User) error {
	endpoint, err := url.JoinPath(k.url, "admin", "realms", k.realm, "users")
	if err != nil {
		return err
	}
	token, err := k.newToken()
	if err != nil {
		return err
	}
	r, err := newRequest(http.MethodPost, endpoint, user,
		http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {fmt.Sprintf("Bearer %s", token.AccessToken)},
		},
		url.Values{},
	)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if is4XX(res.StatusCode) {
		return newError(res.Body)
	}
	return err
}

func (k *keycloak) EnableUser(email string) error {
	token, err := k.newToken()
	if err != nil {
		return err
	}
	id, err := k.GetUserID(email)
	if err != nil {
		return err
	}
	endpoint, err := url.JoinPath(k.url, "admin", "realms", k.realm, "users", id)
	if err != nil {
		return err
	}
	r, err := newRequest(http.MethodPut, endpoint, map[string]bool{"enabled": true},
		http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", token.AccessToken)},
		},
		url.Values{},
	)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if is4XX(res.StatusCode) {
		return newError(res.Body)
	}
	return err
}

func (k *keycloak) GetUserID(email string) (string, error) {
	id := ""
	token, err := k.newToken()
	if err != nil {
		return id, err
	}
	endpoint, err := url.JoinPath(k.url, "admin", "realms", k.realm, "users")
	if err != nil {
		return id, err
	}
	r, err := newRequest[any](http.MethodGet, endpoint, nil,
		http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", token.AccessToken)},
		},
		url.Values{
			"email": {email},
		},
	)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return id, err
	}
	defer res.Body.Close()
	if is4XX(res.StatusCode) {
		return id, newError(res.Body)
	}
	ids := []struct {
		ID string `json:"id"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&ids)
	if err != nil {
		return id, err
	}
	id = ids[0].ID
	return id, nil
}
