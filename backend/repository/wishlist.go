package repository

import (
	"software-metrics-project/model"
)

type WishlistRepository struct {
	Model model.Wishlist
}

func (r *WishlistRepository) GetWishlists(userId int) ([]model.Wishlist, error) {
	rows, err := listCollectionRows(userId, wishlistCollection)
	if err != nil {
		return nil, err
	}
	entries := groupCollectionRows(rows)
	wishlists := make([]model.Wishlist, 0, len(entries))
	for _, e := range entries {
		wishlists = append(wishlists, model.Wishlist{
			UserID:       model.NewID(e.userID),
			ProductID:    model.NewID(e.productID),
			WishlistDate: model.NewDate(e.date),
			Game:         e.game,
		})
	}
	return wishlists, nil
}

func (r *WishlistRepository) AddWishlist(userID, gameID int) (*model.Wishlist, error) {
	if err := insertCollectionEntry(userID, gameID, "wishlists", "wishlist_date"); err != nil {
		return nil, err
	}
	return r.getWishlist(userID, gameID)
}

func (r *WishlistRepository) getWishlist(userID, gameID int) (*model.Wishlist, error) {
	rows, err := getCollectionRows(userID, gameID, wishlistCollection, true)
	if err != nil {
		return nil, err
	}
	if _, err := requireOneRow(rows); err != nil {
		return nil, err
	}
	e := groupCollectionRows(rows)[0]
	return &model.Wishlist{
		UserID:       model.NewID(e.userID),
		ProductID:    model.NewID(e.productID),
		WishlistDate: model.NewDate(e.date),
		Game:         e.game,
	}, nil
}

func (r *WishlistRepository) DeleteWishlist(userID, gameID int) (*model.Wishlist, error) {
	output, err := r.getWishlist(userID, gameID)
	if err != nil {
		return nil, err
	}
	if err := deleteCollectionEntry(userID, gameID, "wishlists"); err != nil {
		return nil, err
	}
	return output, nil
}
