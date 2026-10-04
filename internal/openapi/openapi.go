// Package openapi generates the OpenAPI document from the Go annotations, so
// the committed document is code-derived and cannot drift.
//
// swag emits OpenAPI 3.1 but cannot express two contract choices: that an
// omitted description is a nullable string, and that an error response carries
// application/problem+json. normalize applies those two representations to the
// swag output, so docs/openapi.json is still produced end to end by this
// generator and never edited by hand.
package openapi

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/swaggo/swag/v2/gen"
)

const (
	jsonMediaType    = "application/json"
	problemMediaType = "application/problem+json"
)

// Generate parses the swag annotations under root and returns the OpenAPI 3.1
// document as the committed docs/openapi.json content, with a trailing newline.
func Generate(root string) ([]byte, error) {
	dir, err := os.MkdirTemp("", "swaggo-openapi")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(dir) }()

	generator := gen.New()
	config := &gen.Config{
		SearchDir:           root,
		MainAPIFile:         "main.go",
		OutputDir:           dir,
		OutputTypes:         []string{"json"},
		GenerateOpenAPI3Doc: true,
		ParseInternal:       true,
		ParseGoList:         true,
	}
	if err := generator.Build(config); err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(filepath.Join(dir, "swagger.json"))
	if err != nil {
		return nil, err
	}
	normalized, err := normalize(raw)
	if err != nil {
		return nil, err
	}
	return append(normalized, '\n'), nil
}

func normalize(raw []byte) ([]byte, error) {
	var document map[string]any
	if err := json.Unmarshal(raw, &document); err != nil {
		return nil, err
	}
	// swag emits an externalDocs stub with an empty url when none is set.
	delete(document, "externalDocs")
	markNullableDescription(document)
	useProblemMediaType(document)
	return json.MarshalIndent(document, "", "    ")
}

// markNullableDescription types the optional description property as a nullable
// string, the OpenAPI 3.1 spelling of a value that may be null.
func markNullableDescription(document map[string]any) {
	for _, schema := range componentSchemas(document) {
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			continue
		}
		description, ok := properties["description"].(map[string]any)
		if !ok || description["type"] != "string" {
			continue
		}
		description["type"] = []string{"string", "null"}
	}
}

// useProblemMediaType serves every error response as application/problem+json.
func useProblemMediaType(document map[string]any) {
	paths, ok := document["paths"].(map[string]any)
	if !ok {
		return
	}
	for _, pathValue := range paths {
		pathItem, ok := pathValue.(map[string]any)
		if !ok {
			continue
		}
		for _, operationValue := range pathItem {
			operation, ok := operationValue.(map[string]any)
			if !ok {
				continue
			}
			responses, ok := operation["responses"].(map[string]any)
			if !ok {
				continue
			}
			for status, responseValue := range responses {
				response, ok := responseValue.(map[string]any)
				if !ok || !isErrorStatus(status) {
					continue
				}
				content, ok := response["content"].(map[string]any)
				if !ok {
					continue
				}
				media, ok := content[jsonMediaType]
				if !ok {
					continue
				}
				delete(content, jsonMediaType)
				content[problemMediaType] = media
			}
		}
	}
}

func componentSchemas(document map[string]any) []map[string]any {
	components, ok := document["components"].(map[string]any)
	if !ok {
		return nil
	}
	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(schemas))
	for _, value := range schemas {
		if schema, ok := value.(map[string]any); ok {
			out = append(out, schema)
		}
	}
	return out
}

func isErrorStatus(status string) bool {
	if len(status) == 0 {
		return false
	}
	return status[0] == '4' || status[0] == '5'
}
