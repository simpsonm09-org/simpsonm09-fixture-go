// Package domain holds the item model the service reasons about, free of
// transport and store detail.
package domain

import "fmt"

// Item is an item as the service layer works with it. ID is zero until the
// store assigns it on save. Description is nil when the caller omits it.
type Item struct {
	ID          int64
	Name        string
	Description *string
}

// NotFoundError is raised when an item id has no matching record. The API layer
// maps it to an HTTP 404 problem detail.
type NotFoundError struct {
	ID int64
}

// Error reports the id that was not found.
func (e NotFoundError) Error() string {
	return fmt.Sprintf("Item %d was not found", e.ID)
}
