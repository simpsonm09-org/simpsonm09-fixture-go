package openapi_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/openapi"
)

// TestCommittedDocumentMatchesGenerator regenerates docs/openapi.json from the
// annotations and fails when the committed document differs, so a hand edit or
// a drifting route goes red.
func TestCommittedDocumentMatchesGenerator(t *testing.T) {
	generated := generate(t)

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate the test file")
	}
	root := filepath.Join(filepath.Dir(file), "..", "..")

	committed, err := os.ReadFile(filepath.Join(root, "docs", "openapi.json"))
	if err != nil {
		t.Fatalf("read docs/openapi.json: %v", err)
	}
	if string(committed) != string(generated) {
		t.Fatalf("docs/openapi.json is out of date; run `just spec`")
	}
}

// generatedDocument is the subset of the OpenAPI document the contract test
// reads. It decodes the generator output, so a regression in the normalize pass
// fails here even when the committed file and the generated file agree.
type generatedDocument struct {
	Components struct {
		Schemas map[string]struct {
			Required   []string                  `json:"required"`
			Properties map[string]map[string]any `json:"properties"`
		} `json:"schemas"`
	} `json:"components"`
	Paths map[string]map[string]struct {
		Responses map[string]struct {
			Content map[string]json.RawMessage `json:"content"`
		} `json:"responses"`
	} `json:"paths"`
}

// TestGeneratedDocumentPinsContracts asserts the contract choices that the
// annotations and the normalize pass must produce, reading the generator output
// directly rather than the committed document.
func TestGeneratedDocumentPinsContracts(t *testing.T) {
	var document generatedDocument
	if err := json.Unmarshal(generate(t), &document); err != nil {
		t.Fatalf("decode generated openapi: %v", err)
	}

	assertRequired(t, document, "api.ItemRequest", "name")
	assertRequired(t, document, "api.ItemResponse", "id", "name", "description")

	for _, schemaName := range []string{"api.ItemRequest", "api.ItemResponse"} {
		property, ok := document.Components.Schemas[schemaName].Properties["description"]
		if !ok {
			t.Fatalf("%s has no description property", schemaName)
		}
		types, ok := property["type"].([]any)
		if !ok || len(types) != 2 || types[0] != "string" || types[1] != "null" {
			t.Fatalf("%s description type = %v, want [string null]", schemaName, property["type"])
		}
	}

	assertProblemJSON(t, document)
}

// generate returns the document the generator produces from the annotations.
func generate(t *testing.T) []byte {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate the test file")
	}
	root := filepath.Join(filepath.Dir(file), "..", "..")

	generated, err := openapi.Generate(root)
	if err != nil {
		t.Fatalf("generate openapi: %v", err)
	}
	return generated
}

func assertRequired(t *testing.T, document generatedDocument, schemaName string, want ...string) {
	t.Helper()
	schema, ok := document.Components.Schemas[schemaName]
	if !ok {
		t.Fatalf("generated document has no schema %s", schemaName)
	}
	got := make(map[string]bool, len(schema.Required))
	for _, name := range schema.Required {
		got[name] = true
	}
	if len(schema.Required) != len(want) {
		t.Fatalf("%s required = %v, want %v", schemaName, schema.Required, want)
	}
	for _, name := range want {
		if !got[name] {
			t.Fatalf("%s required = %v, want %v", schemaName, schema.Required, want)
		}
	}
}

// assertProblemJSON fails when any 4xx or 5xx response is not served as
// application/problem+json, including the 404s on the item routes.
func assertProblemJSON(t *testing.T, document generatedDocument) {
	t.Helper()
	errorResponses := 0
	for path, operations := range document.Paths {
		for method, operation := range operations {
			for status, response := range operation.Responses {
				if len(status) == 0 || (status[0] != '4' && status[0] != '5') {
					continue
				}
				errorResponses++
				if _, ok := response.Content["application/problem+json"]; !ok {
					t.Errorf("%s %s %s content = %v, want application/problem+json", method, path, status, keysOf(response.Content))
				}
			}
		}
	}
	if errorResponses == 0 {
		t.Fatal("generated document has no error responses to check")
	}
}

func keysOf(content map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(content))
	for key := range content {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
