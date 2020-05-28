package staker_test

import (
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/testutils"
)

var bc *testutils.BlockchainNode

func TestMain(m *testing.M) {
	var err error
	bc, err = testutils.NewBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
