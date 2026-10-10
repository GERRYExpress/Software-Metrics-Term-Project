package model

type Wishlist struct {
	UserID       ID   `json:"userId"`
	ProductID    ID   `json:"productId"`
	WishlistDate Date `json:"wishlistedDate"`
	Game         Game `json:"game"`
}
