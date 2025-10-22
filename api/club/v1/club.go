package v1

type ClubProfile struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Timezone string `json:"timezone"`
}
