package openapi

type LicenseKeyword int

const (
	MIT LicenseKeyword = iota + 1
	Apache
)

func (l LicenseKeyword) String() string {
	switch l {
	case MIT:
		return "MIT"
	case Apache:
		return "Apache-2.0"
	default:
		return "UNDEFINED"
	}
}

type Info struct {
	Version        string   `json:"version"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier,omitempty"`
	URL        string `json:"url,omitempty"`
}

func newInfo(name, email string, license LicenseKeyword) *Info {
	return &Info{
		License: &License{
			Name: license.String(),
		},
		Contact: &Contact{
			Name:  name,
			Email: email,
		},
	}
}
