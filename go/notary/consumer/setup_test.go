package consumer_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

// var db *sql.DB
var bc *testutils.BlockchainNode
var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = storage.New("postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	bc, err = testutils.NewBlockchainNode()
	if err != nil {
		log.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
	os.Exit(m.Run())
}

func newMarket(t *testing.T) *marketpay.MarketPay {
	market := marketpay.New()
	return market
}
