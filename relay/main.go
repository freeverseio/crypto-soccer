package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/relay/contracts/updates"
	"github.com/freeverseio/crypto-soccer/relay/process"
	"github.com/freeverseio/crypto-soccer/relay/storage"
)

func main() {
	inMemoryDatabase := flag.Bool("memory", false, "use in memory database")
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable", "postgres url")
	debug := flag.Bool("debug", false, "print debug logs")
	ethereumClient := flag.String("ethereum", "http://localhost:8545", "ethereum node")
	privateKeyHex := flag.String("private_key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "private key")
	updatesContractAddress := flag.String("updatesContractAddress", "", "")
	flag.Parse()

	if *updatesContractAddress == "" {
		log.Fatal("no updates contract address")
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

	log.Info("Creating Updates bindings to: ", *updatesContractAddress)
	updatesContract, err := updates.NewUpdates(common.HexToAddress(*updatesContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	var sto *storage.Storage
	if *inMemoryDatabase {
		log.Warning("Using in memory DBMS (no persistence)")
		sto, err = storage.NewSqlite3("./db/00_schema.sql")
	} else {
		log.Info("Connecting to DBMS: ", *postgresURL)
		sto, err = storage.NewPostgres(*postgresURL)
	}
	if err != nil {
		log.Fatalf("Failed to connect to DBMS: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(*privateKeyHex)
	if err != nil {
		log.Fatal("Unable to obtain privateLey")
	}
	process, err := relay.BackgroundProcessNew(client, privateKey, sto, updatesContract)
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
