package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/iancoleman/strcase"

	"github.com/naivary/omp/openapi"
)

func GenOpenAPISpecs(root *openapi.OpenAPI, endpoints ...*Endpoint) error {
	if root.Info == nil {
		return fmt.Errorf("nil info: OpenAPI.Info is required")
	}
	schemas, err := schemaDefs()
	if err != nil {
		return err
	}
	root.Components.Schemas = schemas
	for _, endpoint := range endpoints {
		err := genOpenAPISpecs(root, endpoint)
		if err != nil {
			return err
		}
	}
	file, err := os.Create("api/openapi/openapi.json")
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(&root)
}

func genOpenAPISpecs(root *openapi.OpenAPI, endpoint *Endpoint) error {
	if len(endpoint.Responses) == 0 {
		return fmt.Errorf("no responses defined")
	}
	patternSegments := strings.SplitN(endpoint.Pattern, " ", 2)
	if len(patternSegments) != 2 {
		return fmt.Errorf("incorrect pattern: %s", endpoint.Pattern)
	}
	method, path := patternSegments[0], patternSegments[1]
	op, err := buildOperation(endpoint)
	if err != nil {
		return err
	}
	root.Paths[path] = openapi.NewPathItem(method, op)
	return nil
}

func buildOperation(endpoint *Endpoint) (*openapi.Operation, error) {
	op := &openapi.Operation{
		Summary:     endpoint.Summary,
		Description: endpoint.Description,
		Parameters:  endpoint.Parameters,
		RequestBody: endpoint.RequestBody,
		Responses:   endpoint.Responses,
		Deprecated:  endpoint.Deprecated,
		Tags:        endpoint.Tags,
	}
	if endpoint.OperationID == "" {
		return nil, fmt.Errorf(
			"empty operation id: make sure to define an operation ID. Usually its best to take the function name as the operation ID",
		)
	}
	return op, nil
}

func schemaDefs() (map[string]*jsonschema.Schema, error) {
	entries, err := os.ReadDir("api/openapi/schemas")
	if err != nil {
		return nil, err
	}
	schemas := make(map[string]*jsonschema.Schema)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		ref := fmt.Sprintf("./schemas/%s", filename)
		typeName := strcase.ToCamel(strings.Split(filename, ".")[0])
		schemas[typeName] = &jsonschema.Schema{Ref: ref}
	}
	return schemas, nil
}
