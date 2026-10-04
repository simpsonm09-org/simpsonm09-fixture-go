package service_test

import (
	"errors"
	"testing"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/domain"
	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/service"
)

type stubRepository struct {
	listed  []domain.Item
	found   domain.Item
	foundOK bool
	saved   []domain.Item
	exists  bool
	deleted []int64
}

func (r *stubRepository) List() []domain.Item { return r.listed }

func (r *stubRepository) Find(int64) (domain.Item, bool) { return r.found, r.foundOK }

func (r *stubRepository) Save(item domain.Item) domain.Item {
	r.saved = append(r.saved, item)
	return item
}

func (r *stubRepository) Delete(id int64) { r.deleted = append(r.deleted, id) }

func (r *stubRepository) Exists(int64) bool { return r.exists }

func ptr(value string) *string { return &value }

func TestListReturnsRepositoryItems(t *testing.T) {
	items := []domain.Item{{ID: 1, Name: "Widget", Description: ptr("A small widget")}}
	itemsService := service.NewItemService(&stubRepository{listed: items})

	got := itemsService.List()

	if len(got) != 1 || got[0].Name != "Widget" {
		t.Fatalf("List() = %+v, want one Widget", got)
	}
}

func TestGetReturnsItem(t *testing.T) {
	repository := &stubRepository{found: domain.Item{ID: 1, Name: "Widget"}, foundOK: true}

	item, err := service.NewItemService(repository).Get(1)

	if err != nil {
		t.Fatalf("Get(1) error = %v", err)
	}
	if item.Name != "Widget" {
		t.Fatalf("Get(1).Name = %q, want Widget", item.Name)
	}
}

func TestGetMissingReturnsNotFound(t *testing.T) {
	_, err := service.NewItemService(&stubRepository{}).Get(9)

	var notFound domain.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("Get(9) error = %v, want domain.NotFoundError", err)
	}
	if notFound.ID != 9 {
		t.Fatalf("NotFoundError.ID = %d, want 9", notFound.ID)
	}
}

func TestCreateSavesThroughTheStore(t *testing.T) {
	repository := &stubRepository{}

	created := service.NewItemService(repository).Create("Gadget", ptr("A handy gadget"))

	if len(repository.saved) != 1 {
		t.Fatalf("Save called %d times, want 1", len(repository.saved))
	}
	if saved := repository.saved[0]; saved.ID != 0 || saved.Name != "Gadget" {
		t.Fatalf("saved = %+v, want zero id and Gadget", saved)
	}
	if created.Name != "Gadget" {
		t.Fatalf("Create().Name = %q, want Gadget", created.Name)
	}
}

func TestUpdateReplacesExistingItem(t *testing.T) {
	repository := &stubRepository{exists: true}

	updated, err := service.NewItemService(repository).Update(1, "Renamed", nil)

	if err != nil {
		t.Fatalf("Update(1) error = %v", err)
	}
	if saved := repository.saved[0]; saved.ID != 1 || saved.Name != "Renamed" || saved.Description != nil {
		t.Fatalf("saved = %+v, want id 1 and Renamed", saved)
	}
	if updated.Name != "Renamed" {
		t.Fatalf("Update().Name = %q, want Renamed", updated.Name)
	}
}

func TestUpdateMissingDoesNotSave(t *testing.T) {
	repository := &stubRepository{}

	_, err := service.NewItemService(repository).Update(404, "Nope", nil)

	var notFound domain.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("Update(404) error = %v, want domain.NotFoundError", err)
	}
	if len(repository.saved) != 0 {
		t.Fatalf("Save called %d times, want 0", len(repository.saved))
	}
}

func TestDeleteRemovesExistingItem(t *testing.T) {
	repository := &stubRepository{exists: true}

	err := service.NewItemService(repository).Delete(1)

	if err != nil {
		t.Fatalf("Delete(1) error = %v", err)
	}
	if len(repository.deleted) != 1 || repository.deleted[0] != 1 {
		t.Fatalf("deleted = %v, want [1]", repository.deleted)
	}
}

func TestDeleteMissingDoesNotDelete(t *testing.T) {
	repository := &stubRepository{}

	err := service.NewItemService(repository).Delete(404)

	var notFound domain.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("Delete(404) error = %v, want domain.NotFoundError", err)
	}
	if len(repository.deleted) != 0 {
		t.Fatalf("Delete called %d times, want 0", len(repository.deleted))
	}
}
