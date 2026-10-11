package service

import (
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type OrderService struct {
	Repo *repository.OrderRepository
}

func (s *OrderService) Create(order model.Order) (*model.Order, error) {
	return s.Repo.Create(order)
}
