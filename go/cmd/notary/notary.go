package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	log "github.com/sirupsen/logrus"
)

func main() {
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable", "postgres url")
	ethereumClient := flag.String("ethereum", "http://localhost:8545", "ethereum node")
	namesDatabase := flag.String("namesDatabase", "./names.db", "name database path")
	marketContractAddress := flag.String("market_address", "", "market contract address")
	constantsgettersContractAddress := flag.String("constantsgetters_address", "", "constantsgetters contract address")
	privilegedContractAddress := flag.String("privileged_address", "", "privileged contract address")
	privateKeyHex := flag.String("private_key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "private key")
	debug := flag.Bool("debug", false, "print debug logs")
	bufferSize := flag.Int("buffer_size", 10000, "size of event buffer")
	processWait := flag.Int("process_wait", 5, "secs to wait for next process")
	marketID := flag.String("market_id", "", "WARNING: market identifier. If set connecting the real market")
	googleKey := flag.String("google_key", "", "google credentials")
	iapTestOn := flag.Bool("iap_test", false, "allow purchase of testing iap players")
	flag.Parse()

	log.Infof("[PARAM] postgres                   : %v", *postgresURL)
	log.Infof("[PARAM] ethereum_client            : %v", *ethereumClient)
	log.Infof("[PARAM] market_address             : %v", *marketContractAddress)
	log.Infof("[PARAM] constantsgetters_address   : %v", *constantsgettersContractAddress)
	log.Infof("[PARAM] privileged_address         : %v", *privilegedContractAddress)
	privateKey, err := crypto.HexToECDSA(*privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("[PARAM] Address                    : %v", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	log.Infof("[PARAM] Buffer size                : %v", *bufferSize)
	log.Infof("[PARAM] Process wait               : %v", *processWait)
	log.Infof("[PARAM] debug                      : %v", *debug)
	if *marketID == "" {
		log.Infof("[PARAM] market                     : sandbox")
	} else {
		log.Infof("[PARAM] market                     : REAL")
	}
	log.Infof("[PARAM] google credentials         : %v", *googleKey)
	log.Infof("[PARAM] iap test                   : %v", *iapTestOn)
	log.Infof("-------------------------------------------------------------------")

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if _, err := os.Stat(*namesDatabase); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("no names database file at %v", *namesDatabase)
		}
	}

	if err := func() error {
		log.Info("Create the connection to DBMS")
		db, err := postgres.New(*postgresURL)
		if err != nil {
			return err
		}
		defer db.Close()

		namesdb, err := names.New(*namesDatabase)
		if err != nil {
			return err
		}
		defer namesdb.Close()

		log.Info("Dial the Ethereum client: ", *ethereumClient)
		client, err := ethclient.Dial(*ethereumClient)
		if err != nil {
			return err
		}
		defer client.Close()
		contracts, err := contracts.New(
			client,
			"", "", "", "", "", "",
			*marketContractAddress,
			"", "", "", "",
			*constantsgettersContractAddress,
			*privilegedContractAddress,
			"",
		)
		if err != nil {
			return err
		}

		googleCredentials, err := ioutil.ReadFile(*googleKey)
		if err != nil {
			return err
		}

		ch := make(chan interface{}, *bufferSize)

		go gql.NewServer(ch, *contracts, namesdb, googleCredentials)
		go producer.NewProcessor(ch, time.Duration(*processWait)*time.Second)
		go producer.NewPlaystoreOrderEventProcessor(ch, time.Duration(*processWait)*time.Second)

		var market marketpay.IMarketPay
		if *marketID == "" {
			market = marketpay.NewSandbox()
		} else {
			market = marketpay.New(*marketID)
		}

		cn, err := consumer.New(
			ch,
			market,
			db,
			*contracts,
			privateKey,
			googleCredentials,
			*iapTestOn,
		)
		if err != nil {
			return err
		}
		cn.Start()

		return nil
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
