package club

import (
	"context"
	"fmt"

	clubv1 "github.com/naivary/omp/api/club/v1"
	"github.com/naivary/omp/keycloak"
)

type clubManager struct {
	issuer       string
	clientID     string
	clientSecret string
}

func NewManager(ctx context.Context, issuer, clientID, clientSecret string) (*clubManager, error) {
	clubMngr := &clubManager{
		issuer:       issuer,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
	_, err := keycloak.NewTokenForClient(ctx, issuer, "clubs", clientID, clientSecret)
	return clubMngr, err
}

func (c *clubManager) CreateUser(user *clubv1.User) error {
	user.Enabled = true
	ctx := context.Background()
	token, err := keycloak.NewTokenForClient(ctx, c.issuer, "clubs", c.clientID, c.clientSecret)
	if err != nil {
		return err
	}
	code, err := keycloak.CreateUser(context.Background(), c.issuer, token.AccessToken, "clubs", user)
	fmt.Println(code, err)
	return err
}
