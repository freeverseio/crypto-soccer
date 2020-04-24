package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	log "github.com/sirupsen/logrus"
)

func main() {
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable", "postgres url")
	ethereumClient := flag.String("ethereum", "http://localhost:8545", "ethereum node")
	marketContractAddress := flag.String("market_address", "", "market contract address")
	constantsgettersContractAddress := flag.String("constantsgetters_address", "", "constantsgetters contract address")
	privateKeyHex := flag.String("private_key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "private key")
	debug := flag.Bool("debug", false, "print debug logs")
	bufferSize := flag.Int("buffer_size", 10000, "size of event buffer")
	processWait := flag.Int("process_wait", 5, "secs to wait for next process")
	flag.Parse()

	log.Infof("[PARAM] postgres                   : %v", *postgresURL)
	log.Infof("[PARAM] ethereum_client            : %v", *ethereumClient)
	log.Infof("[PARAM] market_address             : %v", *marketContractAddress)
	log.Infof("[PARAM] constantsgetters_address   : %v", *constantsgettersContractAddress)
	privateKey, err := crypto.HexToECDSA(*privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("[PARAM] Address                    : %v", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	log.Infof("[PARAM] Buffer size                : %v", *bufferSize)
	log.Infof("[PARAM] Process wait               : %v", *processWait)
	log.Infof("[PARAM] debug                      : %v", *debug)
	log.Infof("-------------------------------------------------------------------")

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := func() error {
		log.Info("Create the connection to DBMS")
		db, err := storage.New(*postgresURL)
		if err != nil {
			return err
		}
		defer db.Close()
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
		)
		if err != nil {
			return err
		}

		ch := make(chan interface{}, *bufferSize)

		go gql.NewServer(ch, *contracts)
		go producer.NewProcessor(ch, time.Duration(*processWait)*time.Second)

		market := marketpay.NewSandbox()
		cn, err := consumer.New(
			ch,
			market,
			db,
			*contracts,
			privateKey,
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
