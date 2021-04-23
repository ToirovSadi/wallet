package wallet

import (
	"bufio"
	"errors"
	"github.com/ToirovSadi/wallet/pkg/types"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func (s *Service) ImportFromFile(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err1 := file.Close()
		if err1 != nil {
			err = err1
		}
	}()

	buf := make([]byte, 4096)
	data := make([]byte, 0)
	for {
		nread, err := file.Read(buf)
		if err == io.EOF {
			data = append(data, buf[:nread]...)
			break
		}
		if err != nil {
			return err
		}
		data = append(data, buf[:nread]...)
	}

	accountsDel := strings.Split(string(data), "|")

	var accounts []types.Account
	for _, account := range accountsDel {
		if len(account) == 0 { // empty
			continue
		}
		parts := strings.Split(account, ";")
		if len(parts) != 3 {
			return errors.New("error:ImportFromFile(): account contains of three parts(id, phone, balance)")
		}
		id, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return err
		}
		phone := parts[1]
		balance, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return err
		}
		accounts = append(accounts, types.Account{
			ID:      id,
			Phone:   types.Phone(phone),
			Balance: types.Money(balance),
		})
	}
	sort.Sort(ByID(accounts))

	for _, account := range accounts {
		s.accounts = append(s.accounts, &types.Account{
			ID:      account.ID,
			Phone:   account.Phone,
			Balance: account.Balance,
		})
	}

	return nil
}

func (s *Service) importAccounts(fileName string) error {
	ok, err := s.ExistFile(fileName)
	if ok == false {
		return err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	data, err := readByLine(file)
	if err != nil {
		return err
	}

	accountsDel := strings.Split(string(data), "|")

	var accountsNotSorted []types.Account
	for _, account := range accountsDel {
		if len(account) == 0 { // empty
			continue
		}
		parts := strings.Split(account, ";")
		if len(parts) != 3 {
			return errors.New("error:ImportFromFile(): account contains of three parts(id, phone, balance)")
		}
		id, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return err
		}
		phone := types.Phone(parts[1])
		balance, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return err
		}

		accountsNotSorted = append(accountsNotSorted, types.Account{
			ID:      id,
			Phone:   phone,
			Balance: types.Money(balance),
		})
	}

	sort.Sort(ByID(accountsNotSorted))

	for _, account := range accountsNotSorted {
		oldAccount, err := s.FindAccountByID(account.ID)
		if err == nil {
			account = *oldAccount
		}
		if s.nextAccountID < account.ID {
			s.nextAccountID = account.ID
		}
		s.accounts = append(s.accounts, &types.Account{
			ID:      account.ID,
			Phone:   account.Phone,
			Balance: account.Balance,
		})
	}
	return nil
}

func (s *Service) importPayments(fileName string) error {
	ok, err := s.ExistFile(fileName)
	if ok == false {
		return err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	data, err := readByLine(file)
	if err != nil {
		return err
	}
	paymentsDel := strings.Split(string(data), "|")
	for _, payment := range paymentsDel {
		if len(payment) == 0 { // empty
			continue
		}
		parts := strings.Split(payment, ";")
		if len(parts) != 5 {
			return errors.New("error:ImportFromFile(): payment contains of three parts(id, phone, balance)")
		}
		id := parts[0]                                       // 0 -> ID
		accountID, err := strconv.ParseInt(parts[1], 10, 64) // 1 -> accountID
		if err != nil {
			return err
		}
		amount, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return err
		}
		category := types.PaymentCategory(parts[3])
		status := types.PaymentStatus(parts[4])
		s.payments = append(s.payments, &types.Payment{
			ID:        id,
			AccountID: accountID,
			Amount:    types.Money(amount),
			Category:  category,
			Status:    status,
		})
	}
	return nil
}

func (s *Service) importFavorites(fileName string) error {
	ok, err := s.ExistFile(fileName)
	if ok == false {
		return err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	data, err := readByLine(file)
	if err != nil {
		return err
	}
	favoriteDel := strings.Split(string(data), "|")
	for _, favorite := range favoriteDel {
		if len(favorite) == 0 { // empty
			continue
		}
		parts := strings.Split(favorite, ";")
		if len(parts) != 5 {
			return errors.New("error:ImportFromFile(): payment contains of three parts(id, phone, balance)")
		}
		id := parts[0]                                       // 0 -> ID
		accountID, err := strconv.ParseInt(parts[1], 10, 64) // 1 -> accountID
		if err != nil {
			return err
		}
		name := parts[2]
		amount, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return err
		}
		category := types.PaymentCategory(parts[4])
		s.favorites = append(s.favorites, &types.Favorite{
			ID:        id,
			AccountID: accountID,
			Name:      name,
			Amount:    types.Money(amount),
			Category:  category,
		})
	}
	return nil
}

func readByLine(src *os.File) ([]byte, error) {
	data := make([]byte, 0)
	reader := bufio.NewReader(src)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			data = append(data, []byte(line)...)
			break
		}
		if err != nil {
			return nil, err
		}
		data = append(data, []byte(line)...)
	}
	return data, nil
}

type ByID []types.Account

func (a ByID) Len() int {
	return len(a)
}
func (a ByID) Swap(i int, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByID) Less(i int, j int) bool {
	return a[i].ID < a[j].ID
}
