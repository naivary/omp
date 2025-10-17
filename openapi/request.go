package openapi

import "reflect"

type RequestBody struct {
	Description string                `json:"description,omitempty"`
	Required    bool                  `json:"required"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

func NewReqBody(desc string, required bool, model any) *RequestBody {
	var s *Schema
	if model != nil {
		s = &Schema{Ref: reflect.TypeOf(model).Name()}
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
