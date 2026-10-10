package service

import (
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type GameService struct {
	Repo *repository.GameRepository
}

func (s *GameService) GetGamesCount(keyword, category string) (int, error) {
	return s.Repo.GetGamesCount(keyword, category)
}

func (s *GameService) GetGames(category string, limit, page int) ([]model.Game, error) {
	return s.Repo.GetGames(category, limit, page)
}

func (s *GameService) GetGamesByName(keyword, category string, limit, page int) ([]model.Game, error) {
	return s.Repo.GetGamesByName(keyword, category, limit, page)
}

func (s *GameService) GetGame(id int) (*model.GameDetail, error) {
	return s.Repo.GetGame(id)
}
