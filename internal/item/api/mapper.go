package api

import "github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/domain"

// toResponse maps a domain item to its transport shape.
func toResponse(item domain.Item) ItemResponse {
	return ItemResponse{ID: item.ID, Name: item.Name, Description: item.Description}
}

// toResponses maps domain items to transport shapes, returning an empty slice
// rather than nil so an empty list serializes as [].
func toResponses(items []domain.Item) []ItemResponse {
	responses := make([]ItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toResponse(item))
	}
	return responses
}
