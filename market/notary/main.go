package main

import (
	"flag"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/market/notary/contracts/market"
	"github.com/freeverseio/crypto-soccer/market/notary/processor"
	"github.com/freeverseio/crypto-soccer/market/notary/storage"

	log "github.com/sirupsen/logrus"
)

func main() {
	inMemoryDatabase := flag.Bool("memory", false, "use in memory database")
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable", "postgres url")
	ethereumClient := flag.String("ethereum_client", "https://devnet.busyverse.com/web3", "ethereum node")
	marketContractAddress := flag.String("market_address", "", "market contract address")
	privateKeyHex := flag.String("private_key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "private key")
	debug := flag.Bool("debug", false, "print debug logs")
	flag.Parse()

	log.Infof("[PARAM] memory            : %v", *inMemoryDatabase)
	log.Infof("[PARAM] postgres          : %v", *postgresURL)
	log.Infof("[PARAM] ethereum_client   : %v", *ethereumClient)
	log.Infof("[PARAM] market_address    : %v", *marketContractAddress)
	log.Infof("[PARAM] debug             : %v", *debug)
	log.Infof("[PARAM] privatekey        : %v", *privateKeyHex)
	log.Infof("-------------------------------------------------------------------")

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Create the connection to DBMS")
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

	log.Infof("Dial the Ethereum client: %v", *ethereumClient)
	client, err := ethclient.Dial(*ethereumClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Creating Market bindings to: %v", *marketContractAddress)
	marketContract, err := market.NewMarket(common.HexToAddress(*marketContractAddress), client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(*privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	processor, err := processor.NewProcessor(sto, client, marketContract, privateKey)
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
