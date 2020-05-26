package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/staker"
)

func main() {
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable", "postgres url")
	namesDatabase := flag.String("namesDatabase", "./names.db", "name database path")
	debug := flag.Bool("debug", false, "print debug logs")
	ethereumClient := flag.String("ethereum", "http://localhost:8545", "ethereum node")
	directoryContractAddress := flag.String("directory_address", "", "")
	stakerPrivateKey := flag.String("staker", "", "the private key if it's a staker")
	ipfsURL := flag.String("ipfs", "localhost:5001", "ipfs node url")
	delta := flag.Int("delta", 10, "number of block to process at maximum")
	flag.Parse()

	if _, err := os.Stat(*namesDatabase); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("no names database file at %v", *namesDatabase)
		}
	}

	if *directoryContractAddress == "" {
		log.Fatal("no directory contract address")
	}

	log.Infof("ipfs URL: %v", *ipfsURL)

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Starting ...")

	if err := func() error {
		log.Info("Dial the Ethereum client: ", *ethereumClient)
		client, err := ethclient.Dial(*ethereumClient)
		if err != nil {
			return err
		}
		defer client.Close()

		var stkr *staker.Staker
		if *stakerPrivateKey != "" {
			log.Info("WARNING: STAKER address set")
			privateKey, err := crypto.HexToECDSA(*stakerPrivateKey)
			if err != nil {
				return err
			}
			stkr, err = staker.New(privateKey)
			if err != nil {
				return err
			}

		}

		log.Info("Connecting to universe DBMS: ", *postgresURL)
		universedb, err := storage.New(*postgresURL)
		if err != nil {
			return err
		}
		defer universedb.Close()

		namesdb, err := names.New(*namesDatabase)
		if err != nil {
			return err
		}
		defer namesdb.Close()

		processor := process.NewEventProcessor(
			client,
			*directoryContractAddress,
			namesdb,
			*ipfsURL,
			stkr,
		)

		log.Info("On Going ...")

		for {
			tx, err := universedb.Begin()
			if err != nil {
				return err
			}
			processedBlocks, err := processor.Process(tx, uint64(*delta))
			if err != nil {
				log.Error(err)
				tx.Rollback()
				continue
			}
			if err := tx.Commit(); err != nil {
				return err
			}

			if processedBlocks == 0 {
				time.Sleep(2 * time.Second)
			}
		}
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
