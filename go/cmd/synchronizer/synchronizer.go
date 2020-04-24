package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
)

func main() {
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable", "postgres url")
	namesDatabase := flag.String("namesDatabase", "./names.db", "name database path")
	debug := flag.Bool("debug", false, "print debug logs")
	ethereumClient := flag.String("ethereum", "http://localhost:8545", "ethereum node")
	leaguesContractAddress := flag.String("leaguesContractAddress", "", "")
	assetsContractAddress := flag.String("assetsContractAddress", "", "")
	evolutionContractAddress := flag.String("evolutionContractAddress", "", "")
	engineContractAddress := flag.String("engineContractAddress", "", "")
	enginePreCompContractAddress := flag.String("enginePreCompContractAddress", "", "")
	updatesContractAddress := flag.String("updatesContractAddress", "", "")
	marketContractAddress := flag.String("marketContractAddress", "", "")
	utilsContractAddress := flag.String("utilsContractAddress", "", "")
	playandevolveContractAddress := flag.String("playandevolveContractAddress", "", "")
	shopContractAddress := flag.String("shopContractAddress", "", "")
	trainingpointsContractAddress := flag.String("trainingpointsContractAddress", "", "")
	constantsgettersContractAddress := flag.String("constantsgettersContractAddress", "", "")
	ipfsURL := flag.String("ipfs", "localhost:5001", "ipfs node url")
	delta := flag.Int("delta", 10, "number of block to process at maximum")
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
	if *utilsContractAddress == "" {
		log.Fatal("no utils contract address")
	}
	if *playandevolveContractAddress == "" {
		log.Fatal("no playandevolve contract address")
	}
	if *shopContractAddress == "" {
		log.Fatal("no shop contract address")
	}
	if *trainingpointsContractAddress == "" {
		log.Fatal("no trainingpoints contract address")
	}
	if *constantsgettersContractAddress == "" {
		log.Fatal("no constantsgetters contract address")
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

		contracts, err := contracts.New(
			client,
			*leaguesContractAddress,
			*assetsContractAddress,
			*evolutionContractAddress,
			*engineContractAddress,
			*enginePreCompContractAddress,
			*updatesContractAddress,
			*marketContractAddress,
			*utilsContractAddress,
			*playandevolveContractAddress,
			*shopContractAddress,
			*trainingpointsContractAddress,
			*constantsgettersContractAddress,
		)
		if err != nil {
			return err
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

		processor, err := process.NewEventProcessor(contracts, namesdb, *ipfsURL)
		if err != nil {
			return err
		}
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
