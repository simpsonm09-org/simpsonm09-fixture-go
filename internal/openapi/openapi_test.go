package openapi_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/openapi"
)

// TestCommittedDocumentMatchesGenerator regenerates docs/openapi.json from the
// annotations and fails when the committed document differs, so a hand edit or
// a drifting route goes red.
func TestCommittedDocumentMatchesGenerator(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate the test file")
	}
	root := filepath.Join(filepath.Dir(file), "..", "..")

	generated, err := openapi.Generate(root)
	if err != nil {
		t.Fatalf("generate openapi: %v", err)
	}
	committed, err := os.ReadFile(filepath.Join(root, "docs", "openapi.json"))
	if err != nil {
		t.Fatalf("read docs/openapi.json: %v", err)
	}
	if string(committed) != string(generated) {
		t.Fatalf("docs/openapi.json is out of date; run `just spec`")
	}
}
