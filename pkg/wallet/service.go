package wallet

import (
	"errors"

	"github.com/ToirovSadi/wallet/pkg/types"
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	// payments      []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistred
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) FindAccountById(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

// Errors that can occur in these functions
var ErrAccountNotFound = errors.New("account that you want doesn't exist")
var ErrPhoneRegistred = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
