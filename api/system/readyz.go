package system

// +openapi:schema:title="readyz request model"
type ReadyzRequest struct {
	Verbose int
}

// +openapi:schema:title="readyz response"
type ReadyzResponse struct {
	Status string
}
