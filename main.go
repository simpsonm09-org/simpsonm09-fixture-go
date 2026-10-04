// @title           Simpsonm09 Fixture Go API
// @version         0.1.0
// @description     Item CRUD service for the simpsonm09 repository fixture
package main

import (
	"log"
	"os"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if err := app.NewRouter().Run("127.0.0.1:" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
