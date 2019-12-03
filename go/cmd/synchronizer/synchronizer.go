package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func main() {
	inMemoryDatabase := flag.Bool("memory", false, "use in memory database")
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable", "postgres url")
	relayPostgresURL := flag.String("relayPostgres", "postgres://freeverse:freeverse@universe.db:5432/relay?sslmode=disable", "postgres url")
	namesDatabase := flag.String("namesDatabase", "./names.db", "name database path")
	debug := flag.Bool("debug", false, "print debug logs")
	ethereumClient := flag.String("ethereum", "http://localhost:8545", "ethereum node")
	leaguesContractAddress := flag.String("leaguesContractAddress", "", "")
	assetsContractAddress := flag.String("assetsContractAddress", "", "")
	evolutionContractAddress := flag.String("evolutionContractAddress", "", "")
	marketContractAddress := flag.String("marketContractAddress", "", "")
	updatesContractAddress := flag.String("updatesContractAddress", "", "")
	engineContractAddress := flag.String("engineContractAddress", "", "")
	enginePreCompContractAddress := flag.String("enginePreCompContractAddress", "", "")
	flag.Parse()

	if _, err := os.Stat(*namesDatabase); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("no names database file at %v", *namesDatabase)
		}
	}

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
	contracts, err := contracts.New(
		client,
		*leaguesContractAddress,
		*assetsContractAddress,
		*evolutionContractAddress,
		*engineContractAddress,
		*enginePreCompContractAddress,
		*updatesContractAddress,
		*marketContractAddress,
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var universedb *storage.Storage
	var relaydb *relay.Storage
	if *inMemoryDatabase {
		log.Warning("Using in memory DBMS (no persistence)")
		universedb, err = storage.NewSqlite3("./../../universe.db/00_schema.sql")
		relaydb, err = relay.NewSqlite3("./../../universe.db/00_schema.sql")
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

	namesdb, err := names.New(*namesDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to names DBMS: %v", err)
	}

	process, err := process.BackgroundProcessNew(
		contracts,
		universedb,
		relaydb,
		namesdb,
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
