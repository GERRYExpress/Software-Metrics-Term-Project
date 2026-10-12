package service

import (
	"software-metrics-project/model"
)

type GameService struct {
	Repo GameRepository
}

type GameRepository interface {
	GetGamesCount(string, string, *int) (int, error)
	GetGames(string, *int, int, int) ([]model.Game, error)
	GetGamesByName(string, string, *int, int, int) ([]model.Game, error)
	GetGame(int) (*model.GameDetail, error)
	Create(model.GameDetail) (*model.GameDetail, error)
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

func (s *GameService) Create(game model.GameDetail) (*model.GameDetail, error) {
	return s.Repo.Create(game)
}
