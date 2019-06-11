package main

import (
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	connStr := "user=postgres dbname=cryptosoccer"
	err := storage.Init(connStr)
	if err != nil {
		log.Fatal(err)
	}
}