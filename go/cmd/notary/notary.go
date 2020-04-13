package main

import (
	"flag"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/processor"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	log "github.com/sirupsen/logrus"
)

func main() {
	inMemoryDatabase := flag.Bool("memory", false, "use in memory database")
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable", "postgres url")
	ethereumClient := flag.String("ethereum", "http://localhost:8545", "ethereum node")
	marketContractAddress := flag.String("market_address", "", "market contract address")
	constantsgettersContractAddress := flag.String("constantsgetters_address", "", "constantsgetters contract address")
	privateKeyHex := flag.String("private_key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "private key")
	debug := flag.Bool("debug", false, "print debug logs")
	flag.Parse()

	log.Infof("[PARAM] memory            : %v", *inMemoryDatabase)
	log.Infof("[PARAM] postgres          : %v", *postgresURL)
	log.Infof("[PARAM] ethereum_client   : %v", *ethereumClient)
	log.Infof("[PARAM] market_address    : %v", *marketContractAddress)
	log.Infof("[PARAM] constantsgetters_address    : %v", *constantsgettersContractAddress)
	privateKey, err := crypto.HexToECDSA(*privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("[PARAM] Address : %v", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	log.Infof("[PARAM] debug             : %v", *debug)
	log.Infof("-------------------------------------------------------------------")

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Create the connection to DBMS")
	var sto *storage.Storage
	if *inMemoryDatabase {
		log.Warning("Using in memory DBMS (no persistence)")
		sto, err = storage.NewSqlite3("../../../market.db/00_schema.sql")
	} else {
		log.Info("Connecting to DBMS: ", *postgresURL)
		sto, err = storage.NewPostgres(*postgresURL)
	}
	if err != nil {
		log.Fatalf("Failed to connect to DBMS: %v", err)
	}
	log.Info("Dial the Ethereum client: ", *ethereumClient)
	client, err := ethclient.Dial(*ethereumClient)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	contracts, err := contracts.New(
		client,
		"", "", "", "", "", "",
		*marketContractAddress,
		"", "", "", "",
		*constantsgettersContractAddress,
	)
	if err != nil {
		log.Fatal(err)
	}
	processor, err := processor.NewProcessor(sto, contracts, privateKey)
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

	ch := make(chan interface{}, 100000)

	go gql.NewServer(ch)
	// go submitactions.NewSubmitTimer(ch, 5*time.Second)

	consumer.NewConsumer(
		ch,
	).Start()
}
