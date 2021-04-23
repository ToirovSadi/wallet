package wallet

import (
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ToirovSadi/wallet/pkg/types"
	"github.com/google/uuid"
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

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

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
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

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
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

func (s *Service)Export(dir string) (err error) {
	err = exportAccounts(s.accounts, dir + "\\accounts.dump")
	if err != nil{
		return err
	}
	err = exportPayments(s.payments, dir + "\\payments.dump")
	if err != nil{
		return err
	}
	err = exportFavorites(s.favorites, dir + "\\favorites.dump")
	return err
}

// Export accounts to indicated file
func exportAccounts(accounts []*types.Account, fileName string) (err error) {
	if len(accounts) == 0{
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil{
		return err
	}
	defer func(){
		cerr := file.Close()
		if err == nil{
			err = cerr
		}
	}()

	for _, account := range accounts{
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
	if len(payments) == 0{
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil{
		return err
	}
	defer func(){
		cerr := file.Close()
		if err == nil{
			err = cerr
		}
	}()

	for _, payment := range payments{
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
	if len(favorites) == 0{
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil{
		return err
	}
	defer func(){
		cerr := file.Close()
		if err == nil{
			err = cerr
		}
	}()

	for _, favorite := range favorites{
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

func getString(data ...string) (res string) {
	for i := 0; i < len(data); i ++{
		res += data[i]
		if i == (len(data) - 1){
			continue
		}
		res += ";"
	}
	return res + "|"
}

// Errors that can occur in these functions
var ErrAccountNotFound = errors.New("account that you want doesn't exist")
var ErrPhoneRegistred = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment that you asked not found")
var ErrFavoriteNotFound = errors.New("favorite payment that you ask not found")

func (s *Service) NumAccount() int {
	return len(s.accounts)
}
