package main

import (
	"flag"
	"time"

	"github.com/freeverseio/crypto-soccer/market/notary/processor"
	"github.com/freeverseio/crypto-soccer/market/notary/storage"

	log "github.com/sirupsen/logrus"
)

func main() {
	inMemoryDatabase := flag.Bool("memory", false, "use in memory database")
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable", "postgres url")
	ethereumClient := flag.String("ethereum_client", "https://devnet.busyverse.com/web3", "ethereum node")
	assetsContractAddress := flag.String("assets_address", "0x893D66E9eBd8Db7cb4d565f4e135fc87D2547696", "assets contract address")
	debug := flag.Bool("debug", false, "print debug logs")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	var err error
	var sto *storage.Storage
	if *inMemoryDatabase {
		log.Warning("Using in memory DBMS (no persistence)")
		sto, err = storage.NewSqlite3("../db/00_schema.sql")
	} else {
		log.Info("Connecting to DBMS: ", *postgresURL)
		sto, err = storage.NewPostgres(*postgresURL)
	}
	if err != nil {
		log.Fatalf("Failed to connect to DBMS: %v", err)
	}

	processor, err := processor.NewProcessor(sto, *ethereumClient, *assetsContractAddress)
	if err != nil {
		log.Fatal(err)
	}

	for {
		err = processor.Process()
		if err != nil {
			log.Error(err)
		}
		time.Sleep(2 * time.Second)
	}
}
