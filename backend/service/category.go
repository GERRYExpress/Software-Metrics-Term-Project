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

func (s *CategoryService) Create(name string) (model.Category, error) {
	return s.Repo.Create(name)
}

func (s *CategoryService) Update(id int, name string, isHidden bool) (model.Category, error) {
	return s.Repo.Update(id, name, isHidden)
}

func (s *CategoryService) Delete(id int) error {
	return s.Repo.Delete(id)
}
