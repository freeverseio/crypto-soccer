package main

import (
	"flag"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

func main() {
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable", "postgres url")
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
	ipfsURL := flag.String("ipfs", "localhost:5001", "ipfs node url")
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
	if *enginePreCompContractAddress == "" {
		log.Fatal("no enginePreComp contract address")
	}

	log.Infof("ipfs URL: %v", *ipfsURL)

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

	log.Info("Connecting to universe DBMS: ", *postgresURL)
	universedb, err := storage.New(*postgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to universe DBMS: %v", err)
	}

	namesdb, err := names.New(*namesDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to names DBMS: %v", err)
	}

	log.Info("All is ready ... 5 seconds to start ...")
	time.Sleep(5 * time.Second)

	processor, err := process.NewEventProcessor(contracts, namesdb)
	if err != nil {
		log.Fatal(err)
	}

	delta := uint64(1000)
	for {
		tx, err := universedb.Begin()
		if err != nil {
			log.Fatal(err)
		}
		processedBlocks, err := processor.Process(tx, delta)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		tx.Commit()
		if processedBlocks == 0 {
			time.Sleep(2 * time.Second)
		}
	}
}
