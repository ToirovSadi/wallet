package wallet

import (
	"github.com/ToirovSadi/wallet/pkg/types"
	"reflect"
	"testing"
)

func TestDelN(t *testing.T) {
	type args struct {
		n int
		m int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			args: args{n: 10, m: 3},
			want: []int{3, 6, 9, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DelN(tt.args.n, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DelN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Deposit(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}
	type args struct {
		accountID int64
		amount    types.Money
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				accounts: []*types.Account{
					{ID: 1, Balance: 0},
				},
			},
			args: args{
				accountID: 1,
				amount:    10,
			},
			wantErr: false,
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
			if err := s.Deposit(tt.args.accountID, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("Deposit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ExistFile(t *testing.T) {
	s := Service{}
	ok, _ := s.ExistFile("nowhere/nothing")
	if ok == true {
		t.Fatal("wrong")
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

func TestService_GetAccounts(t *testing.T) {
	s := Service{}
	if !reflect.DeepEqual(s.GetAccounts(), s.accounts) {
		t.Fatal("wrong")
	}
}

func TestService_GetFavorites(t *testing.T) {
	s := Service{}
	if !reflect.DeepEqual(s.GetFavorites(), s.favorites) {
		t.Fatal("wrong")
	}
}

func TestService_GetPayments(t *testing.T) {
	s := Service{}
	if !reflect.DeepEqual(s.GetPayments(), s.payments) {
		t.Fatal("wrong")
	}
}

func TestService_NumAccount(t *testing.T) {
	s := Service{}
	if len(s.accounts) != s.NumAccount() {
		t.Fatal("wrong")
	}
}

func TestService_Pay(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}
	type args struct {
		accountID int64
		amount    types.Money
		category  types.PaymentCategory
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Payment
		wantErr bool
	}{
		{
			fields: fields{
				nextAccountID: 1,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				accountID: 1,
				amount:    5,
				category:  "",
			},
			wantErr: false,
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
			_, err := s.Pay(tt.args.accountID, tt.args.amount, tt.args.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
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

func Test_getString(t *testing.T) {
	type args struct {
		data []string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		{
			args: args{
				data: []string{"a", "b"},
			},
			wantRes: "a;b|",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := getString(tt.args.data...); gotRes != tt.wantRes {
				t.Errorf("getString() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

var accounts = []*types.Account{
	{ID: 1, Balance: 10, Phone: "+1"},
	{ID: 2, Balance: 10, Phone: "+2"},
	{ID: 3, Balance: 10, Phone: "+3"},
}

var payments = []*types.Payment{
	{
		ID:        "1",
		AccountID: 1,
		Amount:    10,
		Category:  "",
		Status:    "INPROGRESS",
	},
	{
		ID:        "2",
		AccountID: 2,
		Amount:    10,
		Category:  "",
		Status:    "INPROGRESS",
	},
	{
		ID:        "3",
		AccountID: 3,
		Amount:    10,
		Category:  "",
		Status:    "INPROGRESS",
	},
}

var favorites = []*types.Favorite{
	{
		ID:        "1",
		AccountID: 1,
		Name:      "",
		Amount:    10,
		Category:  "",
	},
	{
		ID:        "2",
		AccountID: 2,
		Name:      "",
		Amount:    10,
		Category:  "",
	},
	{
		ID:        "3",
		AccountID: 3,
		Name:      "",
		Amount:    10,
		Category:  "",
	},
}
