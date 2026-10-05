package service

import (
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type WishlistService struct {
	Repo *repository.WishlistRepository
}

func (s *WishlistService) GetWishlists(userID int) ([]model.Wishlist, error) {
	return s.Repo.GetWishlists(userID)
}

func (s *WishlistService) AddWishlist(userID, gameID int) (*model.Wishlist, error) {
	return s.Repo.AddWishlist(userID, gameID)
}

func (s *WishlistService) DeleteWishlist(userID, gameID int) (*model.Wishlist, error) {
	return s.Repo.DeleteWishlist(userID, gameID)
}
