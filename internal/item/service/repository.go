package service

import "github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/domain"

// ItemRepository is the store port the service depends on. It speaks the
// domain type, so the service never sees a store record and the adapter can be
// swapped without touching business logic.
type ItemRepository interface {
	// List returns every item, ordered by id.
	List() []domain.Item
	// Find returns the item with the given id and whether it exists.
	Find(id int64) (domain.Item, bool)
	// Save stores the item. When its ID is zero the store assigns the next id.
	Save(item domain.Item) domain.Item
	// Delete removes the item with the given id. It is a no-op when absent.
	Delete(id int64)
	// Exists reports whether an item with the given id is stored.
	Exists(id int64) bool
}
