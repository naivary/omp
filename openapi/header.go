package openapi

type Header struct {
	Description string                `json:"description,omitempty"`
	Required    bool                  `json:"required,omitempty"`
	Deprecated  bool                  `json:"deprecated,omitempty"`
	Example     any                   `json:"example,omitempty"`
	Examples    map[string]*Example   `json:"examples,omitempty"`
	Schema      *Schema               `json:"schema,omitempty"`
	Style       Style                 `json:"style,omitempty"`
	Explode     bool                  `json:"explode,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

func NewHeader(desc string, required bool, s *Schema) *Header {
	h := &Header{
		Description: desc,
		Required:    required,
		Schema:      s,
	}
	return h
}

func (h *Header) Deprecate() *Header {
	h.Deprecated = true
	return h
}

func (h *Header) AddExample(exp any) *Header {
	h.Example = exp
	return h
}
