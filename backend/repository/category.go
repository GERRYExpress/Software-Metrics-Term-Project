package repository

import (
	"software-metrics-project/config"
	"software-metrics-project/model"

	"github.com/Masterminds/squirrel"
)

type CategoryRepository struct {
	Model model.Category
}

func (r *CategoryRepository) GetCategories() ([]model.Category, error) {
	query, args, err := squirrel.
		Select("*").
		From("categories").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.IsHidden); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
