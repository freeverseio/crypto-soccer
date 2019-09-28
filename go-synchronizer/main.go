package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/config"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/market"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
)

func main() {
	configFile := flag.String("config", "./config.json", "configuration file")
	inMemoryDatabase := flag.Bool("memory", false, "use in memory database")
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable", "postgres url")
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
	client, err := ethclient.Dial(config.EthereumClient)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Leagues bindings to: ", config.LeaguesContractAddress)
	leaguesContract, err := leagues.NewLeagues(common.HexToAddress(config.LeaguesContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Market bindings to: ", config.MarketContractAddress)
	marketContract, err := market.NewMarket(common.HexToAddress(config.MarketContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Updates bindings to: ", config.UpdatesContractAddress)
	updatesContract, err := updates.NewUpdates(common.HexToAddress(config.UpdatesContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	var sto *storage.Storage
	if *inMemoryDatabase {
		log.Warning("Using in memory DBMS (no persistence)")
		sto, err = storage.NewSqlite3("./sql/00_schema.sql")
	} else {
		log.Info("Connecting to DBMS: ", *postgresURL)
		sto, err = storage.NewPostgres(*postgresURL)
	}
	if err != nil {
		log.Fatalf("Failed to connect to DBMS: %v", err)
	}

	process := process.BackgroundProcessNew(client, sto, marketContract, leaguesContract, updatesContract)

	log.Info("Start processing events ...")
	process.Start()

	log.Info("Press 'ctrl + c' to interrupt")
	waitForInterrupt()

	log.Info("Stop processing events ...")
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
