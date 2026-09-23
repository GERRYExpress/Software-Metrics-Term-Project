package repository

import (
	"software-metrics-project/config"
	"software-metrics-project/model"

	"github.com/Masterminds/squirrel"
)

type GameRepository struct {
	Model model.Game
}

func (r *GameRepository) GetGames() ([]model.Game, error) {
	query, args, err := squirrel.
		Select(
			"p.product_id", "p.name", "p.description", "p.price", "p.poster_url",
			"c.category_id", "c.name", "c.is_hidden",
		).
		From("products p").
		LeftJoin("product_categories pc ON pc.product_id = p.product_id").
		LeftJoin("categories c ON c.category_id = pc.category_id").
		OrderBy("p.product_id", "c.category_id").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]model.Game, 0)
	gameIndex := make(map[int]int)
	for rows.Next() {
		var (
			gameID       int
			name         string
			description  string
			price        float64
			posterURL    string
			categoryID   *int
			categoryName *string
			isHidden     *bool
		)

		if err := rows.Scan(&gameID, &name, &description, &price, &posterURL,
			&categoryID, &categoryName, &isHidden); err != nil {
			return nil, err
		}

		_, ok := gameIndex[gameID]
		if !ok {
			games = append(games, model.Game{
				ID:          gameID,
				Name:        name,
				Description: description,
				Price:       price,
				PosterUrl:   posterURL,
				Categories:  make([]model.Category, 0),
			})
			gameIndex[gameID] = gameID
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return games, nil
}
