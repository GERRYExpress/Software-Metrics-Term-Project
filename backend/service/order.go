package service

import (
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type OrderService struct {
	Repo *repository.OrderRepository
}

func (s *OrderService) Create(userID int, paymentMethod string, finalPrice float64) (*model.Order, error) {
	return s.Repo.Create(userID, paymentMethod, finalPrice)
}
