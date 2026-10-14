package service

import (
	"errors"
	"os"
	"software-metrics-project/model"
	"testing"
)

// Mock implementation of AccountRepository
type mockAccountRepo struct {
	mockAccount model.Account
	mockErr     error
}

func (m *mockAccountRepo) GetAccounts() ([]model.Account, error) {
	if m.mockErr != nil {
		return nil, m.mockErr
	}
	return []model.Account{m.mockAccount}, nil
}

func (m *mockAccountRepo) Authenticate(username, password string) (model.Account, error) {
	return m.mockAccount, m.mockErr
}

func (m *mockAccountRepo) Create(username, password string) (model.Account, error) {
	return m.mockAccount, m.mockErr
}

func (m *mockAccountRepo) Delete(userID int) error {
	return m.mockErr
}

func TestAccountService_Login(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	t.Run("success", func(t *testing.T) {
		repo := &mockAccountRepo{mockAccount: model.Account{ID: 1}}
		svc := &AccountService{Repo: repo}

		token, acc, err := svc.Login("user", "pass")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == "" {
			t.Errorf("expected token, got empty string")
		}
		if acc.ID != 1 {
			t.Errorf("expected account ID 1, got %v", acc.ID)
		}
	})

	t.Run("auth failure", func(t *testing.T) {
		repo := &mockAccountRepo{mockErr: errors.New("invalid credentials")}
		svc := &AccountService{Repo: repo}

		_, _, err := svc.Login("user", "wrongpass")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestAccountService_GetAccounts(t *testing.T) {
	repo := &mockAccountRepo{mockAccount: model.Account{ID: 1}}
	svc := &AccountService{Repo: repo}

	accounts, err := svc.GetAccounts()
	if err != nil || len(accounts) == 0 {
		t.Errorf("failed to get accounts")
	}
}

func TestAccountService_Create(t *testing.T) {
	repo := &mockAccountRepo{mockAccount: model.Account{ID: 2}}
	svc := &AccountService{Repo: repo}

	acc, err := svc.Create("newuser", "pass")
	if err != nil || acc.ID != 2 {
		t.Errorf("failed to create account")
	}
}

func TestAccountService_Delete(t *testing.T) {
	repo := &mockAccountRepo{}
	svc := &AccountService{Repo: repo}

	if err := svc.Delete(1); err != nil {
		t.Errorf("failed to delete account")
	}
}
