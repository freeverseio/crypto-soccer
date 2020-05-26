package contracts_test

import (
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/contracts/truffle"
	log "github.com/sirupsen/logrus"
)

var bc *contracts.Contracts

func TestMain(m *testing.M) {
	var err error
	bc, err = truffle.New()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
