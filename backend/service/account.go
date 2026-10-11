package service

import (
	"software-metrics-project/auth"
	"software-metrics-project/model"
	"software-metrics-project/repository"
)

type AccountService struct {
	Repo *repository.AccountRepository
}

func (r *AccountService) GetAccounts() ([]model.Account, error) {
	return r.Repo.GetAccounts()
}

func (s *AccountService) Login(username, password string) (string, *model.Account, error) {
	account, err := s.Repo.Authenticate(username, password)
	if err != nil {
		return "", nil, err
	}
	token, err := auth.Issue(int(account.ID))
	if err != nil {
		return "", nil, err
	}
	return token, &account, nil
}

func (s *AccountService) Register(username, password string) (model.Account, error) {
	return s.Repo.Create(username, password)
}
