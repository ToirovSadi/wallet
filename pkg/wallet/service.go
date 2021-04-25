package wallet

import (
	"github.com/ToirovSadi/wallet/pkg/types"
	"sync"
)

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
