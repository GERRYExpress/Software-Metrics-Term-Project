package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	a := NewAccount(1, "testuser")
	assert.Equal(t, NewID(1), a.ID)
	assert.Equal(t, NewUsername("testuser"), a.Username)
	assert.Empty(t, a.Password)
}

func TestNewCategory(t *testing.T) {
	c := NewCategory(5, "Action", true)
	assert.Equal(t, NewID(5), c.ID)
	assert.Equal(t, NewName("Action"), c.Name)
	assert.True(t, bool(c.IsHidden))
}

func TestNewGame(t *testing.T) {
	g := NewGame(10, "Test Game", "Description", 29.99, "poster.jpg")
	assert.Equal(t, NewID(10), g.ID)
	assert.Equal(t, NewName("Test Game"), g.Name)
	assert.Equal(t, NewDescription("Description"), g.Description)
	assert.Equal(t, NewPrice(29.99), g.Price)
	assert.Equal(t, NewPosterURL("poster.jpg"), g.PosterUrl)
	assert.Empty(t, g.Categories)
}

func TestNewGameDetail(t *testing.T) {
	releaseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	g := NewGameDetail(1, "Detail Game", "Desc", 19.99, "poster.jpg", releaseDate)
	assert.Equal(t, NewID(1), g.ID)
	assert.Equal(t, NewName("Detail Game"), g.Name)
	assert.Equal(t, NewReleaseDate(releaseDate), g.ReleaseDate)
}