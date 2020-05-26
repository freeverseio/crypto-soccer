package contracts_test

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts/truffle"
	log "github.com/sirupsen/logrus"
)

var dump spew.ConfigState
var client *ethclient.Client
var directoryAddress string

func TestMain(m *testing.M) {
	var err error
	client, err = ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}
	directoryAddress, err = truffle.Deploy()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
