package main

import (
	"flag"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/xolo"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var ethereumClient *string
var privateKeyHex *string

func must(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func serverStart() {

	rpclient, err := ethclient.Dial(*ethereumClient)
	if err != nil {
		log.Fatalf("Failed to connect to RPC: %v %v", *ethereumClient, err)
	}
	privateKey, err := crypto.HexToECDSA(*privateKeyHex)
	if err != nil {
		log.Fatal("Unable to obtain privateKey")
	}
	signer := bind.NewKeyedTransactor(privateKey)
	xserver, err := xolo.NewServer(signer, rpclient)
	if err != nil {
		log.Fatalf("Cannot create server", err)
	}

	engine := gin.Default()
	srv := &http.Server{
		Addr:    ":8004",
		Handler: engine,
	}
	engine.POST("/tx", xserver.HttpPostTx)
	srv.ListenAndServe()
}

func main() {
	privateKeyHex = flag.String("private_key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "private key")
	ethereumClient = flag.String("ethereum", "http://localhost:8545", "ethereum node")
	flag.Parse()
	serverStart()
}
