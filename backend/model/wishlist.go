package model

import "time"

type Wishlist struct {
	UserID       int       `json:"userId"`
	ProductID    int       `json:"productId"`
	WishlistDate time.Time `json:"wishlistedDate"`
	Game         Game      `json:"game"`
}
