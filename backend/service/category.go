package service

import (
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type CategoryService struct {
	Repo *repository.CategoryRepository
}

func (s *CategoryService) GetCategories() ([]model.Category, error) {
	return s.Repo.GetCategories()
}
