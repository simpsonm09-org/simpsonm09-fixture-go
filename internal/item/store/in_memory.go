// Package store holds the domain-facing repository adapter. The service depends
// on the port; this package implements it.
package store

import (
	"sort"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/domain"
)

// DefaultItems are the three items seeded on startup so the API has something
// to return.
var DefaultItems = []domain.Item{
	{Name: "Widget", Description: strPtr("A small widget")},
	{Name: "Gadget", Description: strPtr("A handy gadget")},
	{Name: "Gizmo", Description: strPtr("A clever gizmo")},
}

// InMemoryItemRepository is the in-memory store adapter. It converges to the
// seeded state on restart, so a restart discards every item the caller created
// and restores the three seeds.
type InMemoryItemRepository struct {
	items  map[int64]domain.Item
	nextID int64
}

// NewInMemoryItemRepository returns a store seeded with DefaultItems.
func NewInMemoryItemRepository() *InMemoryItemRepository {
	repository := &InMemoryItemRepository{}
	repository.Seed(DefaultItems)
	return repository
}

// List returns every item ordered by id.
func (r *InMemoryItemRepository) List() []domain.Item {
	items := make([]domain.Item, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

// Find returns the item with the given id and whether it exists.
func (r *InMemoryItemRepository) Find(id int64) (domain.Item, bool) {
	item, ok := r.items[id]
	return item, ok
}

// Save stores the item, assigning the next id when its ID is zero.
func (r *InMemoryItemRepository) Save(item domain.Item) domain.Item {
	id := item.ID
	if id == 0 {
		r.nextID++
		id = r.nextID
	}
	stored := domain.Item{ID: id, Name: item.Name, Description: item.Description}
	r.items[id] = stored
	return stored
}

// Delete removes the item with the given id.
func (r *InMemoryItemRepository) Delete(id int64) {
	delete(r.items, id)
}

// Exists reports whether an item with the given id is stored.
func (r *InMemoryItemRepository) Exists(id int64) bool {
	_, ok := r.items[id]
	return ok
}

// Seed replaces every record with the given items and resets the id sequence.
func (r *InMemoryItemRepository) Seed(seeds []domain.Item) {
	r.items = make(map[int64]domain.Item, len(seeds))
	r.nextID = 0
	for _, seed := range seeds {
		r.Save(domain.Item{Name: seed.Name, Description: seed.Description})
	}
}

// Clear empties the store and resets the id sequence.
func (r *InMemoryItemRepository) Clear() {
	r.items = make(map[int64]domain.Item)
	r.nextID = 0
}

func strPtr(value string) *string {
	return &value
}
