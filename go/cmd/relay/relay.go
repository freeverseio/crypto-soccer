package main

import (
	"flag"
	"math/big"
	"time"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/useractions/ipfscluster"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/relay/consumer"
	"github.com/freeverseio/crypto-soccer/go/relay/producer"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

func main() {
	postgresURL := flag.String("postgres", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable", "postgres url")
	debug := flag.Bool("debug", false, "print debug logs")
	ethereumClient := flag.String("ethereum", "http://localhost:8545", "ethereum node")
	proxyAddress := flag.String("proxy_address", "", "")
	privateKeyHex := flag.String("private_key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "private key")
	ipfsURL := flag.String("ipfs", "localhost:5001", "ipfs node url")
	bufferSize := flag.Int("buffer_size", 10000, "size of event buffer")
	flag.Parse()

	log.Infof("[PARAM] proxy_address              : %v", *proxyAddress)
	log.Infof("[PARAM] ipfs url                   : %v", *ipfsURL)

	privateKey, err := crypto.HexToECDSA(*privateKeyHex)
	if err != nil {
		log.Fatal("Unable to obtain privateLey")
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
	auth := bind.NewKeyedTransactor(privateKey)
	auth.GasPrice = big.NewInt(10000000000) // in xdai is fixe to 3 GWei
	log.Infof("Address : %v", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())

	bc, err := contracts.NewByProxyAddress(client, *proxyAddress)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("ipfs URL: %v", *ipfsURL)

	log.Info("Connecting to DBMS: ", *postgresURL)
	db, err := storage.New(*postgresURL)
	if err != nil {
		log.Fatal("Failed to connect to DBMS: %v", err)
	}

	ch := make(chan interface{}, *bufferSize)

	go gql.NewServer(ch, *bc, db)
	go producer.NewSubmitUserActionsTimer(ch, 5*time.Second)

	consumer.NewConsumer(
		ch,
		client,
		auth,
		*bc,
		ipfscluster.NewUserActionsPublishService(*ipfsURL),
		db,
	).Start()
}
