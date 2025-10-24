package v1

type Team struct {
	ID     int64  `json:"id"`
	ClubID int64  `json:"clubID"`
	Name   string `json:"name"`
	League string `json:"league"`
}

// +openapi:schema:title="create team request model"
type CreateTeamRequest struct {
	// +openapi:schema:required
	ClubID int64 `json:"clubID"`

	// +openapi:schema:required
	Name string `json:"name"`

	// +openapi:schema:required
	League string `json:"league"`
}

// +openapi:schema:title="create team response model"
type CreateTeamResponse struct {
	ID int64 `json:"id"`
}
