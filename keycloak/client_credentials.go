package keycloak

import (
	"context"
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
