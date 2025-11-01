//go:generate go tool go-enum --marshal --nocomments
package openapi

import (
	"github.com/google/jsonschema-go/jsonschema"
)

// ENUM(Path, Query, Header, Cookie)
type ParamIn int

type Parameter struct {
	Ref             string                `json:"$ref,omitempty"`
	Name            string                `json:"name,omitempty"`
	ParamIn         ParamIn               `json:"in,omitempty"`
	Description     string                `json:"description,omitempty"`
	Required        bool                  `json:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty"`
	Example         any                   `json:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty"`
	Schema          *jsonschema.Schema    `json:"schema,omitempty"`
	Style           Style                 `json:"style,omitempty"`
	Explode         bool                  `json:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty"`
}

func NewQueryParam[T any](name, desc string, required bool) *Parameter {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		panic(err)
	}
	param := &Parameter{
		Name:        name,
		Description: desc,
		Required:    required,
		ParamIn:     ParamInQuery,
		Schema:      schema,
	}
	return param
}

func NewCookieParam[T any](name, desc string, required bool, s *Schema) *Parameter {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		panic(err)
	}
	param := &Parameter{
		Name:        name,
		Description: desc,
		Required:    required,
		ParamIn:     ParamInCookie,
		Schema:      schema,
	}
	return param
}

func NewHeaderParam[T any](name, desc string, required bool, s *Schema) *Parameter {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		panic(err)
	}
	param := &Parameter{
		Name:        name,
		Description: desc,
		Required:    required,
		ParamIn:     ParamInHeader,
		Schema:      schema,
	}
	return param
}

func NewPathParam[T any](name string) *Parameter {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		panic(err)
	}
	param := &Parameter{
		Name:     name,
		Required: true,
		ParamIn:  ParamInPath,
		Schema:   schema,
	}
	return param
}

func (p *Parameter) Deprecate() *Parameter {
	p.Deprecated = true
	return p
}
