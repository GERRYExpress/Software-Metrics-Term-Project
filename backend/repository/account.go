package repository

import (
	"database/sql"
	"errors"
	"software-metrics-project/config"
	"software-metrics-project/model"
	"time"

	"github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

type AccountRepository struct {
	Model model.Account
}

var ErrInvalidCredentials = errors.New("invalid credentials")

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

func (r *AccountRepository) Authenticate(username, password string) (model.Account, error) {
	var id int
	var hash string
	if err := config.DB.QueryRow(
		"SELECT user_id, password FROM accounts WHERE username = ?",
		username,
	).Scan(&id, &hash); err != nil {
		if err == sql.ErrNoRows {
			return model.Account{}, ErrInvalidCredentials
		}
		return model.Account{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		if hash != password {
			return model.Account{}, ErrInvalidCredentials
		}
		upgraded, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return model.Account{}, err
		}
		if _, err := config.DB.Exec("UPDATE accounts SET password = ? WHERE user_id = ?", upgraded, id); err != nil {
			return model.Account{}, err
		}
	}
	return model.NewAccount(id, username), nil
}

func (r *AccountRepository) Create(username, password string) (model.Account, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.Account{}, err
	}

	result, err := config.DB.Exec(
		"INSERT INTO accounts (username, password, created_date) VALUES (?, ?, ?)",
		username, hash, time.Now(),
	)
	if err != nil {
		return model.Account{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.Account{}, err
	}
	return model.NewAccount(int(id), username), nil
}
