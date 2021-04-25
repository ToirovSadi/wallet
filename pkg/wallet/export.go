package wallet

import (
	"github.com/ToirovSadi/wallet/pkg/types"
	"os"
	"strconv"
)

func (s *Service) Export(dir string) (err error) {
	err = s.exportAccounts(dir + "/accounts.dump")
	if err != nil {
		return err
	}
	err = s.exportPayments(dir + "/payments.dump")
	if err != nil {
		return err
	}
	err = s.exportFavorites(dir + "/favorites.dump")
	return err
}

func (s *Service) ExportToFile(path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err1 := file.Close()
		if err1 != nil {
			err = err1
		}
	}()
	for _, account := range s.accounts {
		id := strconv.FormatInt(account.ID, 10)
		phone := string(account.Phone)
		balance := strconv.FormatInt(int64(account.Balance), 10)
		_, err := file.Write([]byte(id + ";" + phone + ";" + balance + "|"))
		if err != nil {
			return err
		}
	}
	return nil
}

// Export accounts to indicated file
func (s *Service) exportAccounts(fileName string) (err error) {
	if len(s.accounts) == 0 {
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	for _, account := range s.accounts {
		id := strconv.FormatInt(account.ID, 10)
		phone := string(account.Phone)
		balance := strconv.FormatInt(int64(account.Balance), 10)
		_, err := file.Write([]byte(getString(id, phone, balance)))
		if err != nil {
			return err
		}
	}
	return nil
}

// Export payments to indicated file
func (s *Service) exportPayments(fileName string) (err error) {
	if len(s.payments) == 0 {
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	for _, payment := range s.payments {
		id := payment.ID
		accountID := strconv.FormatInt(payment.AccountID, 10)
		amount := strconv.FormatInt(int64(payment.Amount), 10)
		category := string(payment.Category)
		status := string(payment.Status)
		_, err := file.Write([]byte(getString(id, accountID, amount, category, status)))
		if err != nil {
			return err
		}
	}
	return nil
}

// Export payments to indicated file
func (s *Service) exportFavorites(fileName string) (err error) {
	if len(s.favorites) == 0 {
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	for _, favorite := range s.favorites {
		id := favorite.ID
		accountID := strconv.FormatInt(favorite.AccountID, 10)
		name := favorite.Name
		amount := strconv.FormatInt(int64(favorite.Amount), 10)
		category := string(favorite.Category)
		_, err := file.Write([]byte(getString(id, accountID, name, amount, category)))
		if err != nil {
			return err
		}
	}
	return nil
}

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
