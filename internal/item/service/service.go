// Package service holds the item business logic. It works in domain.Item and
// depends on the ItemRepository port, so it never sees a transport DTO or a
// store record.
package service

import "github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/domain"

// ItemService applies the item rules over the repository port.
type ItemService struct {
	repository ItemRepository
}

// NewItemService returns a service backed by the given repository.
func NewItemService(repository ItemRepository) *ItemService {
	return &ItemService{repository: repository}
}

// List returns every item.
func (s *ItemService) List() []domain.Item {
	return s.repository.List()
}

// Get returns the item with the given id, or a domain.NotFoundError.
func (s *ItemService) Get(id int64) (domain.Item, error) {
	item, ok := s.repository.Find(id)
	if !ok {
		return domain.Item{}, domain.NotFoundError{ID: id}
	}
	return item, nil
}

// Create stores a new item and returns it with its assigned id.
func (s *ItemService) Create(name string, description *string) domain.Item {
	return s.repository.Save(domain.Item{Name: name, Description: description})
}

// Update replaces the item with the given id, or returns a
// domain.NotFoundError when it does not exist.
func (s *ItemService) Update(id int64, name string, description *string) (domain.Item, error) {
	if !s.repository.Exists(id) {
		return domain.Item{}, domain.NotFoundError{ID: id}
	}
	return s.repository.Save(domain.Item{ID: id, Name: name, Description: description}), nil
}

// Delete removes the item with the given id, or returns a
// domain.NotFoundError when it does not exist.
func (s *ItemService) Delete(id int64) error {
	if !s.repository.Exists(id) {
		return domain.NotFoundError{ID: id}
	}
	s.repository.Delete(id)
	return nil
}
