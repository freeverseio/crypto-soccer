package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/config"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

var ethereumClient = "https://devnet.busyverse.com/web3"
var assetsContractAddress = "0x05Fdd4d2340bcA823802849c75F385561278c3aB"
var postgresUrl = "user=postgres dbname=cryptosoccer"

func main() {
	configFile := flag.String("config", "./config.json", "configuration file")
	flag.Parse()

	log.Info("Starting ...")
	log.Info("Parsing configuration file: ", *configFile)
	config, err := config.New(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	config.Print()

	log.Info("Dial the Ethereum client: ", config.EthereumClient)
	conn, err := ethclient.Dial(ethereumClient)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Assets bindings to: ", config.AssetsContractAddress)
	assetsContract, err := assets.NewAssets(common.HexToAddress(config.AssetsContractAddress), conn)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Connecting to DBMS: ", config.PostgresUrl)
	storage, err := storage.NewPostgres(config.PostgresUrl)
	if err != nil {
		log.Fatalf("Failed to connect to DBMS: %v", err)
	}

	process := process.BackgroundProcessNew(assetsContract, storage)

	log.Info("Start to process events ...")
	process.Start()

	log.Info("Press 'ctrl + c' to interrupt")
	waitForInterrupt()

	log.Info("Stop to process events ...")
	process.StopAndJoin()

	log.Info("... exiting")
}

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
