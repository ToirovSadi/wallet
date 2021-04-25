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

//func (s *Service) ExportAccountHistory(accountID int64) ([]*types.Payment, error) {
//	_, err := s.FindAccountByID(accountID)
//	if err != nil {
//		return nil, err
//	}
//	var payments []*types.Payment
//
//	for _, payment := range s.payments {
//		if payment.AccountID != accountID {
//			continue
//		}
//		payments = append(payments, payment)
//	}
//	return payments, nil
//}
//
//func (s *Service) HistoryToFiles(payments []*types.Payment, dir string, records int) error {
//	if len(payments) <= records {
//		err := s.exportPayments(dir + "/payments.dump")
//		return err
//	} else {
//		n := len(payments)
//		for i := 0; i < n/records; i++ {
//			l := i * records
//			r := (i + 1) * records // -> not included
//			tempP := s.payments
//			s.payments = s.payments[l:r]
//			err := s.exportPayments(dir + "/payments" + strconv.Itoa(i+1) + ".dump")
//			if err != nil {
//				return err
//			}
//			s.payments = tempP
//		}
//		rem := n % records
//		if rem != 0 {
//			tempP := s.payments
//			s.payments = s.payments[n-rem:]
//			err := s.exportPayments(dir + "/payments" + strconv.Itoa((n+records-1)/records) + ".dump")
//			if err != nil {
//				return err
//			}
//			s.payments = tempP
//		}
//	}
//	return nil
//}
