// Command spec regenerates docs/openapi.json from the Go annotations.
package main

import (
	"log"
	"os"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/openapi"
)

func main() {
	content, err := openapi.Generate(".")
	if err != nil {
		log.Fatalf("generate openapi: %v", err)
	}
	if err := os.WriteFile("docs/openapi.json", content, 0o644); err != nil {
		log.Fatalf("write docs/openapi.json: %v", err)
	}
}
