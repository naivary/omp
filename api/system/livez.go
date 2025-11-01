package system

// +openapi:schema:title="livez request model"
type LivezRequest struct {
	Verbose int
}

// +openapi:schema:title="readyz response"
type LivezResponse struct {
	Status string
}
