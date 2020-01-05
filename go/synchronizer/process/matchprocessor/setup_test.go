package matchprocessor_test

import (
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/testutils"

	log "github.com/sirupsen/logrus"
)

var bc *testutils.BlockchainNode

func TestMain(m *testing.M) {
	var err error
	bc, err = testutils.NewBlockchainNode()
	if err != nil {
		log.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
	os.Exit(m.Run())
}
