package openapi

import (
	"encoding/json"
	"testing"
)

// TestNormalizeTraversalGuards exercises each early return in the useProblemMediaType
// traversal. swag emits well-formed documents, so the guards only fire on shapes
// the generator never produces; the tests below pin the behavior that a malformed
// node is skipped and the rest of the document is still normalized.
func TestNormalizeTraversalGuards(t *testing.T) {
	for _, tt := range normalizeTraversalGuardCases {
		t.Run(tt.name, func(t *testing.T) {
			assertNormalize(t, tt.document, tt.want)
		})
	}
}

// normalizeTraversalGuardCases is the table for TestNormalizeTraversalGuards,
// kept at package scope so the test function stays within the funlen gate.
var normalizeTraversalGuardCases = []struct {
	name     string
	document map[string]any
	want     string
}{
	{
		name:     "no paths",
		document: map[string]any{},
		want:     "{}",
	},
	{
		name:     "paths is not an object",
		document: map[string]any{"paths": "not-an-object"},
		want:     `{"paths":"not-an-object"}`,
	},
	{
		name:     "path item is not an object",
		document: map[string]any{"paths": map[string]any{"/items": "not-an-object"}},
		want:     `{"paths":{"/items":"not-an-object"}}`,
	},
	{
		name: "operation has no responses",
		document: map[string]any{
			"paths": map[string]any{
				"/items": map[string]any{"get": map[string]any{"summary": "list"}},
			},
		},
		want: `{"paths":{"/items":{"get":{"summary":"list"}}}}`,
	},
	{
		name: "responses is not an object",
		document: map[string]any{
			"paths": map[string]any{
				"/items": map[string]any{"get": map[string]any{"responses": "not-an-object"}},
			},
		},
		want: `{"paths":{"/items":{"get":{"responses":"not-an-object"}}}}`,
	},
	{
		name: "response is not an object",
		document: map[string]any{
			"paths": map[string]any{
				"/items": map[string]any{
					"get": map[string]any{"responses": map[string]any{"404": "not-an-object"}},
				},
			},
		},
		want: `{"paths":{"/items":{"get":{"responses":{"404":"not-an-object"}}}}}`,
	},
	{
		name: "success response is untouched",
		document: map[string]any{
			"paths": map[string]any{
				"/items": map[string]any{
					"get": map[string]any{"responses": map[string]any{
						"200": map[string]any{
							"content": map[string]any{jsonMediaType: map[string]any{"schema": map[string]any{}}},
						},
					}},
				},
			},
		},
		want: `{"paths":{"/items":{"get":{"responses":{"200":{"content":{"application/json":{"schema":{}}}}}}}}}`,
	},
	{
		name: "error response has no content",
		document: map[string]any{
			"paths": map[string]any{
				"/items": map[string]any{
					"get": map[string]any{"responses": map[string]any{
						"404": map[string]any{"description": "not found"},
					}},
				},
			},
		},
		want: `{"paths":{"/items":{"get":{"responses":{"404":{"description":"not found"}}}}}}`,
	},
	{
		name: "error response content is not an object",
		document: map[string]any{
			"paths": map[string]any{
				"/items": map[string]any{
					"get": map[string]any{"responses": map[string]any{
						"500": map[string]any{"content": "not-an-object"},
					}},
				},
			},
		},
		want: `{"paths":{"/items":{"get":{"responses":{"500":{"content":"not-an-object"}}}}}}`,
	},
	{
		name: "error response content lacks json media type",
		document: map[string]any{
			"paths": map[string]any{
				"/items": map[string]any{
					"get": map[string]any{"responses": map[string]any{
						"422": map[string]any{
							"content": map[string]any{"text/plain": map[string]any{"schema": map[string]any{}}},
						},
					}},
				},
			},
		},
		want: `{"paths":{"/items":{"get":{"responses":{"422":{"content":{"text/plain":{"schema":{}}}}}}}}}`,
	},
}

// assertNormalize runs normalize over the document and compares the compacted
// result with want. The compaction removes key order and indentation, so the
// assertion pins the traversed values rather than the formatting.
func assertNormalize(t *testing.T, document map[string]any, want string) {
	t.Helper()
	raw, err := json.Marshal(document)
	if err != nil {
		t.Fatalf("marshal input: %v", err)
	}
	got, err := normalize(raw)
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	var compact map[string]any
	if err := json.Unmarshal(got, &compact); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	remarshaled, err := json.Marshal(compact)
	if err != nil {
		t.Fatalf("remarshal output: %v", err)
	}
	if string(remarshaled) != want {
		t.Fatalf("normalize output = %s, want %s", remarshaled, want)
	}
}

// TestNormalizeRewritesErrorResponse checks the positive path in
// useProblemMediaTypeForResponse that the guard table does not reach: a 4xx or
// 5xx response served as application/json becomes application/problem+json.
func TestNormalizeRewritesErrorResponse(t *testing.T) {
	raw := []byte(`{"paths":{"/items":{"get":{"responses":{"404":{"content":{"application/json":{"schema":{}}}}}}}}}`)
	got, err := normalize(raw)
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	var document map[string]any
	if err := json.Unmarshal(got, &document); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	content := document["paths"].(map[string]any)["/items"].(map[string]any)["get"].(map[string]any)["responses"].(map[string]any)["404"].(map[string]any)["content"].(map[string]any)
	if _, ok := content[problemMediaType]; !ok {
		t.Fatalf("content = %v, want application/problem+json", content)
	}
	if _, ok := content[jsonMediaType]; ok {
		t.Fatalf("content = %v, want application/json removed", content)
	}
}
