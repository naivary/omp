package openapi

type SecurityType int

const (
	API_KEY SecurityType = iota + 1
	HTTP
	MUTUAL_TLS
	OAUTH2
	OPENID_CONNECT
)

type Style int

const (
	MATRIX Style = iota + 1
	LABEL
	SIMPLE
	FORM
	SPACE_DELIM
	PIPE_DELIM
	DEEPOBJ
	S_COOKIE
)

type OpenAPI struct {
	Version           string                 `json:"openapi"`
	Self              string                 `json:"$self,omitempty"`
	Info              *Info                  `json:"info"`
	JSONSchemaDialect string                 `json:"jsonSchemaDialect,omitempty"`
	Servers           []*Server              `json:"servers,omitempty"`
	Paths             map[string]*PathItem   `json:"paths"`
	Webhooks          map[string]*PathItem   `json:"webhooks,omitempty"`
	Components        *Components            `json:"components,omitempty"`
	Security          []*SecurityRequirement `json:"security,omitempty"`
	Tags              []*Tag                 `json:"tags,omitempty"`
	ExternalDocs      *ExternalDoc           `json:"externalDocs,omitempty"`
}

func New(openAPIVersion, contactName, email string, license LicenseKeyword) *OpenAPI {
	return &OpenAPI{
		Version: openAPIVersion,
		Info:    newInfo(contactName, email, license),
		Paths:   make(map[string]*PathItem),
		Components: &Components{
			Schemas:         make(map[string]*Schema),
			Responses:       make(map[string]*Response),
			Parameters:      make(map[string]*Parameter),
			Examples:        make(map[string]*Example),
			RequestBodies:   make(map[string]*RequestBody),
			Headers:         make(map[string]*Header),
			SecuritySchemas: make(map[string]*SecurityScheme),
			Links:           make(map[string]*Link),
			Callbacks:       make(map[string]*PathItem),
			PathItems:       make(map[string]*PathItem),
			MediaTypes:      make(map[string]*MediaType),
		},
	}
}

type Server struct {
	URL         string                     `json:"url"`
	Description string                     `json:"description,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Default     string   `json:"default"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Operation struct {
	Tags         []string              `json:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty"`
	Description  string                `json:"description,omitempty"`
	ExternalDocs *ExternalDoc          `json:"externalDocs,omitempty"`
	OperationID  string                `json:"operationId,omitempty"`
	Parameters   []*Parameter          `json:"parameters,omitempty"`
	RequestBody  *RequestBody          `json:"requestBody,omitempty"`
	Responses    map[string]*Response  `json:"responses"`
	Callbacks    map[string]*PathItem  `json:"callbacks,omitempty"`
	Deprecated   bool                  `json:"deprecated,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty"`
	Servers      []*Server             `json:"servers,omitempty"`
}

type ExternalDoc struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Example struct {
	Ref             string `json:"$ref,omitempty"`
	Summary         string `json:"summary,omitempty"`
	Description     string `json:"description,omitempty"`
	DataValue       any    `json:"dataValue,omitempty"`
	SerializedValue string `json:"serializedValue,omitempty"`
	ExternalValue   string `json:"externalValue,omitempty"`
	Value           any    `json:"value,omitempty"`
}

type MediaType struct {
	Schema         *Schema              `json:"schema,omitempty,omitzero"`
	ItemSchema     *Schema              `json:"itemSchema,omitempty"`
	Example        any                  `json:"example,omitempty"`
	Examples       map[string]*Example  `json:"examples,omitempty"`
	Encoding       map[string]*Encoding `json:"encoding,omitempty"`
	PrefixEncoding []*Encoding          `json:"prefixEncoding,omitempty"`
	ItemEncoding   *Encoding            `json:"itemEncoding,omitempty"`
}

type Encoding struct {
	ContentType    string               `json:"contentType,omitempty"`
	Headers        map[string]*Header   `json:"headers,omitempty"`
	Encoding       map[string]*Encoding `json:"encoding,omitempty"`
	PrefixEncoding []*Encoding          `json:"prefixEncoding,omitempty"`
	ItemEncoding   *Encoding            `json:"itemEncoding,omitempty"`
}

type Link struct {
	OperationRef string         `json:"operationRef,omitempty"`
	OperationID  string         `json:"operationId,omitempty"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	RequestBody  map[string]any `json:"requestBody,omitempty"`
	Description  string         `json:"description,omitempty"`
	Server       *Server        `json:"server,omitempty"`
}

type SecurityRequirement map[string][]string

type Tag struct {
	Name         string       `json:"name"`
	Summary      string       `json:"summary,omitempty"`
	Description  string       `json:"description,omitempty"`
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`
	Parent       string       `json:"parent,omitempty"`
	Kind         string       `json:"kind,omitempty"`
}

type Components struct {
	Schemas         map[string]*Schema         `json:"schemas,omitempty"`
	Responses       map[string]*Response       `json:"responses,omitempty"`
	Parameters      map[string]*Parameter      `json:"parameters,omitempty"`
	Examples        map[string]*Example        `json:"examples,omitempty"`
	RequestBodies   map[string]*RequestBody    `json:"requestBodies,omitempty"`
	Headers         map[string]*Header         `json:"headers,omitempty"`
	SecuritySchemas map[string]*SecurityScheme `json:"securitySchemes,omitempty"`
	Links           map[string]*Link           `json:"links,omitempty"`
	Callbacks       map[string]*PathItem       `json:"callbacks,omitempty"`
	PathItems       map[string]*PathItem       `json:"pathItems,omitempty"`
	MediaTypes      map[string]*MediaType      `json:"mediaTypes,omitempty"`
}

type SecurityScheme struct {
	Type             SecurityType `json:"type"`
	Description      string       `json:"description,omitempty"`
	Name             string       `json:"name,omitempty"`
	In               string       `json:"in,omitempty"`
	Scheme           string       `json:"scheme,omitempty"`
	BearerFormat     string       `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows  `json:"flows,omitempty"`
	OpenIDConnectURL string       `json:"openIdConnectUrl,omitempty"`
	OAuth2MetdataURL string       `json:"oauth2MetadataUrl,omitempty"`
	Deprecated       bool         `json:"deprecated,omitempty"`
}

type OAuthFlows struct {
	Implicit            *OAuthFlow `json:"implicit,omitempty"`
	Password            *OAuthFlow `json:"password,omitempty"`
	ClientCredentials   *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode   *OAuthFlow `json:"authorizationCode,omitempty"`
	DeviceAuthorization *OAuthFlow `json:"deviceAuthorization,omitempty"`
}

type OAuthFlow struct {
	AuthorizationURL       string            `json:"authorizationUrl,omitempty"`
	DeviceAuthorizationURL string            `json:"deviceAuthorizationUrl,omitempty"`
	TokenURL               string            `json:"tokenUrl,omitempty"`
	RefreshURL             string            `json:"refreshUrl,omitempty"`
	Scopes                 map[string]string `json:"scopes"`
}
