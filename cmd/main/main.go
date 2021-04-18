package main

import (
	"log"

	"github.com/ToirovSadi/wallet/pkg/wallet"
)

func main() {
	path := "data/file1.txt"

	s := wallet.Service{}

	err := s.ImportFromFile(path)
	if err != nil {
		log.Println(err)
		return
	}

	for i := 1; i <= int(s.NumAccount()); i++ {
		account, err := s.FindAccountByID(int64(i))
		if err != nil {
			log.Print(err)
			return
		}
		log.Printf("%#v\n", account)

	}

	log.Println("all done!")
}
