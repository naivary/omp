package v1

type User struct {
	Email       string            `json:"email"`
	Enabled     bool              `json:"enabled"`
	Credentials []*Credential     `json:"credentials"`
	Attributes  map[string]string `json:"attributes"`
}

type Credential struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}
