package repository

import (
	"database/sql"
	"software-metrics-project/config"
	"software-metrics-project/model"
	"time"

	"github.com/Masterminds/squirrel"
)

type OrderRepository struct {
	Model *model.Order
}

func (r *OrderRepository) Create(order model.Order) (*model.Order, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert("orders").
		Columns("user_id", "order_date", "payment_method", "final_price").
		Values(order.UserID, squirrel.Expr("NOW()"), order.PaymentMethod, order.FinalPrice).
		ToSql()
	if err != nil {
		return nil, err
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	for _, game := range order.Games {
		if _, err = tx.Exec("INSERT INTO order_items (order_id, product_id) VALUES (?, ?)", orderID, game.ID); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return r.getOrder(int(orderID))
}

func (r *OrderRepository) getOrder(orderID int) (*model.Order, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(append(
			[]string{orderIDColumn, orderUserIDColumn, orderDateColumn, paymentMethodColumn, finalPriceColumn},
			append(qualify("p", productIDField, nameField, descriptionField, priceField, posterURLField, releaseDateField),
				qualify("c", categoryIDField, nameField, categoryHiddenField)...)...,
		)...).
		From("orders o").
		Join("order_items oi ON oi.order_id = o.order_id").
		Join("products p ON p.product_id = oi.product_id").
		LeftJoin("product_categories pc ON pc.product_id = p.product_id").
		LeftJoin("categories c ON c.category_id = pc.category_id").
		Where(squirrel.Eq{"o.order_id": orderID}).
		OrderBy("p.product_id", "c.category_id").
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		id            int
		userID        int
		orderDate     time.Time
		paymentMethod string
		finalPrice    float64
	)

	var order *model.Order
	gameIndexes := make(map[int]int)
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			gameID       int
			gameName     string
			description  string
			price        float64
			posterURL    string
			releaseDate  time.Time
			categoryID   *int
			categoryName *string
			isHidden     *model.Hidden
		)
		if err := rows.Scan(&id, &userID, &orderDate, &paymentMethod, &finalPrice, &gameID, &gameName,
			&description, &price, &posterURL, &releaseDate, &categoryID, &categoryName, &isHidden); err != nil {
			return nil, err
		}
		if order == nil {
			created := model.NewOrder(id, userID, orderDate, paymentMethod, finalPrice)
			order = &created
		}
		index, exists := gameIndexes[gameID]
		if !exists {
			order.Games = append(order.Games, model.NewGame(gameID, gameName, description, price, posterURL))
			index = len(order.Games) - 1
			gameIndexes[gameID] = index
		}
		if categoryID != nil && categoryName != nil {
			order.Games[index].AddCategory(model.NewCategory(*categoryID, *categoryName, isHidden != nil && bool(*isHidden)))
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if order == nil {
		return nil, sql.ErrNoRows
	}
	return order, nil
}
