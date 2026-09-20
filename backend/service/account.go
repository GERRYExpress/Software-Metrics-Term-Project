package service

import (
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type AccountService struct {
	Repo *repository.AccountRepository
}

func (r *AccountService) GetAccounts() ([]model.Account, error) {
	return r.Repo.GetAccounts()
}
