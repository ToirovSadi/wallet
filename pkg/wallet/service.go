package wallet

import (
	"github.com/ToirovSadi/wallet/pkg/types"
	"strconv"
)

func (s *Service) ExportAccountHistory(accountID int64) ([]*types.Payment, error) {
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	var payments []*types.Payment

	for _, payment := range s.payments {
		if payment.AccountID != accountID {
			continue
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (s *Service) HistoryToFiles(payments []*types.Payment, dir string, records int) error {
	if len(payments) <= records {
		err := s.exportPayments(dir + "/payments.dump")
		return err
	} else {
		n := len(payments)
		for i := 0; i < n/records; i++ {
			l := i * records
			r := (i + 1) * records // -> not included
			tempP := s.payments
			s.payments = s.payments[l:r]
			err := s.exportPayments(dir + "/payments" + strconv.Itoa(i+1) + ".dump")
			if err != nil {
				return err
			}
			s.payments = tempP
		}
		rem := n % records
		if rem != 0 {
			tempP := s.payments
			s.payments = s.payments[n-rem:]
			err := s.exportPayments(dir + "/payments" + strconv.Itoa((n+records-1)/records) + ".dump")
			if err != nil {
				return err
			}
			s.payments = tempP
		}
	}
	return nil
}
