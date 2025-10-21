package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func CreateUser[T any](ctx context.Context, issuer, token, realm string, user T) (int, error) {
	endpoint, err := url.JoinPath(issuer, "admin", "realms", realm, "users")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	data, err := json.Marshal(&user)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	body := bytes.NewReader(data)
	r, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := http.DefaultClient.Do(r)
	return res.StatusCode, err
}
