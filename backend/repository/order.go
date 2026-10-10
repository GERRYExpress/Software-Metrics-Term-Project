package repository

import (
	"software-metrics-project/config"
	"software-metrics-project/model"
	"time"

	"github.com/Masterminds/squirrel"
)

type OrderRepository struct {
	Model *model.Order
}

func (r *OrderRepository) Create(userID int, paymentMethod string, finalPrice float64) (*model.Order, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert("orders").
		Columns("user_id", "order_date", "payment_method", "final_price").
		Values(userID, squirrel.Expr("NOW()"), paymentMethod, finalPrice).
		ToSql()
	if err != nil {
		return nil, err
	}

	result, err := config.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.getOrder(int(orderID))
}

func (r *OrderRepository) getOrder(orderID int) (*model.Order, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(orderIDColumn, orderDateColumn, paymentMethodColumn, finalPriceColumn).
		From("orders o").
		Where(squirrel.Eq{"o.order_id": orderID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		id            int
		orderDate     model.OrderDate
		paymentMethod model.PaymentMethod
		finalPrice    model.FinalPrice
	)
	if err := config.DB.QueryRow(query, args...).Scan(&id, &orderDate, &paymentMethod, &finalPrice); err != nil {
		return nil, err
	}

	order := model.NewOrder(int(id), time.Time(orderDate), string(paymentMethod), float64(finalPrice))
	return &order, nil
}
