package main

import (
	_ "github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	_ "github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Setup the synchronizer ...")

	process := process.BackgroundProcessNew()

	log.Info("Start to process events ...")
	process.Start()	

	log.Info("Stop to process events ...")
	process.StopAndJoin()
	
	// connStr := "user=postgres dbname=cryptosoccer"
	// err := storage.Init(connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Info("Exiting ...")
}