package keycloak

import (
	"context"
	"fmt"
	"net/url"

	"golang.org/x/oauth2/clientcredentials"
)

type clubManager struct {
	issuer       string
	clientID     string
	clientSecret string
}

func NewClubManager(ctx context.Context, issuer, clientID, clientSecret string) (*clubManager, error) {
	clubMngr := &clubManager{
		issuer:       issuer,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
	tokenURL, err := url.JoinPath(issuer, "realms", "clubs", "protocol", "openid-connect", "token")
	if err != nil {
		return nil, err
	}
	clientCredsConfig := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}
	jwtToken, err := clientCredsConfig.Token(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println(jwtToken.AccessToken)
	return clubMngr, nil
}
