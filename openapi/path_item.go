package openapi

import "net/http"

type PathItem struct {
	Ref                  string                `json:"$ref,omitempty"`
	Summary              string                `json:"summary,omitempty"`
	Description          string                `json:"description,omitempty"`
	Get                  *Operation            `json:"get,omitempty"`
	Put                  *Operation            `json:"put,omitempty"`
	Post                 *Operation            `json:"post,omitempty"`
	Delete               *Operation            `json:"delete,omitempty"`
	Options              *Operation            `json:"options,omitempty"`
	Head                 *Operation            `json:"head,omitempty"`
	Patch                *Operation            `json:"patch,omitempty"`
	Trace                *Operation            `json:"trace,omitempty"`
	Query                *Operation            `json:"query,omitempty"`
	AdditionalProperties map[string]*Operation `json:"additionalProperties,omitempty"`
	Servers              []*Server             `json:"servers,omitempty"`
	Parameters           []*Parameter          `json:"parameters,omitempty"`
}

func NewPathItem(method string, op *Operation) *PathItem {
	pathItem := &PathItem{}
	pathItem.AddOperation(method, op)
	return pathItem
}

func (p *PathItem) AddOperation(method string, op *Operation) *PathItem {
	switch method {
	case http.MethodGet:
		p.Get = op
	case http.MethodPut:
		p.Put = op
	case http.MethodPost:
		p.Post = op
	case http.MethodDelete:
		p.Delete = op
	case http.MethodOptions:
		p.Options = op
	case http.MethodHead:
		p.Head = op
	case http.MethodPatch:
		p.Patch = op
	case http.MethodTrace:
		p.Trace = op
	case "query":
		p.Query = op
	}
	return p
}
