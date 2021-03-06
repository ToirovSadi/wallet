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

func TestService_SumPayments(t *testing.T) {
	s := Service{}
	res := types.Money(0)
	for i := 0; i < 100; i++ {
		s.payments = append(s.payments, &types.Payment{
			AccountID: 1,
			Amount:    10,
		})
		res += 10
	}
	sum := s.SumPayments(10)
	if res != sum {
		t.Fatalf("want: %v\n, got: %v\n", res, sum)
	}

	sum = s.SumPayments(1)
	if sum != res {
		t.Fatalf("want: %v\n, got: %v\n", res, sum)
	}
}

func TestService_FilterPayments(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}
	type args struct {
		accountID  int64
		goroutines int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []types.Payment
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				accountID:  1,
				goroutines: 2,
			},
			want:    []types.Payment{*payments[0]},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				accountID:  10,
				goroutines: 2,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				accountID:  10,
				goroutines: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
				favorites:     tt.fields.favorites,
			}
			got, err := s.FilterPayments(tt.args.accountID, tt.args.goroutines)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterPayments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterPayments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func filter(payment types.Payment) bool {
	return payment.AccountID == 1
}

func TestService_FilterPaymentsByFn(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}
	type args struct {
		filter     func(payment types.Payment) bool
		goroutines int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []types.Payment
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				filter:     filter,
				goroutines: 2,
			},
			want:    []types.Payment{*payments[0]},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				filter: func(payment types.Payment) bool {
					return payment.AccountID == 10
				},
				goroutines: 2,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				filter: func(payment types.Payment) bool {
					return payment.AccountID == 10
				},
				goroutines: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
				favorites:     tt.fields.favorites,
			}
			got, err := s.FilterPaymentsByFn(tt.args.filter, tt.args.goroutines)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterPaymentsByFn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterPaymentsByFn() got = %v, want %v", got, tt.want)
			}
		})
	}
}
