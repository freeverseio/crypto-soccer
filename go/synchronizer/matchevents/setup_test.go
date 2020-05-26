package matchevents_test

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/freeverseio/crypto-soccer/go/testutils"

	log "github.com/sirupsen/logrus"
)

var bc *testutils.BlockchainNode
var dump spew.ConfigState

func TestMain(m *testing.M) {
	var err error
	bc, err = testutils.NewBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	dump = spew.ConfigState{DisablePointerAddresses: true, Indent: "\t"}
	os.Exit(m.Run())
}
