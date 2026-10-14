package repotestutil

import (
	"time"

	"software-metrics-project/model"
)

// Category creates a test category with optional modifiers
func Category(name string, opts ...func(*model.Category)) model.Category {
	c := model.NewCategory(0, name, false)
	for _, o := range opts {
		o(&c)
	}
	return c
}

// WithCategoryID sets the category ID
func WithCategoryID(id int) func(*model.Category) {
	return func(c *model.Category) { c.ID = model.NewID(id) }
}

// WithHidden sets the hidden flag
func WithHidden(hidden bool) func(*model.Category) {
	return func(c *model.Category) { c.IsHidden = model.NewHidden(hidden) }
}

// Game creates a test game with optional modifiers
func Game(title string, opts ...func(*model.GameDetail)) *model.GameDetail {
	g := model.NewGameDetail(0, title, "Test description", 9.99, "poster.jpg", time.Now())
	for _, o := range opts {
		o(&g)
	}
	return &g
}

// WithGameID sets the game ID
func WithGameID(id int) func(*model.GameDetail) {
	return func(g *model.GameDetail) { g.ID = model.NewID(id) }
}

// WithPrice sets the price
func WithPrice(price float64) func(*model.GameDetail) {
	return func(g *model.GameDetail) { g.Price = model.NewPrice(price) }
}

// WithCategories adds categories
func WithCategories(cats ...model.Category) func(*model.GameDetail) {
	return func(g *model.GameDetail) { g.Categories = cats }
}

// Account creates a test account
func Account(username string, opts ...func(*model.Account)) model.Account {
	a := model.NewAccount(0, username)
	for _, o := range opts {
		o(&a)
	}
	return a
}

// WithAccountID sets the account ID
func WithAccountID(id int) func(*model.Account) {
	return func(a *model.Account) { a.ID = model.NewID(id) }
}

// WithSeller sets the seller flag
func WithSeller(isSeller bool) func(*model.Account) {
	return func(a *model.Account) { a.IsSeller = model.NewHidden(isSeller) }
}

// WithBuyer sets the buyer flag
func WithBuyer(isBuyer bool) func(*model.Account) {
	return func(a *model.Account) { a.IsBuyer = model.NewHidden(isBuyer) }
}