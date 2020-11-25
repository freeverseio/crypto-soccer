package auctionpassmachine_test

import (
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/testutils"

	log "github.com/sirupsen/logrus"
)

// var db *sql.DB
var bc *testutils.BlockchainNode
var googleCredentials []byte
var service storage.StorageService

func TestMain(m *testing.M) {
	var err error
	bc, err = testutils.NewBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	service = postgres.NewStorageService(db)
	os.Exit(m.Run())
}
