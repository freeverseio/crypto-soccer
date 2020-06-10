package playstore_test

import (
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/testutils"

	log "github.com/sirupsen/logrus"
)

// var db *sql.DB
var bc *testutils.BlockchainNode
var googleCredentials []byte

func TestMain(m *testing.M) {
	var err error
	bc, err = testutils.NewBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
