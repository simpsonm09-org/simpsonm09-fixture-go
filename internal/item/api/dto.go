// Package api is the HTTP transport layer. It speaks request and response DTOs
// and never touches the store.
package api

// ItemRequest is the payload used to create or replace an item. ID is assigned
// by the store, so it is absent from requests.
//
// The binding tag is inert at runtime; the request validator reads validate.
// It exists so swag marks the field required in docs/openapi.json, which it
// cannot infer from the custom notblank rule.
type ItemRequest struct {
	Name        string  `json:"name" example:"Widget" maxLength:"200" binding:"required" validate:"notblank,max=200"`
	Description *string `json:"description" example:"A small widget" maxLength:"2000" validate:"omitempty,max=2000"`
}

// ItemResponse is an item returned by the API. Every field is always present,
// so each is marked required for the generated document.
type ItemResponse struct {
	ID          int64   `json:"id" example:"1" format:"int64" binding:"required"`
	Name        string  `json:"name" example:"Widget" binding:"required"`
	Description *string `json:"description" example:"A small widget" binding:"required"`
}
