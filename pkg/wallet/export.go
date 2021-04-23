package wallet

import (
	"github.com/ToirovSadi/wallet/pkg/types"
	"os"
	"strconv"
)

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
func exportAccounts(accounts []*types.Account, fileName string) (err error) {
	if len(accounts) == 0 {
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

	for _, account := range accounts {
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
func exportPayments(payments []*types.Payment, fileName string) (err error) {
	if len(payments) == 0 {
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

	for _, payment := range payments {
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
func exportFavorites(favorites []*types.Favorite, fileName string) (err error) {
	if len(favorites) == 0 {
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

	for _, favorite := range favorites {
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
