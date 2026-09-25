package service

import (
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type GameService struct {
	Repo *repository.GameRepository
}

func (s *GameService) GetGames() ([]model.Game, error) {
	return s.Repo.GetGames()
}

func (s *GameService) GetGame(id int) (*model.GameDetail, error) {
	return s.Repo.GetGame(id)
}
