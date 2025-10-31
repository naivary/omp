package openapi

import "reflect"

type Response struct {
	Ref         string                `json:"$ref,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Headers     map[string]*Header    `json:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
	Links       map[string]*Link      `json:"links,omitempty"`
}

func NewResRef(ref string) *Response {
	return &Response{Ref: ref}
}

func NewResponse[T any](desc string) *Response {
	var model T
	var s *Schema
	typ := reflect.TypeOf(model)
	if typ != nil {
		s = &Schema{Ref: componentRef("schemas", typeName(model))}
	}
	res := &Response{
		Description: desc,
		Headers:     make(map[string]*Header),
		Content: map[string]*MediaType{
			"application/json": {
				Schema: s,
			},
		},
		Links: make(map[string]*Link),
	}
	return res
}

func (r *Response) AddHeader(name string, h *Header) *Response {
	r.Headers[name] = h
	return r
}
