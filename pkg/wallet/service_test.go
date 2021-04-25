package wallet

import (
	"reflect"
	"testing"

	"github.com/ToirovSadi/wallet/pkg/types"
)

func BenchmarkService_SumPayments(b *testing.B) {
	s := Service{}
	for i := 0; i < 1000; i++ {
		s.payments = append(s.payments, &types.Payment{
			AccountID: int64(i),
			Amount:    5,
		})
	}
	op_per_go := 5
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := s.SumPayments(len(s.payments) / op_per_go)
		b.StopTimer()
		want := types.Money(0)
		for _, payment := range s.payments {
			want += payment.Amount
		}
		if got != want {
			b.Fatalf("error:\ngot: %v\nwant: %v\n", got, want)
		}
		b.StartTimer()
	}
}

func TestService_FindAccountByID_success(t *testing.T) {

	svc := &Service{}

	account, err := svc.RegisterAccount("+92349234")
	if err != nil {
		t.Error("error:TestService_FindAccountByID_success(): ", err)
	}

	_, err = svc.FindAccountByID(account.ID)

	if err != nil {
		t.Error("error:TestService_FindAccountByID_success(): ", err)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error("error:TestService_FindAccountByID_notFound(): ", err)
	}

	_, err = s.FindAccountByID(2)
	if err == nil {
		t.Error("error:TestService_FindAccountByID_notFound(): ", err)
	}

}

func TestService_FindPaymentByID(t *testing.T) {
	svc := newTestService()
	_, _, err := svc.addAccount(defaultTestAccount)
	if err != nil {
		t.Error("error:TestService_FindPaymentByID(): ", err)
	}
	// lets check when account.Balance < amount
	_, err = svc.Pay(1, 1000000000, "home")
	if err == nil {
		t.Error("error:TestService_FindPaymentByID(): function Pay not working")
	}

	payment, err := svc.Pay(1, 50, "internet")
	if err != nil {
		t.Error("error:TestService_FindPaymentByID(): ", err)
	}

	_, err = svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Error("error:TestService_FindPaymentByID(): ", err)
	}
	_, err = svc.FindPaymentByID("nothing")
	if err == nil {
		t.Error("error:TestService_FindPaymentByID(): find wrong payment :(")
	}
}

func TestService_Reject(t *testing.T) {
	svc := newTestService()
	_, _, err := svc.addAccount(defaultTestAccount)
	if err != nil {
		t.Error("error:TestService_FindPaymentByID(): ", err)
	}

	accountBeforePay, err := svc.FindAccountByID(1)
	if err != nil {
		t.Error("error:TestService_FindPaymentByID(): ", err)
	}
	payment, err := svc.Pay(1, 50, "internet")
	if err != nil {
		t.Error("error:TestService_FindPaymentByID(): ", err)
	}
	err = svc.Reject(payment.ID)
	if err != nil {
		t.Error("error:TestService_Reject(): ", err)
	}
	if payment.Status != types.PaymentStatusFail {
		t.Error("error:TestService_Reject(): payment status didn't changed")
	}
	account, _ := svc.FindAccountByID(1)
	if account.Balance != accountBeforePay.Balance {
		t.Error("error:TestService_Reject(): Reject not working")
	}
}

func TestService_Repeat(t *testing.T) {
	svc := newTestService()
	account, payments, err := svc.addAccount(defaultTestAccount)
	if err != nil {
		t.Error("error:TestService_Repeat(): ", err)
	}

	_, err = svc.Repeat(payments[0].ID)
	if err != nil {
		t.Error("error:TestService_Repeat(): ", err)
	}
	if account.Balance != (defaultTestAccount.balance - 2*payments[0].Amount) {
		t.Error("error:TestService_Repeat(): Repeat function not working")
	}
}

func TestService_FavoritePayment(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error("error:TestService_FavoritePayment(): ", err)
	}
	favorite, err := s.FavoritePayment(payments[0].ID, "just for checking")
	if err != nil {
		t.Error("error:TestService_FavoritePayment(): ", err)
	}
	tempFavorite, err := s.FindFavoriteByID(favorite.ID)
	if err != nil {
		t.Error("error:TestService_FavoritePayment(): ", err)
	}
	if !reflect.DeepEqual(favorite, tempFavorite) {
		t.Errorf("error:TestService_FavoritePayment():\ngot:%v\nwant:%v\n", tempFavorite, favorite)
	}
	_, err = s.FindFavoriteByID("nothing")
	if err == nil {
		t.Error("error:TestService_FavoritePayment(): find non-existent favorite payment")
	}
}

func TestService_PayFromFavorite(t *testing.T) {
	s := newTestService()
	account, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error("error:TestService_FavoritePayment(): ", err)
	}
	accountBeforePay := *account
	favorite, err := s.FavoritePayment(payments[0].ID, "just for checking")
	if err != nil {
		t.Error("error:TestService_FavoritePayment(): ", err)
	}
	payment, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Error("error:TestService_FavoritePayment(): ", err)
	}

	if (accountBeforePay.Balance - payment.Amount) != account.Balance {
		t.Error("error:TestService_FavoritePayment(): PayFromFavorite not working")
	}
}

// func regAccounts() *Service {
// 	svc := &Service{}

// 	svc.RegisterAccount("1")
// 	svc.RegisterAccount("2")
// 	svc.RegisterAccount("3")
// 	svc.RegisterAccount("4")
// 	svc.Deposit(1, 100)
// 	svc.Deposit(2, 100)
// 	// test for Service.Deposit
// 	account, _ := svc.FindAccountByID(1)
// 	if account.Balance != 100 {
// 		panic("error:regAccounts():Service.Deposit dosit not working(")
// 	}

// 	// test for Service.Pay
// 	_, err := svc.Pay(1, 50, "internet")
// 	if err != nil {
// 		panic("error:regAccounts():Service.Pay can't withdraw money")
// 	}

// 	/// test for Service.RegisterAccount
// 	_, err = svc.RegisterAccount("1")
// 	if err == nil {
// 		panic("error:regAccounts():Service.RegisterAccount can't match by phone!")
// 	}

// 	return svc
// }

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+19242352545",
	balance: 100_000_000,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 100, category: "cafe"},
	},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, err
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, err
	}
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, err
		}
	}
	return account, payments, nil
}
