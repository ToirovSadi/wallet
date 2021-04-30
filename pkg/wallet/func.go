package wallet

import (
	"github.com/ToirovSadi/wallet/pkg/types"
	"github.com/google/uuid"
	"os"
	"sync"
)

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	payment.Status = types.PaymentStatusFail
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}
	account.Balance += payment.Amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	newPayment, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}
	return newPayment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	favoriteID := uuid.New().String()
	favorite := &types.Favorite{
		ID:        favoriteID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}
	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return err
	}
	account.Balance += amount
	return nil
}
func (s *Service) SumPayments(goroutines int) types.Money {
	var sumPayments types.Money = 0
	if goroutines <= 1 {
		for _, payment := range s.payments {
			sumPayments += payment.Amount
		}
	} else {
		l := 0
		mu := sync.Mutex{}
		n := len(s.payments)
		wg := sync.WaitGroup{}
		wg.Add((n + goroutines - 1) / goroutines)
		for _, r := range DelN(n, goroutines) {
			go func(payments []*types.Payment) {
				defer wg.Done()
				var sum types.Money = 0
				for _, payment := range payments {
					sum += payment.Amount
				}
				mu.Lock()
				sumPayments += sum
				mu.Unlock()
			}(s.payments[l:r])
			l = r
		}
		wg.Wait()
	}
	return sumPayments
}

func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
	var resPayments []types.Payment
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if goroutines <= 1 {
		for _, payment := range s.payments {
			if payment.AccountID == accountID {
				resPayments = append(resPayments, *payment)
			}
		}
	} else {
		l := 0
		mu := sync.Mutex{}
		n := len(s.payments)
		wg := sync.WaitGroup{}
		wg.Add((n + goroutines - 1) / goroutines)
		for _, r := range DelN(n, goroutines) {
			go func(payments []*types.Payment) {
				defer wg.Done()
				var tempPayments []types.Payment
				for _, payment := range payments {
					if payment.AccountID == accountID {
						tempPayments = append(tempPayments, *payment)
					}
				}
				mu.Lock()
				resPayments = append(resPayments, tempPayments...)
				mu.Unlock()
			}(s.payments[l:r])
			l = r
		}
		wg.Wait()
	}
	return resPayments, nil
}

func (s *Service) FilterPaymentsByFn(
	filter func(payment types.Payment) bool,
	goroutines int) ([]types.Payment, error) {
	var resPayments []types.Payment
	if goroutines <= 1 {
		for _, payment := range s.payments {
			if filter(*payment) {
				resPayments = append(resPayments, *payment)
			}
		}
	} else {
		l := 0
		mu := sync.Mutex{}
		n := len(s.payments)
		wg := sync.WaitGroup{}
		wg.Add((n + goroutines - 1) / goroutines)
		for _, r := range DelN(n, goroutines) {
			go func(payments []*types.Payment) {
				defer wg.Done()
				var tempPayments []types.Payment
				for _, payment := range payments {
					if filter(*payment) {
						tempPayments = append(tempPayments, *payment)
					}
				}
				mu.Lock()
				resPayments = append(resPayments, tempPayments...)
				mu.Unlock()
			}(s.payments[l:r])
			l = r
		}
		wg.Wait()
	}
	if resPayments == nil {
		return nil, ErrPaymentNotFound
	}
	return resPayments, nil
}

func getString(data ...string) (res string) {
	for i := 0; i < len(data); i++ {
		res += data[i]
		if i == (len(data) - 1) {
			continue
		}
		res += ";"
	}
	return res + "|"
}

func (s *Service) NumAccount() int {
	return len(s.accounts)
}

func (s *Service) GetAccounts() []*types.Account {
	return s.accounts
}
func (s *Service) GetPayments() []*types.Payment {
	return s.payments
}
func (s *Service) GetFavorites() []*types.Favorite {
	return s.favorites
}
func (s *Service) ExistFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// getN - devide n into m equal parts
func DelN(n int, m int) []int {
	a := []int{}
	if m == 0 {
		panic("you can't devide intiger by zero")
	}
	for i := 1; i <= n/m; i++ {
		a = append(a, i*m)
	}
	if n%m != 0 {
		a = append(a, n)
	}
	return a
}
