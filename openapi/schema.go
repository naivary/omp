package openapi

type JSONType string

const (
	InvalidType JSONType = "invalid"
	NullType    JSONType = "null"
	BooleanType JSONType = "boolean"
	NumberType  JSONType = "number"
	IntegerType JSONType = "integer"
	StringType  JSONType = "string"
	ArrayType   JSONType = "array"
	ObjectType  JSONType = "object"
)

type Schema struct {
	// metadata
	ID    string   `json:"$id,omitempty"`
	Draft string   `json:"$schema,omitempty"`
	Ref   string   `json:"$ref,omitempty"`
	Type  JSONType `json:"type,omitempty"`

	OneOf []*Schema `json:"oneOf,omitempty"`
	AnyOf []*Schema `json:"anyOf,omitempty"`
	Not   *Schema   `json:"not,omitempty"`

	// agnostic
	Enum []any `json:"enum,omitempty"`

	// annotations
	Title      string `json:"title,omitempty"`
	Desc       string `json:"description,omitempty"`
	Examples   []any  `json:"examples,omitempty"`
	Deprecated bool   `json:"deprecated,omitempty"`
	WriteOnly  bool   `json:"writeOnly,omitempty"`
	ReadOnly   bool   `json:"readOnly,omitempty"`
	Default    string `json:"default,omitempty"`

	// array
	MaxItems    int64   `json:"maxItems,omitempty"`
	MinItems    int64   `json:"minItems,omitempty"`
	UniqueItems bool    `json:"uniqueItems,omitempty"`
	Items       *Schema `json:"items,omitempty"`

	// object
	Properties           map[string]*Schema  `json:"properties,omitempty"`
	Required             []string            `json:"required,omitempty"`
	AdditionalProperties *Schema             `json:"additionalProperties,omitempty"`
	PatternProperties    map[string]*Schema  `json:"patternProperties,omitempty"`
	DependentRequired    map[string][]string `json:"dependentRequired,omitempty"`

	// string
	MinLength        int64  `json:"minLength,omitempty"`
	MaxLength        int64  `json:"maxLength,omitempty"`
	Pattern          string `json:"pattern,omitempty"`
	ContentEncoding  string `json:"contentEnconding,omitempty"`
	ContentMediaType string `json:"contentMediaType,omitempty"`
	Format           string `json:"format,omitempty"`

	// number
	Maximum          int64 `json:"maximum,omitempty"`
	Minimum          int64 `json:"minimum,omitempty"`
	ExclusiveMaximum int64 `json:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum int64 `json:"exclusiveMinimum,omitempty"`
	MultipleOf       int64 `json:"multipleOf,omitempty"`
}
