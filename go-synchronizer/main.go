package main

import (
	_ "github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	_ "github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/fifo"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Start ...")

	bufferSize := 1000;

	log.Info("Creating the fifo with buffer size ", bufferSize)

	fifo := fifo.FifoNew(bufferSize)

	log.Info(fifo)
	
	// connStr := "user=postgres dbname=cryptosoccer"
	// err := storage.Init(connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Info("Exiting ...")
}