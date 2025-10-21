package v1

import (
	"encoding/json"
	"strconv"
)

type Club struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Timezone string `json:"timezone"`
}

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Enabled   bool   `json:"enabled"`
	ClubID    int
}

type keycloakUser struct {
	Email         string            `json:"email"`
	FirstName     string            `json:"firstName"`
	LastName      string            `json:"lastName"`
	Enabled       bool              `json:"enabled"`
	EmailVerified bool              `json:"emailVerified"`
	Attributes    map[string]string `json:"attributes"`
}

var _ json.Marshaler = (*User)(nil)

func (u User) MarshalJSON() ([]byte, error) {
	kcUser := keycloakUser{
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Enabled:   u.Enabled,
		Attributes: map[string]string{
			"clubID": strconv.Itoa(u.ClubID),
		},
	}
	return json.Marshal(&kcUser)
}
