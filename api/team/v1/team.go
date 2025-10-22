package v1

type Team struct {
	ID     int64  `json:"id"`
	ClubID int64  `json:"clubID"`
	Name   string `json:"name"`
	League string `json:"league"`
}
