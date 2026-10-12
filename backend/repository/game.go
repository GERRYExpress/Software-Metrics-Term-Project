package repository

import (
	"database/sql"
	"slices"
	"software-metrics-project/config"
	"software-metrics-project/model"
	"time"

	"github.com/Masterminds/squirrel"
)

type GameRepository struct {
	Model model.Game
}

var ErrGameNotFound = sql.ErrNoRows

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

func pagedProducts(keyword, category string, sellerID *int, limit, page int) squirrel.SelectBuilder {
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
	if sellerID != nil {
		query = query.Where(squirrel.Eq{userIDField: *sellerID})
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

func (r *GameRepository) GetGamesCount(keyword, category string, sellerID *int) (int, error) {
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
	if sellerID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{userIDField: *sellerID})
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

func (r *GameRepository) GetGames(category string, sellerID *int, limit, page int) ([]model.Game, error) {
	query, args, err := gamesQuery(pagedProducts("", category, sellerID, limit, page)).ToSql()
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

func (r *GameRepository) GetGamesByName(keyword, category string, sellerID *int, limit, page int) ([]model.Game, error) {
	query, args, err := gamesQuery(pagedProducts(keyword, category, sellerID, limit, page)).ToSql()

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
		LeftJoin("media m ON m.product_id = p.product_id").
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
	categorySeen := make(map[int]struct{})
	requirementSeen := make(map[int]struct{})
	mediaSeen := make(map[int]struct{})
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
			if _, exists := categorySeen[*categoryID]; exists {
				categoryID = nil
			} else {
				categorySeen[*categoryID] = struct{}{}
			}
		}
		if categoryID != nil {
			output.Categories = append(output.Categories, model.NewCategory(*categoryID, *categoryName, false))
		}

		if userID != nil {
			output.Account = model.NewAccount(*userID, *username)
		}
		if requirementID != nil {
			if _, exists := requirementSeen[*requirementID]; exists {
				requirementID = nil
			} else {
				requirementSeen[*requirementID] = struct{}{}
			}
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
		if mediaID != nil && mediaUrl != nil {
			if _, exists := mediaSeen[*mediaID]; exists {
				continue
			}
			mediaSeen[*mediaID] = struct{}{}
			output.Media = append(output.Media,
				model.NewMedia(*mediaID, *mediaUrl))
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	slices.Reverse(output.Media)
	return &output, nil
}

func (r *GameRepository) Create(game model.GameDetail) (*model.GameDetail, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	columns := []string{"name", "description", "price", "poster_url", "user_id"}
	values := []any{game.Name, game.Description, game.Price, game.PosterUrl, game.Account.ID}
	if !time.Time(game.ReleaseDate).IsZero() {
		columns = append(columns, "release_date")
		values = append(values, time.Time(game.ReleaseDate))
	}
	insert := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert("products").
		Columns(columns...).
		Values(values...)
	query, args, err := insert.ToSql()
	if err != nil {
		return nil, err
	}
	result, err := tx.Exec(query, args...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := insertGameRelations(tx, id, game); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	created := model.NewGameDetail(int(id), string(game.Name), string(game.Description), float64(game.Price), string(game.PosterUrl), time.Time(game.ReleaseDate))
	created.Account = game.Account
	created.Categories = game.Categories
	created.HardwareRequirement = game.HardwareRequirement
	created.Media = game.Media
	return &created, nil
}

func (r *GameRepository) Update(id int, game model.GameDetail) (*model.GameDetail, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(
		"UPDATE products SET name = ?, description = ?, release_date = ?, price = ?, poster_url = ?, user_id = ? WHERE product_id = ?",
		game.Name, game.Description, time.Time(game.ReleaseDate), game.Price, game.PosterUrl, game.Account.ID, id,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	var exists bool
	if err := tx.QueryRow("SELECT EXISTS (SELECT 1 FROM products WHERE product_id = ?)", id).Scan(&exists); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if !exists {
		_ = tx.Rollback()
		return nil, ErrGameNotFound
	}

	for _, query := range []string{
		"DELETE FROM product_categories WHERE product_id = ?",
		"DELETE FROM hardware_requirement WHERE product_id = ?",
		"DELETE FROM media WHERE product_id = ?",
	} {
		if _, err := tx.Exec(query, id); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}
	if err := insertGameRelations(tx, int64(id), game); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	updated := model.NewGameDetail(id, string(game.Name), string(game.Description), float64(game.Price), string(game.PosterUrl), time.Time(game.ReleaseDate))
	updated.Account = game.Account
	updated.Categories = game.Categories
	updated.HardwareRequirement = game.HardwareRequirement
	updated.Media = game.Media
	return &updated, nil
}

func insertGameRelations(tx *sql.Tx, productID int64, game model.GameDetail) error {
	for _, category := range game.Categories {
		if _, err := tx.Exec("INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)", productID, category.ID); err != nil {
			return err
		}
	}
	for _, requirement := range game.HardwareRequirement {
		if _, err := tx.Exec(`INSERT INTO hardware_requirement
			(operating_system, processor, memory, graphics_card, directx_version, network, storage, soundcard, additional_note, type, product_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			requirement.OS, requirement.Processor, requirement.Memory, requirement.GraphicsCard,
			requirement.DirectXVersion, requirement.Network, requirement.Storage, requirement.Soundcard,
			requirement.AdditionalNote, requirement.Type, productID,
		); err != nil {
			return err
		}
	}
	mediaSeen := make(map[string]struct{}, len(game.Media))
	for _, media := range game.Media {
		url := string(media.MediaUrl)
		if _, exists := mediaSeen[url]; exists {
			continue
		}
		mediaSeen[url] = struct{}{}
		if _, err := tx.Exec("INSERT INTO media (media_url, product_id) VALUES (?, ?)", media.MediaUrl, productID); err != nil {
			return err
		}
	}
	return nil
}
