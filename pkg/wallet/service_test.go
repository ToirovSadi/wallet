package wallet

import (
	"testing"

	"github.com/ToirovSadi/wallet/pkg/types"
)

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
	svc := regAccounts()

	_, err := svc.FindAccountByID(5)
	if err == nil {
		t.Error("error:TestService_FindAccountByID_notFound(): ", err)
	}

}

func TestService_FindPaymentByID(t *testing.T) {
	svc := regAccounts()

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
	svc := regAccounts()

	// initial balance of account was 100
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
	if account.Balance != 50 {
		t.Error("error:TestService_Reject(): Reject not working")
	}
}

func TestService_Repeat(t *testing.T) {
	svc := regAccounts()

	payment, _ := svc.Pay(2, 50, "nothing")
	newPayment, err := svc.Repeat(payment.ID)
	if err != nil {
		t.Error("error:TestService_Repeat(): ", err)
	}
	account, err := svc.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Error("error:TestService_Repeat(): ", err)
	}
	if account.Balance != (100 - payment.Amount - newPayment.Amount) {
		t.Error("error:TestService_Repeat(): Repeat function not working")
	}
}

func regAccounts() *Service {
	svc := &Service{}

	svc.RegisterAccount("1")
	svc.RegisterAccount("2")
	svc.RegisterAccount("3")
	svc.RegisterAccount("4")
	svc.Deposit(1, 100)
	svc.Deposit(2, 100)
	// test for Service.Deposit
	account, _ := svc.FindAccountByID(1)
	if account.Balance != 100 {
		panic("error:regAccounts():Service.Deposit dosit not working(")
	}

	// test for Service.Pay
	_, err := svc.Pay(1, 50, "internet")
	if err != nil {
		panic("error:regAccounts():Service.Pay can't withdraw money")
	}

	/// test for Service.RegisterAccount
	_, err = svc.RegisterAccount("1")
	if err == nil {
		panic("error:regAccounts():Service.RegisterAccount can't match by phone!")
	}

	return svc
}
