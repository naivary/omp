package v1

type Profile struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Timezone string `json:"timezone"`
}

// +openapi:schema:title="create club profile request"
type CreateClubRequest struct {
	// +openapi:schema:required
	// +openapi:schema:format="email"
	Email string `json:"email"`

	// +openapi:schema:required
	// +openapi:schema:minLength=16
	// +openapi:schema:maxLength=32
	Password string `json:"password"`

	// +openapi:schema:required
	Name string `json:"name"`

	// +openapi:schema:required
	Location string `json:"location"`

	// +openapi:schema:required
	Timezone string `json:"timezone"`
}

// +openapi:schema:title="create club response"
type CreateClubResponse struct {
	ID int64 `json:"id"`
}

// +openapi:schema:title="delete club request"
type DeleteClubRequest struct {
	// +openapi:schema:required
	// +openapi:schema:format="email"
	Email string

	// +openapi:schema:required
	Name string
}
