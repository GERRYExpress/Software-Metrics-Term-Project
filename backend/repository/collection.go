package repository

import (
	"database/sql"
	"software-metrics-project/config"
	"software-metrics-project/model"
	"time"

	"github.com/Masterminds/squirrel"
)

type collectionConfig struct {
	table      string
	from       string
	userCol    string
	productCol string
	dateCol    string
	productOn  string
}

var cartCollection = collectionConfig{
	table: "carts", from: "carts ca",
	userCol: cartUserIDColumn, productCol: cartProductIDColumn, dateCol: cartDateColumn,
	productOn: "ca.product_id = p.product_id",
}

var wishlistCollection = collectionConfig{
	table: "wishlists", from: "wishlists w",
	userCol: wishlistUserIDColumn, productCol: wishlistProductIDColumn, dateCol: wishlistDateColumn,
	productOn: "w.product_id = p.product_id",
}

type collectionRow struct {
	userID       int
	productID    int
	date         time.Time
	gameName     *string
	gamePrice    *float64
	posterURL    *string
	categoryID   *int
	categoryName *string
}

func listCollectionRows(userID int, cfg collectionConfig) ([]collectionRow, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(cfg.userCol, cfg.productCol, cfg.dateCol, "p.name", "p.price", "p.poster_url", "c.category_id", "c.name").
		From(cfg.from).
		LeftJoin("products p on " + cfg.productOn).
		LeftJoin("product_categories pc ON p.product_id = pc.product_id").
		LeftJoin("categories c ON c.category_id = pc.category_id AND c.is_hidden = false").
		Where(squirrel.Eq{cfg.userCol: userID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	return scanCollectionRows(query, args)
}

func getCollectionRows(userID, gameID int, cfg collectionConfig, innerJoin bool) ([]collectionRow, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(cfg.userCol, cfg.productCol, cfg.dateCol, "p.name", "p.price", "p.poster_url", "c.category_id", "c.name").
		From(cfg.from)
	if innerJoin {
		builder = builder.Join("products p ON " + cfg.productOn)
	} else {
		builder = builder.LeftJoin("products p on " + cfg.productOn)
	}
	query, args, err := builder.
		LeftJoin("product_categories pc ON p.product_id = pc.product_id").
		LeftJoin("categories c ON c.category_id = pc.category_id AND c.is_hidden = false").
		Where(squirrel.Eq{cfg.userCol: userID, cfg.productCol: gameID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	return scanCollectionRows(query, args)
}

func scanCollectionRows(query string, args []any) ([]collectionRow, error) {
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]collectionRow, 0)
	for rows.Next() {
		var r collectionRow
		if err := rows.Scan(&r.userID, &r.productID, &r.date, &r.gameName, &r.gamePrice,
			&r.posterURL, &r.categoryID, &r.categoryName); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// groupedEntry is the shared shape behind model.Cart and model.Wishlist.
type groupedEntry struct {
	userID    int
	productID int
	date      time.Time
	game      model.Game
}

func buildGameFromRow(r collectionRow) model.Game {
	if r.gameName == nil {
		return model.Game{}
	}
	price, poster := 0.0, ""
	if r.gamePrice != nil {
		price = *r.gamePrice
	}
	if r.posterURL != nil {
		poster = *r.posterURL
	}
	return model.NewGame(r.productID, *r.gameName, "", price, poster)
}

func groupCollectionRows(rows []collectionRow) []groupedEntry {
	entries := make([]groupedEntry, 0, len(rows))
	index := make(map[int]int)
	for _, r := range rows {
		pos, ok := index[r.productID]
		if !ok {
			game := buildGameFromRow(r)
			entries = append(entries, groupedEntry{userID: r.userID, productID: r.productID, date: r.date, game: game})
			pos = len(entries) - 1
			index[r.productID] = pos
		}
		if r.categoryID != nil {
			entries[pos].game.Categories = append(entries[pos].game.Categories,
				model.NewCategory(*r.categoryID, *r.categoryName, false))
		}
	}
	return entries
}

func insertCollectionEntry(userID, gameID int, table, dateColumn string) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert(table).
		Columns("user_id", "product_id", dateColumn).
		Values(userID, gameID, time.Now()).
		Suffix("ON DUPLICATE KEY UPDATE " + dateColumn + " = " + dateColumn).
		ToSql()
	if err != nil {
		return err
	}
	_, err = config.DB.Exec(query, args...)
	return err
}

func deleteCollectionEntry(userID, gameID int, table string) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Delete(table).
		Where(squirrel.Eq{"user_id": userID, "product_id": gameID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = config.DB.Exec(query, args...)
	return err
}

func requireOneRow(rows []collectionRow) ([]collectionRow, error) {
	if len(rows) == 0 {
		return nil, sql.ErrNoRows
	}
	return rows, nil
}
