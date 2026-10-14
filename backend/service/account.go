package service

import (
	"software-metrics-project/auth"
	"software-metrics-project/model"
)

type AccountService struct {
	Repo AccountRepository
}

type AccountRepository interface {
	GetAccounts() ([]model.Account, error)
	Authenticate(string, string) (model.Account, error)
	Create(string, string) (model.Account, error)
	Delete(int) error
}

func (s *AccountService) GetAccounts() ([]model.Account, error) {
	return s.Repo.GetAccounts()
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
	return s.Create(username, password)
}

func (s *AccountService) Create(username, password string) (model.Account, error) {
	return s.Repo.Create(username, password)
}

func (s *AccountService) Delete(userID int) error {
	return s.Repo.Delete(userID)
}
