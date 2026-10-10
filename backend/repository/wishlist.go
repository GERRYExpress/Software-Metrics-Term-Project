package repository

import (
	"database/sql"
	"software-metrics-project/config"
	"software-metrics-project/model"
	"time"

	"github.com/Masterminds/squirrel"
)

type WishlistRepository struct {
	Model model.Wishlist
}

func (r *WishlistRepository) GetWishlists(userId int) ([]model.Wishlist, error) {
	mysqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	query, args, err := mysqlBuilder.
		Select("w.user_id", "w.product_id", "w.wishlist_date", "p.name", "p.price", "p.poster_url", "c.category_id", "c.name").
		From("wishlists w").
		LeftJoin("products p on w.product_id = p.product_id").
		LeftJoin("product_categories pc ON p.product_id = pc.product_id").
		LeftJoin("categories c ON c.category_id = pc.category_id AND c.is_hidden = false").
		Where(squirrel.Eq{"w.user_id": userId}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wishlists := make([]model.Wishlist, 0)
	wishlistIndex := make(map[int]int)
	for rows.Next() {
		var (
			wishlistUserID int
			productID      int
			wishlistDate   time.Time
			gameName       *string
			gamePrice      *float64
			posterUrl      *string
			categoryID     *int
			categoryName   *string
		)

		if err := rows.Scan(&wishlistUserID, &productID, &wishlistDate, &gameName, &gamePrice,
			&posterUrl, &categoryID, &categoryName); err != nil {
			return nil, err
		}

		wishlistPosition, ok := wishlistIndex[productID]
		if !ok {
			wishlist := model.Wishlist{
				UserID:       model.NewID(wishlistUserID),
				ProductID:    model.NewID(productID),
				WishlistDate: model.NewDate(wishlistDate),
				Game:         model.Game{},
			}
			if gameName != nil {
				wishlist.Game = model.NewGame(productID, *gameName, "", *gamePrice, *posterUrl)
			}
			wishlists = append(wishlists, wishlist)
			wishlistPosition = len(wishlists) - 1
			wishlistIndex[productID] = wishlistPosition
		}

		if categoryID != nil {
			wishlists[wishlistPosition].Game.Categories = append(wishlists[wishlistPosition].Game.Categories,
				model.NewCategory(*categoryID, *categoryName, false))
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return wishlists, nil
}

func (r *WishlistRepository) AddWishlist(userID, gameID int) (*model.Wishlist, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert("wishlists").
		Columns("user_id", "product_id", "wishlist_date").
		Values(userID, gameID, time.Now()).
		Suffix("ON DUPLICATE KEY UPDATE wishlist_date = wishlist_date").
		ToSql()

	if err != nil {
		return nil, err
	}

	_, err = config.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return r.getWishlist(userID, gameID)
}

func (r *WishlistRepository) getWishlist(userID, gameID int) (*model.Wishlist, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select("w.user_id", "w.product_id", "w.wishlist_date", "p.name", "p.price", "p.poster_url", "c.category_id", "c.name").
		From("wishlists w").
		Join("products p ON w.product_id = p.product_id").
		LeftJoin("product_categories pc ON p.product_id = pc.product_id").
		LeftJoin("categories c ON c.category_id = pc.category_id AND c.is_hidden = false").
		Where(squirrel.Eq{"w.user_id": userID, "w.product_id": gameID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var output model.Wishlist
	initialized := false
	for rows.Next() {
		var (
			wishlistUserID int
			productID      int
			wishlistDate   time.Time
			gameName       string
			gamePrice      float64
			posterURL      string
			categoryID     *int
			categoryName   *string
		)

		if err := rows.Scan(&wishlistUserID, &productID, &wishlistDate, &gameName, &gamePrice,
			&posterURL, &categoryID, &categoryName); err != nil {
			return nil, err
		}

		if !initialized {
			output = model.Wishlist{
				UserID:       model.NewID(wishlistUserID),
				ProductID:    model.NewID(productID),
				WishlistDate: model.NewDate(wishlistDate),
				Game:         model.NewGame(gameID, gameName, "", gamePrice, posterURL),
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

func (r *WishlistRepository) DeleteWishlist(userID, gameID int) (*model.Wishlist, error) {
	output, err := r.getWishlist(userID, gameID)
	if err != nil {
		return nil, err
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Delete("wishlists").
		Where(squirrel.Eq{"user_id": userID, "product_id": gameID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	if _, err := config.DB.Exec(query, args...); err != nil {
		return nil, err
	}

	return output, nil
}
