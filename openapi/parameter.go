package openapi

import "bytes"

type In int

const (
	PATH In = iota + 1
	QUERY
	HEADER
	COOKIE
)

func (i In) String() string {
	switch i {
	case PATH:
		return "path"
	case QUERY:
		return "query"
	case HEADER:
		return "header"
	case COOKIE:
		return "cookie"
	default:
		return "UNDEFINED"
	}
}

func (i In) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(i.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

type Parameter struct {
	Ref             string                `json:"$ref,omitempty"`
	Name            string                `json:"name,omitempty"`
	In              In                    `json:"in,omitempty"`
	Description     string                `json:"description,omitempty"`
	Required        bool                  `json:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty"`
	Example         any                   `json:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty"`
	Schema          *Schema               `json:"schema,omitempty"`
	Style           Style                 `json:"style,omitempty"`
	Explode         bool                  `json:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty"`
}

func NewQueryParam(name, desc string, required bool) *Parameter {
	param := &Parameter{
		Name:        name,
		Description: desc,
		Required:    required,
		In:          QUERY,
	}
	return param
}

func NewCookieParam(name, desc string, required bool) *Parameter {
	param := &Parameter{
		Name:        name,
		Description: desc,
		Required:    required,
		In:          COOKIE,
	}
	return param
}

func NewHeaderParam(name, desc string, required bool) *Parameter {
	param := &Parameter{
		Name:        name,
		Description: desc,
		Required:    required,
		In:          HEADER,
	}
	return param
}

func NewPathParam(name string, s *Schema) *Parameter {
	param := &Parameter{
		Name:     name,
		Required: true,
		In:       PATH,
		Schema:   s,
	}
	return param
}

func (p *Parameter) Deprecate() *Parameter {
	p.Deprecated = true
	return p
}
