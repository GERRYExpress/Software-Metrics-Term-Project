package repository

import (
	"database/sql"
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

func (r *CategoryRepository) Create(name string) (model.Category, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert("categories").
		Columns("name", "is_hidden").
		Values(name, false).
		ToSql()
	if err != nil {
		return model.Category{}, err
	}

	result, err := config.DB.Exec(query, args...)
	if err != nil {
		return model.Category{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.Category{}, err
	}

	return model.NewCategory(int(id), name, false), nil
}

func (r *CategoryRepository) Update(id int, name string, isHidden bool) (model.Category, error) {
	result, err := config.DB.Exec(
		"UPDATE categories SET name = ?, is_hidden = ? WHERE category_id = ?",
		name, isHidden, id,
	)
	if err != nil {
		return model.Category{}, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return model.Category{}, err
	}
	if affected == 0 {
		return model.Category{}, sql.ErrNoRows
	}
	return model.NewCategory(id, name, isHidden), nil
}
