package main

import (
	"fmt"
	"github.com/ToirovSadi/wallet/pkg/wallet"
	"log"
)

func main() {
	s := wallet.Service{}

	fmt.Println("accounts are:")
	for _, account := range s.GetAccounts() {
		fmt.Printf("%+v\n", account)
	}
	err := s.Import("./data")
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Println("accounts are:")
	for _, account := range s.GetAccounts() {
		fmt.Printf("%+v\n", *account)
	}

	//fmt.Println("payments are:")
	//for _, payment := range s.GetPayments(){
	//	fmt.Printf("%+v\n", payment)
	//}
	//
	//fmt.Println("favorites are:")
	//for _, favorite := range s.GetFavorites(){
	//	fmt.Printf("%+v\n", favorite)
	//}

	//account, err := s.RegisterAccount("+9921231234")
	//account, err = s.RegisterAccount("+9921231232")
	//account, err = s.RegisterAccount("+9921231233")
	//if err != nil{
	//	log.Print(err)
	//	return
	//}
	//
	//err = s.Deposit(account.ID, 1000000)
	//if err != nil{
	//	log.Print(err)
	//	return
	//}
	//
	//for i := 0; i < 15; i ++{
	//	payment, err := s.Pay(account.ID, 5, "hichi")
	//
	//	if err != nil{
	//		log.Print(err)
	//		return
	//	}
	//	if i % 3 == 0{
	//		_, err = s.FavoritePayment(payment.ID, "nothing")
	//		if err != nil{
	//			log.Print(err)
	//			return
	//		}
	//	}
	//}
	//
	//err = s.Export("./data")
	//if err != nil{
	//	log.Print(err)
	//	return
	//}

	log.Print("Done!!!")
}
