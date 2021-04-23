package main

import (
	"github.com/ToirovSadi/wallet/pkg/wallet"
	"log"
)

func main() {
	s := wallet.Service{}

	account, err := s.RegisterAccount("+9921231234")
	if err != nil{
		log.Print(err)
		return
	}

	err = s.Deposit(account.ID, 1000000)
	if err != nil{
		log.Print(err)
		return
	}

	for i := 0; i < 15; i ++{
		payment, err := s.Pay(account.ID, 5, "hichi")

		if err != nil{
			log.Print(err)
			return
		}
		if i % 3 == 0{
			_, err = s.FavoritePayment(payment.ID, "nothing")
			if err != nil{
				log.Print(err)
				return
			}
		}
	}

	err = s.Export("./data")
	if err != nil{
		log.Print(err)
		return
	}

	log.Print("Done!!!")
}