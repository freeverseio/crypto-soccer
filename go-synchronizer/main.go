package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	log "github.com/sirupsen/logrus"
)

func waitForInterrupt() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	log.Info("ctrl + c to interrupt")
	<-done
}

func main() {
	log.Info("Setup the synchronizer ...")

	process := process.BackgroundProcessNew()

	// log.Info("Start to process events ...")
	process.Start()

	waitForInterrupt()

	log.Info("Stop to process events ...")
	process.StopAndJoin()

	// connStr := "user=postgres dbname=cryptosoccer"
	// err := storage.Init(connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Info("Exiting ...")
}
