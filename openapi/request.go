package openapi

import "reflect"

type RequestBody struct {
	Description string                `json:"description,omitempty"`
	Required    bool                  `json:"required"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

func NewReqBody[T any](desc string, required bool) *RequestBody {
	var model T
	var s *Schema
	typ := reflect.TypeOf(model)
	if typ != nil {
		s = &Schema{Ref: componentRef("schemas", typ.Name())}
	}
	req := &RequestBody{
		Description: desc,
		Required:    required,
		Content: map[string]*MediaType{
			"application/json": {
				Schema: s,
			},
		},
	}
	return req
}
