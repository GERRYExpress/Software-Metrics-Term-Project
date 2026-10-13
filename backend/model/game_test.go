package model

import (
	"reflect"
	"testing"
)

func TestGame_AddCategory(t *testing.T) {
	t.Run("successfully appends category", func(t *testing.T) {
		game := &Game{
			Categories: make(Categories, 0),
		}

		newCat := Category{
			ID:   ID(1),
			Name: Name("Action"),
		}

		game.AddCategory(newCat)

		if len(game.Categories) != 1 {
			t.Fatalf("AddCategory() failed to append, expected length 1, got %d", len(game.Categories))
		}

		if !reflect.DeepEqual(game.Categories[0], newCat) {
			t.Errorf("AddCategory() appended wrong category. got %v, want %v", game.Categories[0], newCat)
		}
	})
}
