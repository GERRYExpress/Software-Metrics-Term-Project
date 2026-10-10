package repository

import (
	"database/sql"
	"software-metrics-project/config"
	"software-metrics-project/model"
	"time"

	"github.com/Masterminds/squirrel"
)

type CartRepository struct {
	Cart model.Cart
}

type cartFields struct {
	userID        string
	cartGameID    string
	cartDate      string
	gameID        string
	gameName      string
	gamePrice     string
	gamePosterUrl string
	categoryID    string
	categoryName  string
}

func (r *CartRepository) GetCarts(userID int) ([]model.Cart, error) {
	mysqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	query, args, err := mysqlBuilder.
		Select(cartUserIDColumn, cartProductIDColumn, cartDateColumn, "p.name", "p.price", "p.poster_url", "c.category_id", "c.name").
		From("carts ca").
		LeftJoin("products p on ca.product_id = p.product_id").
		LeftJoin("product_categories pc ON p.product_id = pc.product_id").
		LeftJoin("categories c ON c.category_id = pc.category_id AND c.is_hidden = false").
		Where(squirrel.Eq{cartUserIDColumn: userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carts := make([]model.Cart, 0)
	cartIndex := make(map[int]int)
	for rows.Next() {
		var (
			cartUserID   int
			productID    int
			cartDate     time.Time
			gameName     *string
			gamePrice    *float64
			posterUrl    *string
			categoryID   *int
			categoryName *string
		)

		if err := rows.Scan(&cartUserID, &productID, &cartDate, &gameName, &gamePrice,
			&posterUrl, &categoryID, &categoryName); err != nil {
			return nil, err
		}

		cartPosition, ok := cartIndex[productID]
		if !ok {
			cart := model.Cart{
				UserID:    model.NewID(cartUserID),
				ProductID: model.NewID(productID),
				CartDate:  model.NewDate(cartDate),
				Game:      model.Game{},
			}
			if gameName != nil {
				cart.Game = model.NewGame(productID, *gameName, "", *gamePrice, *posterUrl)
			}
			carts = append(carts, cart)
			cartPosition = len(carts) - 1
			cartIndex[productID] = cartPosition
		}

		if categoryID != nil {
			carts[cartPosition].Game.Categories = append(carts[cartPosition].Game.Categories,
				model.NewCategory(*categoryID, *categoryName, false))
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return carts, nil
}

func (r *CartRepository) AddCart(userID, gameID int) (*model.Cart, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert("carts").
		Columns("user_id", "product_id", "cart_date").
		Values(userID, gameID, time.Now()).
		Suffix("ON DUPLICATE KEY UPDATE cart_date = cart_date").
		ToSql()

	if err != nil {
		return nil, err
	}

	_, err = config.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return r.getCart(userID, gameID)
}

func (r *CartRepository) getCart(userID, gameID int) (*model.Cart, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(cartUserIDColumn, cartProductIDColumn, cartDateColumn, "p.name", "p.price", "p.poster_url", "c.category_id", "c.name").
		From("carts ca").
		Join("products p ON ca.product_id = p.product_id").
		LeftJoin("product_categories pc ON p.product_id = pc.product_id").
		LeftJoin("categories c ON c.category_id = pc.category_id AND c.is_hidden = false").
		Where(squirrel.Eq{cartUserIDColumn: userID, cartProductIDColumn: gameID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var output model.Cart
	initialized := false
	for rows.Next() {
		var (
			cartUserID   int
			productID    int
			cartDate     time.Time
			gameName     string
			gamePrice    float64
			posterURL    string
			categoryID   *int
			categoryName *string
		)

		if err := rows.Scan(&cartUserID, &productID, &cartDate, &gameName, &gamePrice,
			&posterURL, &categoryID, &categoryName); err != nil {
			return nil, err
		}

		if !initialized {
			output = model.Cart{
				UserID:    model.NewID(cartUserID),
				ProductID: model.NewID(productID),
				CartDate:  model.NewDate(cartDate),
				Game:      model.NewGame(gameID, gameName, "", gamePrice, posterURL),
			}
			initialized = true
		}

		if categoryID != nil {
			output.Game.Categories = append(output.Game.Categories,
				model.NewCategory(*categoryID, *categoryName, false))
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !initialized {
		return nil, sql.ErrNoRows
	}

	return &output, nil
}

func (r *CartRepository) DeleteCart(userID, gameID int) (*model.Cart, error) {
	output, err := r.getCart(userID, gameID)
	if err != nil {
		return nil, err
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Delete("carts").
		Where(squirrel.Eq{cartUserFilter: userID, cartProductFilter: gameID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	if _, err := config.DB.Exec(query, args...); err != nil {
		return nil, err
	}

	return output, nil
}
