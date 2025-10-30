package keycloak

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	keycloakv1 "github.com/naivary/omp/api/keycloak/v1"
)

type Keycloak interface {
	RemoveUser(ctx context.Context, email string) error
	CreateUser(ctx context.Context, user *keycloakv1.User) error
	EnableUser(ctx context.Context, email string) error
	GetUserID(ctx context.Context, email string) (string, error)
	GetUser(ctx context.Context, email string) (*keycloakv1.User, error)
	IsEmailUsed(ctx context.Context, email string) (bool, error)
}

var _ Keycloak = (*keycloak)(nil)

type keycloak struct {
	// url of the keycloak server
	url string
	// realm managed by the provided client credentials
	realm        string
	clientID     string
	clientSecret string

	credsConfig *clientcredentials.Config
	cl          *http.Client
}

func New(ctx context.Context, urll, realm, clientID, clientSecret string, insecureSkipVerify bool) (Keycloak, error) {
	tokenURL, err := url.JoinPath(urll, "realms", realm, "protocol", "openid-connect", "token")
	if err != nil {
		return nil, err
	}
	ccConfig := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}
	k := &keycloak{
		url:          urll,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		credsConfig:  ccConfig,
		cl: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
			},
		},
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, k.cl)
	_, err = ccConfig.Token(ctx)
	return k, err
}

func (k *keycloak) newToken(ctx context.Context) (*oauth2.Token, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, k.cl)
	return k.credsConfig.Token(ctx)
}

func (k *keycloak) CreateUser(ctx context.Context, user *keycloakv1.User) error {
	endpoint, err := url.JoinPath(k.url, "admin", "realms", k.realm, "users")
	if err != nil {
		return err
	}
	token, err := k.newToken(ctx)
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
	res, err := k.cl.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return newError(res)
}

func (k *keycloak) EnableUser(ctx context.Context, email string) error {
	token, err := k.newToken(ctx)
	if err != nil {
		return err
	}
	id, err := k.GetUserID(ctx, email)
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
	if err != nil {
		return err
	}
	res, err := k.cl.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return newError(res)
}

func (k *keycloak) GetUserID(ctx context.Context, email string) (string, error) {
	user, err := k.GetUser(ctx, email)
	if err != nil {
		return "", err
	}
	return user.ID, err
}

func (k *keycloak) RemoveUser(ctx context.Context, email string) error {
	token, err := k.newToken(ctx)
	if err != nil {
		return err
	}
	userID, err := k.GetUserID(ctx, email)
	if err != nil {
		return err
	}
	endpoint, err := url.JoinPath(k.url, "admin", "realms", k.realm, "users", userID)
	if err != nil {
		return err
	}
	r, err := newRequest[any](http.MethodDelete, endpoint, nil, http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", token.AccessToken)},
	}, nil)
	if err != nil {
		return err
	}
	res, err := k.cl.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return newError(res)
}

func (k *keycloak) GetUser(ctx context.Context, email string) (*keycloakv1.User, error) {
	token, err := k.newToken(ctx)
	if err != nil {
		return nil, err
	}
	endpoint, err := url.JoinPath(k.url, "admin", "realms", k.realm, "users")
	if err != nil {
		return nil, err
	}
	r, err := newRequest[any](http.MethodGet, endpoint, nil,
		http.Header{
			"Authorization": {fmt.Sprintf("Bearer %s", token.AccessToken)},
		},
		url.Values{
			"email": {email},
		},
	)
	if err != nil {
		return nil, err
	}
	res, err := k.cl.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if isError(res) {
		return nil, newError(res)
	}
	users := make([]*keycloakv1.User, 0)
	err = json.NewDecoder(res.Body).Decode(&users)
	if len(users) == 0 {
		return nil, nil
	}
	if len(users) > 1 {
		return nil, fmt.Errorf("more users found with the same email: %s", email)
	}
	return users[0], err
}

func (k *keycloak) IsEmailUsed(ctx context.Context, email string) (bool, error) {
	user, err := k.GetUser(ctx, email)
	return user != nil, err
}
