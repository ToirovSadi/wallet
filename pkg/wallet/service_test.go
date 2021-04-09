package wallet

import "testing"

func TestService_FindAccountByID_success(t *testing.T) {

	svc := &Service{}

	account, err := svc.RegisterAccount("+92349234")
	if err != nil {
		t.Error("error:TestService_FindAccountByID_success(): ", err)
	}

	_, err = svc.FindAccountById(account.ID)

	if err != nil {
		t.Error("error:TestService_FindAccountByID_success(): ", err)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := regAccounts()

	_, err := svc.FindAccountById(5)
	if err == nil {
		t.Error("error:TestService_FindAccountByID_notFound(): ", err)
	}

}

func regAccounts() *Service {
	svc := &Service{}

	svc.RegisterAccount("1")
	svc.RegisterAccount("2")
	svc.RegisterAccount("3")
	svc.RegisterAccount("4")
	
	/// lets test Service.RegisterAccount
	_, err := svc.RegisterAccount("1")
	if err == nil{
		panic("error:regAccounts():Service.RegisterAccount can't match by phone!")
	}

	return svc
}
