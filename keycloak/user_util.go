package keycloak

import (
	"strconv"

	kcv1 "github.com/naivary/omp/api/keycloak/v1"
)

type Attributes struct {
	ProfileID int64
}

func (a Attributes) Map() map[string][]string {
	return map[string][]string{
		"profileID": {strconv.FormatInt(a.ProfileID, 10)},
	}
}

func NewUser(email, password string, roles []string, attr *Attributes) *kcv1.User {
	user := &kcv1.User{
		Email: email,
		Credentials: []*kcv1.Credential{
			{Type: "password", Value: password},
		},
		Attributes: attr.Map(),
	}
	return user
}
