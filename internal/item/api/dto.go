// Package api is the HTTP transport layer. It speaks request and response DTOs
// and never touches the store.
package api

// ItemRequest is the payload used to create or replace an item. ID is assigned
// by the store, so it is absent from requests.
type ItemRequest struct {
	Name        string  `json:"name" example:"Widget" maxLength:"200" validate:"notblank,max=200"`
	Description *string `json:"description" example:"A small widget" maxLength:"2000" validate:"omitempty,max=2000"`
}

// ItemResponse is an item returned by the API.
type ItemResponse struct {
	ID          int64   `json:"id" example:"1" format:"int64"`
	Name        string  `json:"name" example:"Widget"`
	Description *string `json:"description" example:"A small widget"`
}
