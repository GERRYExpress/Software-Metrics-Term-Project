package service

import (
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type CartService struct {
	Repo *repository.CartRepository
}

func (s *CartService) GetCarts(userID int) ([]model.Cart, error) {
	return s.Repo.GetCarts(userID)
}

func (s *CartService) AddCart(userID, gameID int) (*model.Cart, error) {
	return s.Repo.AddCart(userID, gameID)
}

func (s *CartService) DeleteCart(userID, gameID int) (*model.Cart, error) {
	return s.Repo.DeleteCart(userID, gameID)
}
