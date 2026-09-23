package repository

import (
	"software-metrics-project/config"
	"software-metrics-project/model"

	"github.com/Masterminds/squirrel"
)

type AccountRepository struct {
	Model model.Account
}

func (r *AccountRepository) GetAccounts() ([]model.Account, error) {
	query, args, err := squirrel.
		Select("*").
		From("accounts").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account
	for rows.Next() {
		var account model.Account
		if err := rows.Scan(&account.ID, &account.Username, &account.Password,
			&account.CreatedDate, &account.IsBuyer, &account.IsSeller,
			&account.TotalGame, &account.TotalSpend, &account.TotalSale); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}
