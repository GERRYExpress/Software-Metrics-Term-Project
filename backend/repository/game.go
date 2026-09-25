package repository

import (
	"software-metrics-project/config"
	"software-metrics-project/model"
	"time"

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

		gamePosition, ok := gameIndex[gameID]
		if !ok {
			games = append(games, model.Game{
				ID:          gameID,
				Name:        name,
				Description: description,
				Price:       price,
				PosterUrl:   posterURL,
				Categories:  make([]model.Category, 0),
			})
			gamePosition = len(games) - 1
			gameIndex[gameID] = gamePosition
		}

		if categoryID != nil {
			games[gamePosition].Categories = append(games[gamePosition].Categories, model.Category{
				ID:       *categoryID,
				Name:     *categoryName,
				IsHidden: *isHidden,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return games, nil
}

func (r *GameRepository) GetGame(id int) (*model.GameDetail, error) {
	mysqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	query, args, err := mysqlBuilder.
		Select("p.product_id", "p.name", "p.description", "p.price", "p.poster_url", "p.release_date",
			"c.category_id", "c.name", "a.user_id", "a.username", "hr.requirement_id", "hr.operating_system", "hr.processor", "hr.memory", "hr.graphics_card", "hr.directx_version",
			"hr.network", "hr.storage", "hr.soundcard", "hr.additional_note", "hr.type",
			"m.media_id", "m.media_url").
		From("products p").
		LeftJoin("product_categories pc ON pc.product_id = p.product_id").
		LeftJoin("categories c ON pc.category_id = c.category_id").
		LeftJoin("accounts a ON a.user_id = p.user_id").
		LeftJoin("hardware_requirement hr ON hr.product_id = p.product_id").
		LeftJoin("media m ON m.media_id = p.product_id").
		Where(squirrel.Eq{"c.is_hidden": false, "p.product_id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var output model.GameDetail
	initialized := false
	for rows.Next() {
		var (
			gameID          int
			name            string
			description     string
			price           float64
			posterUrl       string
			releaseDate     time.Time
			categoryID      *int
			categoryName    *string
			userID          *int
			username        *string
			requirementID   *int
			operatingSystem *string
			processor       *string
			memory          *string
			graphicsCard    *string
			directXVersion  *string
			network         *string
			storage         *string
			soundcard       *string
			additionalNote  *string
			reqType         *string
			mediaID         *int
			mediaUrl        *string
		)
		if err := rows.Scan(&gameID, &name, &description, &price, &posterUrl, &releaseDate, &categoryID,
			&categoryName, &userID, &username, &requirementID, &operatingSystem, &processor, &memory, &graphicsCard, &directXVersion,
			&network, &storage, &soundcard, &additionalNote, &reqType, &mediaID, &mediaUrl); err != nil {
			return nil, err
		}
		if !initialized {
			output = model.GameDetail{
				ID:                  gameID,
				Name:                name,
				Description:         description,
				Price:               price,
				PosterUrl:           posterUrl,
				ReleaseDate:         releaseDate,
				Account:             model.Account{},
				Categories:          make([]model.Category, 0),
				HardwareRequirement: make([]model.Hardware, 0),
				Media:               make([]model.Media, 0),
			}
			initialized = true
		}

		if categoryID != nil {
			output.Categories = append(output.Categories, model.Category{
				ID:   *categoryID,
				Name: *categoryName,
			})
		}

		if userID != nil {
			output.Account = model.Account{
				ID:       *userID,
				Username: *username,
			}
		}
		if requirementID != nil {
			output.HardwareRequirement = append(output.HardwareRequirement, model.Hardware{
				ID:             *requirementID,
				OS:             *operatingSystem,
				Processor:      *processor,
				Memory:         *memory,
				GraphicsCard:   *graphicsCard,
				DirectXVersion: *directXVersion,
				Network:        *network,
				Storage:        *storage,
				Soundcard:      *soundcard,
				AdditionalNote: *additionalNote,
				Type:           *reqType,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &output, nil
}
