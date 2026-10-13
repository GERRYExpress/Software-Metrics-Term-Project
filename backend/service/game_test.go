package service

import (
	"software-metrics-project/model"
	"testing"
)

type mockGameRepo struct {
	mockGames []model.Game
	mockGame  *model.GameDetail
	mockErr   error
	mockCount int
}

func (m *mockGameRepo) GetGamesCount(keyword, category string, sellerID *int) (int, error) {
	return m.mockCount, m.mockErr
}
func (m *mockGameRepo) GetGames(category string, sellerID *int, limit, page int) ([]model.Game, error) {
	return m.mockGames, m.mockErr
}
func (m *mockGameRepo) GetGamesByName(keyword, category string, sellerID *int, limit, page int) ([]model.Game, error) {
	return m.mockGames, m.mockErr
}
func (m *mockGameRepo) GetGame(id int) (*model.GameDetail, error) {
	return m.mockGame, m.mockErr
}
func (m *mockGameRepo) Create(game model.GameDetail) (*model.GameDetail, error) {
	return m.mockGame, m.mockErr
}
func (m *mockGameRepo) Update(id int, game model.GameDetail) (*model.GameDetail, error) {
	return m.mockGame, m.mockErr
}
func (m *mockGameRepo) Delete(id, userID int) error {
	return m.mockErr
}

func TestGameService_GetGame(t *testing.T) {
	expectedGame := &model.GameDetail{ID: 1}
	repo := &mockGameRepo{mockGame: expectedGame}
	svc := &GameService{Repo: repo}

	game, err := svc.GetGame(1)
	if err != nil || game.ID != 1 {
		t.Errorf("failed to get game")
	}
}

func TestGameService_GetGamesCount(t *testing.T) {
	repo := &mockGameRepo{mockCount: 5}
	svc := &GameService{Repo: repo}

	count, err := svc.GetGamesCount("", "", nil)
	if err != nil || count != 5 {
		t.Errorf("failed to get game count, got %v", count)
	}
}
