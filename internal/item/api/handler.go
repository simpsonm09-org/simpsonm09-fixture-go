package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/domain"
	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/service"
)

// ItemHandler serves the item CRUD endpoints. It speaks DTOs and delegates to
// the service; it never touches the store.
type ItemHandler struct {
	service *service.ItemService
}

// NewItemHandler returns a handler over the given service.
func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{service: itemService}
}

// Register mounts the item routes on the router.
func (h *ItemHandler) Register(router gin.IRouter) {
	items := router.Group("/items")
	items.GET("", h.list)
	items.GET("/:id", h.get)
	items.POST("", h.create)
	items.PUT("/:id", h.update)
	items.DELETE("/:id", h.delete)
}

// list godoc
//
// @Summary      List every item
// @Tags         Items
// @Produce      json
// @Success      200  {array}   api.ItemResponse
// @Router       /items [get]
func (h *ItemHandler) list(c *gin.Context) {
	c.JSON(http.StatusOK, toResponses(h.service.List()))
}

// get godoc
//
// @Summary      Get one item by id
// @Tags         Items
// @Produce      json
// @Param        id   path      int64  true  "Item id"
// @Success      200  {object}  api.ItemResponse
// @Failure      404  {object}  api.Problem  "Item not found"
// @Router       /items/{id} [get]
func (h *ItemHandler) get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	item, err := h.service.Get(id)
	if err != nil {
		writeNotFound(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(item))
}

// create godoc
//
// @Summary      Create an item
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        request  body      api.ItemRequest  true  "Item to create"
// @Success      201      {object}  api.ItemResponse
// @Failure      400      {object}  api.Problem  "Validation failed"
// @Router       /items [post]
func (h *ItemHandler) create(c *gin.Context) {
	var request ItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeProblem(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}
	c.JSON(http.StatusCreated, toResponse(h.service.Create(request.Name, request.Description)))
}

// update godoc
//
// @Summary      Replace an item
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        id       path      int64            true  "Item id"
// @Param        request  body      api.ItemRequest  true  "Item to store"
// @Success      200      {object}  api.ItemResponse
// @Failure      400      {object}  api.Problem  "Validation failed"
// @Failure      404      {object}  api.Problem  "Item not found"
// @Router       /items/{id} [put]
func (h *ItemHandler) update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var request ItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeProblem(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}
	item, err := h.service.Update(id, request.Name, request.Description)
	if err != nil {
		writeNotFound(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(item))
}

// delete godoc
//
// @Summary      Delete an item
// @Tags         Items
// @Param        id   path  int64  true  "Item id"
// @Success      204  "Item deleted"
// @Failure      404  {object}  api.Problem  "Item not found"
// @Router       /items/{id} [delete]
func (h *ItemHandler) delete(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(id); err != nil {
		writeNotFound(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func pathID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		writeProblem(c, http.StatusBadRequest, "Invalid item id", "Item id must be an integer")
		return 0, false
	}
	return id, true
}

func writeNotFound(c *gin.Context, err error) {
	var notFound domain.NotFoundError
	if errors.As(err, &notFound) {
		writeProblem(c, http.StatusNotFound, "Item not found", notFound.Error())
		return
	}
	writeProblem(c, http.StatusInternalServerError, "Internal error", "unexpected error")
}

func writeProblem(c *gin.Context, status int, title, detail string) {
	c.Header("Content-Type", "application/problem+json")
	c.JSON(status, Problem{Type: "about:blank", Title: title, Status: status, Detail: detail})
}
