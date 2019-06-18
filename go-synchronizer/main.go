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

	<-done
}

func main() {
	log.Info("Setup the synchronizer ...")

	process := process.BackgroundProcessNew()

	log.Info("Start to process events ...")
	process.Start()

	log.Info("Press 'ctrl + c' to interrupt")
	waitForInterrupt()

	log.Info("Stop to process events ...")
	process.StopAndJoin()

	log.Info("... exiting")
}
