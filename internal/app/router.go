// Package app wires the item layers into an HTTP engine.
package app

import (
	"github.com/gin-gonic/gin"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/api"
	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/service"
	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/store"
)

// NewRouter builds a Gin engine with the item routes, a fresh in-memory store
// seeded with the three dev items, and the request validator installed.
func NewRouter() *gin.Engine {
	api.UseValidator()
	router := gin.New()
	router.Use(gin.Recovery())
	handler := api.NewItemHandler(service.NewItemService(store.NewInMemoryItemRepository()))
	handler.Register(router)
	return router
}
