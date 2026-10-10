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

const (
	productIDField       = "product_id"
	nameField            = "name"
	descriptionField     = "description"
	priceField           = "price"
	posterURLField       = "poster_url"
	releaseDateField     = "release_date"
	categoryIDField      = "category_id"
	categoryHiddenField  = "is_hidden"
	userIDField          = "user_id"
	usernameField        = "username"
	requirementIDField   = "requirement_id"
	operatingSystemField = "operating_system"
	processorField       = "processor"
	memoryField          = "memory"
	graphicsCardField    = "graphics_card"
	directXVersionField  = "directx_version"
	networkField         = "network"
	storageField         = "storage"
	soundcardField       = "soundcard"
	additionalNoteField  = "additional_note"
	requirementTypeField = "type"
	mediaIDField         = "media_id"
	mediaURLField        = "media_url"
)

const existQuery = "EXISTS (SELECT 1 FROM product_categories pc JOIN categories c ON c.category_id = pc.category_id WHERE pc.product_id = products.product_id AND c.name = ?)"

func qualify(alias string, fields ...string) []string {
	columns := make([]string, len(fields))
	prefix := ""
	if alias != "" {
		prefix = alias + "."
	}
	for i, field := range fields {
		columns[i] = prefix + field
	}
	return columns
}

func productColumns(alias string) []string {
	return qualify(alias, productIDField, nameField, descriptionField, priceField, posterURLField)
}

func listColumns() []string {
	return append(productColumns("p"), qualify("c", categoryIDField, nameField, categoryHiddenField)...)
}

func detailColumns() []string {
	columns := append(
		qualify("p", productIDField, nameField, descriptionField, priceField, posterURLField, releaseDateField),
		qualify("c", categoryIDField, nameField)...,
	)
	columns = append(columns, qualify("a", userIDField, usernameField)...)
	columns = append(columns, qualify("hr",
		requirementIDField, operatingSystemField, processorField, memoryField, graphicsCardField,
		directXVersionField, networkField, storageField, soundcardField, additionalNoteField, requirementTypeField,
	)...)
	return append(columns, qualify("m", mediaIDField, mediaURLField)...)
}

func pagedProducts(keyword, category string, limit, page int) squirrel.SelectBuilder {
	query := squirrel.
		Select(productColumns("")...).
		From("products").
		OrderBy("product_id").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit))
	if keyword != "" {
		query = query.Where(squirrel.Like{"name": "%" + keyword + "%"})
	}
	if category != "" {
		query = query.Where(squirrel.Expr(existQuery, category))
	}
	return query
}

func gamesQuery(products squirrel.SelectBuilder) squirrel.SelectBuilder {
	return squirrel.
		Select(listColumns()...).
		FromSelect(products, "p").
		LeftJoin("product_categories pc ON pc.product_id = p.product_id").
		LeftJoin("categories c ON c.category_id = pc.category_id").
		OrderBy("p.product_id", "c.category_id").
		PlaceholderFormat(squirrel.Question)
}

func (r *GameRepository) GetGamesCount(keyword, category string) (int, error) {
	queryBuilder := squirrel.
		Select("COUNT(*)").
		From("products").
		PlaceholderFormat(squirrel.Question)
	if keyword != "" {
		queryBuilder = queryBuilder.Where(squirrel.Like{"name": "%" + keyword + "%"})
	}
	if category != "" {
		queryBuilder = queryBuilder.Where(squirrel.Expr(existQuery, category))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	if err := config.DB.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GameRepository) GetGames(category string, limit, page int) ([]model.Game, error) {
	query, args, err := gamesQuery(pagedProducts("", category, limit, page)).ToSql()
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
			games = append(games, model.NewGame(gameID, name, description, price, posterURL))
			gamePosition = len(games) - 1
			gameIndex[gameID] = gamePosition
		}

		if categoryID != nil {
			games[gamePosition].Categories = append(games[gamePosition].Categories,
				model.NewCategory(*categoryID, *categoryName, *isHidden))
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return games, nil
}

func (r *GameRepository) GetGamesByName(keyword, category string, limit, page int) ([]model.Game, error) {
	query, args, err := gamesQuery(pagedProducts(keyword, category, limit, page)).ToSql()

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
			games = append(games, model.NewGame(gameID, name, description, price, posterURL))
			gamePosition = len(games) - 1
			gameIndex[gameID] = gamePosition
		}

		if categoryID != nil {
			games[gamePosition].Categories = append(games[gamePosition].Categories,
				model.NewCategory(*categoryID, *categoryName, *isHidden))
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
		Select(detailColumns()...).
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
			output = model.NewGameDetail(gameID, name, description, price, posterUrl, releaseDate)
			initialized = true
		}

		if categoryID != nil {
			output.Categories = append(output.Categories, model.NewCategory(*categoryID, *categoryName, false))
		}

		if userID != nil {
			output.Account = model.NewAccount(*userID, *username)
		}
		if requirementID != nil {
			output.HardwareRequirement = append(output.HardwareRequirement,
				model.NewHardware(model.HardwareData{
					ID: *requirementID, OS: *operatingSystem, Processor: *processor,
					Memory: *memory, GraphicsCard: *graphicsCard,
					DirectXVersion: *directXVersion, Network: *network, Storage: *storage,
					Soundcard: *soundcard, AdditionalNote: *additionalNote, Type: *reqType,
				}))
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &output, nil
}
