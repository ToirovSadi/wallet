package main

import (
	"log"

	"github.com/ToirovSadi/wallet/pkg/wallet"
)

func main() {
	path := "data/file1.txt"

	s := wallet.Service{}

	s.RegisterAccount("+1")
	s.RegisterAccount("+2")
	s.RegisterAccount("+3")

	err := s.ExportToFile(path)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("all done!")
}
