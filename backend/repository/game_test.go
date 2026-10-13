package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/config"
	"software-metrics-project/model"
	"software-metrics-project/repository/repotestutil"
)

func setupGameRepoTest(t *testing.T) (*GameRepository, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	originalDB := config.DB
	config.DB = db

	repo := &GameRepository{}
	return repo, mock, func() {
		config.DB = originalDB
		cleanup()
	}
}

// --- Pure Function / SQL Builder Tests ---

func Test_qualify(t *testing.T) {
	t.Run("with alias", func(t *testing.T) {
		got := qualify("p", "id", "name")
		want := []string{"p.id", "p.name"}
		assert.Equal(t, want, got)
	})

	t.Run("without alias", func(t *testing.T) {
		got := qualify("", "id", "name")
		want := []string{"id", "name"}
		assert.Equal(t, want, got)
	})
}

func Test_productColumns(t *testing.T) {
	got := productColumns("p")
	assert.Contains(t, got, "p.product_id")
	assert.Contains(t, got, "p.name")
}

func Test_pagedProducts(t *testing.T) {
	sellerID := 42
	builder := pagedProducts("zelda", "RPG", &sellerID, 10, 2)

	sqlStr, args, err := builder.ToSql()
	assert.NoError(t, err)
	assert.Contains(t, sqlStr, "LIMIT 10 OFFSET 10")
	assert.Contains(t, sqlStr, "name LIKE ?")
	assert.Contains(t, sqlStr, "user_id = ?")
	assert.Len(t, args, 3) // %zelda%, RPG, 42
}

// --- Database Mock Tests ---

func TestGameRepository_GetGamesCount(t *testing.T) {
	repo, mock, cleanup := setupGameRepoTest(t)
	defer cleanup()

	t.Run("success without filters", func(t *testing.T) {
		mock.ExpectQuery("^SELECT COUNT\\(\\*\\) FROM products").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(15))

		count, err := repo.GetGamesCount("", "", nil)
		assert.NoError(t, err)
		assert.Equal(t, 15, count)
	})

	t.Run("success with filters", func(t *testing.T) {
		seller := 1
		// Squirrel generates predictable parameterized queries
		mock.ExpectQuery("^SELECT(.*)FROM products WHERE name LIKE(.*)AND EXISTS(.*)AND user_id").
			WithArgs("%test%", "Action", 1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		count, err := repo.GetGamesCount("test", "Action", &seller)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})
}

func TestGameRepository_GetGame(t *testing.T) {
	repo, mock, cleanup := setupGameRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"product_id", "name", "description", "price", "poster_url", "release_date",
			"category_id", "category_name", "user_id", "username",
			"requirement_id", "operating_system", "processor", "memory", "graphics_card",
			"directx_version", "network", "storage", "soundcard", "additional_note", "type",
			"media_id", "media_url",
		}).AddRow(
			1, "Game 1", "Desc", 59.99, "url", now,
			10, "RPG", 42, "seller1",
			100, "Win10", "i7", "16GB", "RTX", "12", "Broadband", "50GB", "Yes", "", "Minimum",
			200, "http://img.com",
		)

		mock.ExpectQuery("^SELECT(.*)FROM products p").WithArgs(1).WillReturnRows(rows)

		game, err := repo.GetGame(1)

		assert.NoError(t, err)
		assert.NotNil(t, game)
		assert.Equal(t, 1, int(game.ID))
		assert.Equal(t, "Game 1", string(game.Name))
		assert.Len(t, game.Categories, 1)
		assert.Len(t, game.HardwareRequirement, 1)
		assert.Len(t, game.Media, 1)
	})
}

func TestGameRepository_Delete(t *testing.T) {
	repo, mock, cleanup := setupGameRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("^DELETE FROM products").
			WithArgs(1, 42).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1, 42)
		assert.NoError(t, err)
	})

	t.Run("not owned by user", func(t *testing.T) {
		mock.ExpectExec("^DELETE FROM products").
			WithArgs(1, 99). // Wrong user tries to delete
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Fallback check to see if the game exists at all
		mock.ExpectQuery("^SELECT EXISTS").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err := repo.Delete(1, 99)
		assert.ErrorIs(t, err, ErrGameNotOwned)
	})
}

func TestGameRepository_Create(t *testing.T) {
	repo, mock, cleanup := setupGameRepoTest(t)
	defer cleanup()

	t.Run("success with transactions", func(t *testing.T) {
		game := model.GameDetail{
			Name:       "New Game",
			Account:    model.Account{ID: 42},
			Categories: []model.Category{{ID: 1}},
			Media:      []model.Media{{MediaUrl: "url"}},
		}

		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO products").WillReturnResult(sqlmock.NewResult(10, 1))

		// Sub-inserts via insertGameRelationsTx
		mock.ExpectExec("^INSERT INTO product_categories").WithArgs(10, 1).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("^INSERT INTO media").WithArgs("url", 10).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		created, err := repo.Create(game)
		assert.NoError(t, err)
		assert.Equal(t, 10, int(created.ID))
	})
}
