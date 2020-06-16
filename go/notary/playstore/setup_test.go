package playstore_test

import (
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/testutils"

	log "github.com/sirupsen/logrus"
)

// var db *sql.DB
var bc *testutils.BlockchainNode
var googleCredentials []byte
var namesdb *names.Generator

func TestMain(m *testing.M) {
	var err error
	bc, err = testutils.NewBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	namesdb, err = names.New("../../names/sql/names.db")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
