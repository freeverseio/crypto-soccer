package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

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
	postgresUniverseURL := flag.String("postgresUniverse", "postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable", "postgres universe url")
	postgresGameURL := flag.String("postgresGame", "postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable", "postgres game url")
	debug := flag.Bool("debug", false, "print debug logs")
	bufferSize := flag.Int("buffer_size", 10000, "size of event buffer")
	processWait := flag.Int("process_wait", 5, "secs to wait for next process")
	flag.Parse()

	log.Infof("[PARAM] postgres universe          : %v", *postgresUniverseURL)
	log.Infof("[PARAM] postgres game              : %v", *postgresGameURL)
	log.Infof("[PARAM] Buffer size                : %v", *bufferSize)
	log.Infof("[PARAM] Process wait               : %v", *processWait)
	log.Infof("[PARAM] debug                      : %v", *debug)
	log.Infof("-------------------------------------------------------------------")

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := func() error {
		log.Info("Create the connection to DBMS")
		marketdb, err := postgres.New(*postgresUniverseURL)
		if err != nil {
			return err
		}
		defer marketdb.Close()

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
		contracts, err := contracts.NewByProxyAddress(
			client,
			*proxyAddress,
		)
		if err != nil {
			return err
		}

		googleCredentials, err := ioutil.ReadFile(*googleKey)
		if err != nil {
			return err
		}

		ch := make(chan interface{}, *bufferSize)

		go gql.ListenAndServe(
			ch,
			*contracts,
			namesdb,
			googleCredentials,
			marketdb,
		)
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
			marketdb,
			*contracts,
			privateKey,
			googleCredentials,
			namesdb,
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
