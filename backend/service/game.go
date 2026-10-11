package service

import (
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type GameService struct {
	Repo *repository.GameRepository
}

func (s *GameService) GetGamesCount(keyword, category string, sellerID *int) (int, error) {
	return s.Repo.GetGamesCount(keyword, category, sellerID)
}

func (s *GameService) GetGames(category string, sellerID *int, limit, page int) ([]model.Game, error) {
	return s.Repo.GetGames(category, sellerID, limit, page)
}

func (s *GameService) GetGamesByName(keyword, category string, sellerID *int, limit, page int) ([]model.Game, error) {
	return s.Repo.GetGamesByName(keyword, category, sellerID, limit, page)
}

func (s *GameService) GetGame(id int) (*model.GameDetail, error) {
	return s.Repo.GetGame(id)
}
