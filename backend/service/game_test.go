package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"software-metrics-project/model"
)

type mockGameRepo struct {
	mockGames []model.Game
	mockGame  *model.GameDetail
	mockErr   error
	mockCount int

	gotKeyword  string
	gotCategory string
	gotSellerID *int
	gotLimit    int
	gotPage     int
	gotID       int
	gotUserID   int
	gotGame     model.GameDetail
}

func (m *mockGameRepo) GetGamesCount(keyword, category string, sellerID *int) (int, error) {
	m.gotKeyword = keyword
	m.gotCategory = category
	m.gotSellerID = sellerID
	return m.mockCount, m.mockErr
}

func (m *mockGameRepo) GetGames(category string, sellerID *int, limit, page int) ([]model.Game, error) {
	m.gotCategory = category
	m.gotSellerID = sellerID
	m.gotLimit = limit
	m.gotPage = page
	return m.mockGames, m.mockErr
}

func (m *mockGameRepo) GetGamesByName(keyword, category string, sellerID *int, limit, page int) ([]model.Game, error) {
	m.gotKeyword = keyword
	m.gotCategory = category
	m.gotSellerID = sellerID
	m.gotLimit = limit
	m.gotPage = page
	return m.mockGames, m.mockErr
}

func (m *mockGameRepo) GetGame(id int) (*model.GameDetail, error) {
	m.gotID = id
	return m.mockGame, m.mockErr
}

func (m *mockGameRepo) Create(game model.GameDetail) (*model.GameDetail, error) {
	m.gotGame = game
	return m.mockGame, m.mockErr
}

func (m *mockGameRepo) Update(id int, game model.GameDetail) (*model.GameDetail, error) {
	m.gotID = id
	m.gotGame = game
	return m.mockGame, m.mockErr
}

func (m *mockGameRepo) Delete(id, userID int) error {
	m.gotID = id
	m.gotUserID = userID
	return m.mockErr
}

func TestGameService(t *testing.T) {
	sellerID := 9
	expectedGames := []model.Game{{ID: 1}}
	expectedGame := &model.GameDetail{ID: 1}
	expectedErr := errors.New("repository error")
	input := model.GameDetail{ID: 10}

	tests := []struct {
		name string
		run  func(svc *GameService) error
	}{
		{
			name: "GetGamesCount success",
			run: func(svc *GameService) error {
				count, err := svc.GetGamesCount("search", "Action", &sellerID)
				assert.Equal(t, 5, count)
				return err
			},
		},
		{
			name: "GetGames success",
			run: func(svc *GameService) error {
				games, err := svc.GetGames("Action", &sellerID, 10, 2)
				assert.Equal(t, expectedGames, games)
				return err
			},
		},
		{
			name: "GetGamesByName success",
			run: func(svc *GameService) error {
				games, err := svc.GetGamesByName("search", "Action", &sellerID, 10, 2)
				assert.Equal(t, expectedGames, games)
				return err
			},
		},
		{
			name: "GetGame success",
			run: func(svc *GameService) error {
				game, err := svc.GetGame(1)
				assert.Equal(t, expectedGame, game)
				return err
			},
		},
		{
			name: "Create success",
			run: func(svc *GameService) error {
				game, err := svc.Create(input)
				assert.Equal(t, expectedGame, game)
				return err
			},
		},
		{
			name: "Update success",
			run: func(svc *GameService) error {
				game, err := svc.Update(1, input)
				assert.Equal(t, expectedGame, game)
				return err
			},
		},
		{
			name: "Delete success",
			run: func(svc *GameService) error {
				return svc.Delete(1, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockGameRepo{
				mockGames: expectedGames,
				mockGame:  expectedGame,
				mockCount: 5,
			}
			svc := &GameService{Repo: repo}

			assert.NoError(t, tt.run(svc))
		})
	}

	t.Run("all repository errors are passed through", func(t *testing.T) {
		repo := &mockGameRepo{mockErr: expectedErr}
		svc := &GameService{Repo: repo}

		_, err := svc.GetGamesCount("", "", nil)
		assert.ErrorIs(t, err, expectedErr)

		_, err = svc.GetGames("", nil, 10, 1)
		assert.ErrorIs(t, err, expectedErr)

		_, err = svc.GetGamesByName("", "", nil, 10, 1)
		assert.ErrorIs(t, err, expectedErr)

		_, err = svc.GetGame(1)
		assert.ErrorIs(t, err, expectedErr)

		_, err = svc.Create(model.GameDetail{})
		assert.ErrorIs(t, err, expectedErr)

		_, err = svc.Update(1, model.GameDetail{})
		assert.ErrorIs(t, err, expectedErr)

		err = svc.Delete(1, 2)
		assert.ErrorIs(t, err, expectedErr)
	})
}
