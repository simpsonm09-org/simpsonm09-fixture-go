package store_test

import (
	"testing"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/domain"
	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/item/store"
)

func ptr(value string) *string { return &value }

func TestNewSeedsTheDefaultItems(t *testing.T) {
	repository := store.NewInMemoryItemRepository()

	items := repository.List()

	if len(items) != len(store.DefaultItems) {
		t.Fatalf("List() length = %d, want %d", len(items), len(store.DefaultItems))
	}
	names := []string{items[0].Name, items[1].Name, items[2].Name}
	want := []string{"Widget", "Gadget", "Gizmo"}
	for i, name := range names {
		if name != want[i] {
			t.Fatalf("seed %d = %q, want %q", i, name, want[i])
		}
	}
}

func TestSaveAssignsIncreasingIDs(t *testing.T) {
	repository := &store.InMemoryItemRepository{}
	repository.Clear()

	first := repository.Save(domain.Item{Name: "One"})
	second := repository.Save(domain.Item{Name: "Two", Description: ptr("second")})

	if first.ID != 1 || second.ID != 2 {
		t.Fatalf("ids = %d, %d, want 1, 2", first.ID, second.ID)
	}
	stored, ok := repository.Find(2)
	if !ok || stored.Name != "Two" || stored.Description == nil || *stored.Description != "second" {
		t.Fatalf("Find(2) = %+v, %v", stored, ok)
	}
}

func TestSaveKeepsTheIDWhenReplacing(t *testing.T) {
	repository := &store.InMemoryItemRepository{}
	repository.Seed([]domain.Item{{Name: "One"}})

	repository.Save(domain.Item{ID: 1, Name: "Renamed"})

	items := repository.List()
	if len(items) != 1 || items[0].Name != "Renamed" {
		t.Fatalf("List() = %+v, want one renamed item", items)
	}
}

func TestDeleteAndExists(t *testing.T) {
	repository := &store.InMemoryItemRepository{}
	repository.Seed([]domain.Item{{Name: "One"}})

	if !repository.Exists(1) {
		t.Fatal("Exists(1) = false, want true")
	}
	repository.Delete(1)
	if repository.Exists(1) {
		t.Fatal("Exists(1) = true after delete, want false")
	}
	if len(repository.List()) != 0 {
		t.Fatalf("List() = %+v, want empty", repository.List())
	}
}

func TestSeedResetsToTheGivenItems(t *testing.T) {
	repository := &store.InMemoryItemRepository{}
	repository.Clear()

	repository.Seed([]domain.Item{{Name: "Only"}})

	items := repository.List()
	if len(items) != 1 || items[0].Name != "Only" {
		t.Fatalf("List() = %+v, want one Only", items)
	}
}
