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

func main() {
	configFile := flag.String("config", "./config.json", "configuration file")
	inMemoryDatabase := flag.Bool("memory", false, "use in memory database")
	debug := flag.Bool("debug", false, "print debug logs")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Starting ...")
	log.Info("Parsing configuration file: ", *configFile)
	config, err := config.New(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	config.Print()

	log.Info("Dial the Ethereum client: ", config.EthereumClient)
	conn, err := ethclient.Dial(config.EthereumClient)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Assets bindings to: ", config.AssetsContractAddress)
	assetsContract, err := assets.NewAssets(common.HexToAddress(config.AssetsContractAddress), conn)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	var sto *storage.Storage
	if *inMemoryDatabase {
		log.Warning("Using in memory DBMS (no persistence)")
		sto, err = storage.NewSqlite3("./sql/00_schema.sql")
	} else {
		log.Info("Connecting to DBMS: ", config.PostgresUrl)
		sto, err = storage.NewPostgres(config.PostgresUrl)
	}
	if err != nil {
		log.Fatalf("Failed to connect to DBMS: %v", err)
	}

	process := process.BackgroundProcessNew(assetsContract, sto)

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
