package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go/contracts/evolution"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/contracts/updates"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func main() {
	inMemoryDatabase := flag.Bool("memory", false, "use in memory database")
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable", "postgres url")
	relayPostgresURL := flag.String("relayPostgres", "postgres://freeverse:freeverse@relay.db:5432/relay?sslmode=disable", "postgres url")
	debug := flag.Bool("debug", false, "print debug logs")
	ethereumClient := flag.String("ethereum", "http://localhost:8545", "ethereum node")
	leaguesContractAddress := flag.String("leaguesContractAddress", "", "")
	assetsContractAddress := flag.String("assetsContractAddress", "", "")
	evolutionContractAddress := flag.String("evolutionContractAddress", "", "")
	marketContractAddress := flag.String("marketContractAddress", "", "")
	updatesContractAddress := flag.String("updatesContractAddress", "", "")
	engineContractAddress := flag.String("engineContractAddress", "", "")
	flag.Parse()

	if *leaguesContractAddress == "" {
		log.Fatal("no league contract address")
	}
	if *assetsContractAddress == "" {
		log.Fatal("no assets contract address")
	}
	if *evolutionContractAddress == "" {
		log.Fatal("no evolution contract address")
	}
	if *marketContractAddress == "" {
		log.Fatal("no market contract address")
	}
	if *updatesContractAddress == "" {
		log.Fatal("no updates contract address")
	}
	if *engineContractAddress == "" {
		log.Fatal("no engine contract address")
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Starting ...")
	log.Info("Dial the Ethereum client: ", *ethereumClient)
	client, err := ethclient.Dial(*ethereumClient)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Leagues bindings to: ", *leaguesContractAddress)
	leaguesContract, err := leagues.NewLeagues(common.HexToAddress(*leaguesContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Assets bindings to: ", *assetsContractAddress)
	assetsContract, err := assets.NewAssets(common.HexToAddress(*assetsContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Evolution bindings to: ", *evolutionContractAddress)
	evolutionContract, err := evolution.NewEvolution(common.HexToAddress(*evolutionContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Engine bindings to: ", *engineContractAddress)
	engineContract, err := engine.NewEngine(common.HexToAddress(*engineContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Updates bindings to: ", *updatesContractAddress)
	updatesContract, err := updates.NewUpdates(common.HexToAddress(*updatesContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	log.Info("Creating Market bindings to: ", *marketContractAddress)
	marketContract, err := market.NewMarket(common.HexToAddress(*marketContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	var universedb *storage.Storage
	var relaydb *relay.Storage
	if *inMemoryDatabase {
		log.Warning("Using in memory DBMS (no persistence)")
		universedb, err = storage.NewSqlite3("./../../universe.db/00_schema.sql")
		relaydb, err = relay.NewSqlite3("./../../relay.db/00_schema.sql")
	} else {
		log.Info("Connecting to universe DBMS: ", *postgresURL, " and relay DBMS: ", *relayPostgresURL)
		universedb, err = storage.NewPostgres(*postgresURL)
		if err != nil {
			log.Fatalf("Failed to connect to universe DBMS: %v", err)
		}
		relaydb, err = relay.NewPostgres(*relayPostgresURL)
		if err != nil {
			log.Fatalf("Failed to connect to relay DBMS: %v", err)
		}
	}

	process, err := process.BackgroundProcessNew(
		client,
		universedb,
		relaydb,
		engineContract,
		assetsContract,
		leaguesContract,
		updatesContract,
		marketContract,
		evolutionContract,
	)
	if err != nil {
		log.Fatal(err)
	}

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
