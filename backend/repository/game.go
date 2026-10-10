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

type gameFields struct {
	productID       string
	name            string
	description     string
	price           string
	posterURL       string
	releaseDate     string
	categoryID      string
	categoryName    string
	categoryHidden  string
	userID          string
	username        string
	requirementID   string
	operatingSystem string
	processor       string
	memory          string
	graphicsCard    string
	directXVersion  string
	network         string
	storage         string
	soundcard       string
	additionalNote  string
	requirementType string
	mediaID         string
	mediaURL        string
}

var gameQueryFields = gameFields{
	productID:       "product_id",
	name:            "name",
	description:     "description",
	price:           "price",
	posterURL:       "poster_url",
	releaseDate:     "release_date",
	categoryID:      "category_id",
	categoryName:    "name",
	categoryHidden:  "is_hidden",
	userID:          "user_id",
	username:        "username",
	requirementID:   "requirement_id",
	operatingSystem: "operating_system",
	processor:       "processor",
	memory:          "memory",
	graphicsCard:    "graphics_card",
	directXVersion:  "directx_version",
	network:         "network",
	storage:         "storage",
	soundcard:       "soundcard",
	additionalNote:  "additional_note",
	requirementType: "type",
	mediaID:         "media_id",
	mediaURL:        "media_url",
}

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

func (f gameFields) productColumns(alias string) []string {
	return qualify(alias, f.productID, f.name, f.description, f.price, f.posterURL)
}

func (f gameFields) listColumns() []string {
	return append(f.productColumns("p"), qualify("c", f.categoryID, f.categoryName, f.categoryHidden)...)
}

func (f gameFields) detailColumns() []string {
	columns := append(
		qualify("p", f.productID, f.name, f.description, f.price, f.posterURL, f.releaseDate),
		qualify("c", f.categoryID, f.categoryName)...,
	)
	columns = append(columns, qualify("a", f.userID, f.username)...)
	columns = append(columns, qualify("hr",
		f.requirementID, f.operatingSystem, f.processor, f.memory, f.graphicsCard,
		f.directXVersion, f.network, f.storage, f.soundcard, f.additionalNote, f.requirementType,
	)...)
	return append(columns, qualify("m", f.mediaID, f.mediaURL)...)
}

func pagedProducts(keyword, category string, limit, page int) squirrel.SelectBuilder {
	query := squirrel.
		Select(gameQueryFields.productColumns("")...).
		From("products").
		OrderBy("product_id").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit))
	if keyword != "" {
		query = query.Where(squirrel.Like{"name": "%" + keyword + "%"})
	}
	if category != "" {
		query = query.Where(squirrel.Expr(
			"EXISTS (SELECT 1 FROM product_categories pc JOIN categories c ON c.category_id = pc.category_id WHERE pc.product_id = products.product_id AND c.name = ?)",
			category,
		))
	}
	return query
}

func gamesQuery(products squirrel.SelectBuilder) squirrel.SelectBuilder {
	return squirrel.
		Select(gameQueryFields.listColumns()...).
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
		queryBuilder = queryBuilder.Where(squirrel.Expr(
			"EXISTS (SELECT 1 FROM product_categories pc JOIN categories c ON c.category_id = pc.category_id WHERE pc.product_id = products.product_id AND c.name = ?)",
			category,
		))
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
		Select(gameQueryFields.detailColumns()...).
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
				model.NewHardware(*requirementID, *operatingSystem, *processor, *memory,
					*graphicsCard, *directXVersion, *network, *storage, *soundcard,
					*additionalNote, *reqType))
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &output, nil
}
