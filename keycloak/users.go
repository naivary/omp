package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func CreateUser[T any](issuer, token, realm string, user T) (int, error) {
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
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := http.DefaultClient.Do(r)
	return res.StatusCode, err
}

type User struct {
	ID string `json:"id"`
}

func SetPasswordOfUser(issuer, token, realm, email, password string) (int, error) {
	endpoint, err := url.JoinPath(issuer, "admin", "realms", realm, "users")
	if err != nil {
		return 0, err
	}
	r, err := http.NewRequest(http.MethodGet, endpoint, nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	q := r.URL.Query()
	q.Set("email", email)
	r.URL.RawQuery = q.Encode()
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	users := make([]User, 0)
	err = json.NewDecoder(res.Body).Decode(&users)
	if err != nil {
		return 0, err
	}
	userID := users[0].ID

	return res.StatusCode, err
}
