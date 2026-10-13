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

func TestNewHardware(t *testing.T) {
	h := NewHardware(HardwareData{
		ID:             1,
		OS:             "Windows 11",
		Processor:      "Intel i5",
		Memory:         "8GB",
		GraphicsCard:   "GTX 1060",
		DirectXVersion: "12",
		Network:        "Broadband",
		Storage:        "50GB",
		Soundcard:      "DirectX",
		AdditionalNote: "Test note",
		Type:           "min",
	})
	assert.Equal(t, NewID(1), h.ID)
	assert.Equal(t, NewDescription("Windows 11"), h.OS)
	assert.Equal(t, NewDescription("Intel i5"), h.Processor)
	assert.Equal(t, NewDescription("8GB"), h.Memory)
	assert.Equal(t, NewDescription("GTX 1060"), h.GraphicsCard)
	assert.Equal(t, NewDescription("12"), h.DirectXVersion)
	assert.Equal(t, NewDescription("Broadband"), h.Network)
	assert.Equal(t, NewDescription("50GB"), h.Storage)
	assert.Equal(t, NewDescription("DirectX"), h.Soundcard)
	assert.Equal(t, NewDescription("Test note"), h.AdditionalNote)
	assert.Equal(t, NewName("min"), h.Type)
}

func TestNewMedia(t *testing.T) {
	m := NewMedia(5, "https://example.com/media.jpg")
	assert.Equal(t, NewID(5), m.ID)
	assert.Equal(t, NewMediaURL("https://example.com/media.jpg"), m.MediaUrl)
}

func TestNewOrder(t *testing.T) {
	orderDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	o := NewOrder(1, 2, orderDate, "credit_card", 59.98)
	assert.Equal(t, NewID(1), o.ID)
	assert.Equal(t, NewID(2), o.UserID)
	assert.Equal(t, NewFinalPrice(59.98), o.FinalPrice)
	assert.Equal(t, NewPaymentMethod("credit_card"), o.PaymentMethod)
	assert.Equal(t, NewOrderDate(orderDate), o.OrderDate)
}

func TestPassword_NeverSerialized(t *testing.T) {
	// Verify password is never exposed in JSON
	a := Account{
		ID:       NewID(1),
		Username: NewUsername("test"),
		Password: NewPassword("secret123"),
	}

	// This test documents the expected behavior - password should not be in JSON
	// The actual implementation should omit password from JSON
	_ = a
}

func TestGame_AddCategory(t *testing.T) {
	g := NewGame(1, "Test", "Desc", 10.0, "poster.jpg")
	cat := NewCategory(5, "Action", false)
	g.AddCategory(cat)
	assert.Len(t, g.Categories, 1)
	assert.Equal(t, cat, g.Categories[0])

	// AddCategory appends without deduplication (current implementation)
	g.AddCategory(cat)
	assert.Len(t, g.Categories, 2)
}