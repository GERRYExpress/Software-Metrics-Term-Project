package repository

import (
	"software-metrics-project/model"
)

type CartRepository struct {
	Cart model.Cart
}

func (r *CartRepository) GetCarts(userID int) ([]model.Cart, error) {
	rows, err := listCollectionRows(userID, cartCollection)
	if err != nil {
		return nil, err
	}
	entries := groupCollectionRows(rows)
	carts := make([]model.Cart, 0, len(entries))
	for _, e := range entries {
		carts = append(carts, model.Cart{
			UserID:    model.NewID(e.userID),
			ProductID: model.NewID(e.productID),
			CartDate:  model.NewDate(e.date),
			Game:      e.game,
		})
	}
	return carts, nil
}

func (r *CartRepository) AddCart(userID, gameID int) (*model.Cart, error) {
	if err := insertCollectionEntry(userID, gameID, "carts", "cart_date"); err != nil {
		return nil, err
	}
	return r.getCart(userID, gameID)
}

func (r *CartRepository) getCart(userID, gameID int) (*model.Cart, error) {
	rows, err := getCollectionRows(userID, gameID, cartCollection, true)
	if err != nil {
		return nil, err
	}
	if _, err := requireOneRow(rows); err != nil {
		return nil, err
	}
	e := groupCollectionRows(rows)[0]
	return &model.Cart{
		UserID:    model.NewID(e.userID),
		ProductID: model.NewID(e.productID),
		CartDate:  model.NewDate(e.date),
		Game:      e.game,
	}, nil
}

func (r *CartRepository) DeleteCart(userID, gameID int) (*model.Cart, error) {
	output, err := r.getCart(userID, gameID)
	if err != nil {
		return nil, err
	}
	if err := deleteCollectionEntry(userID, gameID, "carts"); err != nil {
		return nil, err
	}
	return output, nil
}
