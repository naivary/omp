//go:generate go tool go-enum --marshal --sql --nocomments
package v1

import (
	"context"
	"net/mail"
)

// ENUM(Right, Left)
type StrongFoot int

type PlayerProfile struct {
	ID           int64      `json:"id"`
	TeamID       int64      `json:"teamID"`
	Name         string     `json:"name"`
	JerseyNumber uint       `json:"jerseyNumber"`
	StrongFoot   StrongFoot `json:"strongFoot"`
	Position     string     `json:"position"`
}

// +openapi:schema:title="create new player request"
type CreatePlayerRequest struct {
	// +openapi:schema:format="email"
	// +openapi:schema:required
	Email string `json:"email"`

	// +openapi:schema:minLength=16
	// +openapi:schema:maxLength=32
	// +openapi:schema:required
	Password string `json:"password"`

	// +openapi:schema:required
	ClubID int64 `json:"clubID"`
}

func (c CreatePlayerRequest) Validate(ctx context.Context) map[string]string {
	problems := make(map[string]string, 3)
	if c.Email == "" {
		problems["email"] = "email is empty"
	}
	if c.Password == "" {
		problems["password"] = "password is empty"
	}
	_, err := mail.ParseAddress(c.Email)
	if err != nil {
		problems["email_format"] = "cannot parse email: " + err.Error()
	}
	return problems
}
