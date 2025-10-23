package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	keycloakv1 "github.com/naivary/omp/api/keycloak/v1"
)

type Keycloak interface {
	CreateUser(user *keycloakv1.User) error
}

var _ Keycloak = (*keycloak)(nil)

type keycloak struct {
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
	data, err := json.Marshal(&user)
	if err != nil {
		return err
	}
	body := bytes.NewReader(data)
	r, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}
	token, err := k.newToken()
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	bearer := fmt.Sprintf("Bearer %s", token.AccessToken)
	r.Header.Add("Authorization", bearer)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if is4XX(res.StatusCode) {
		var buf bytes.Buffer
		io.Copy(&buf, res.Body)
		return fmt.Errorf("4XX status code: %s", buf.String())
	}
	return err
}

func is4XX(code int) bool {
	return code >= 400 && code < 500
}
