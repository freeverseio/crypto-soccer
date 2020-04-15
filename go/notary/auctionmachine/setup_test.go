package auctionmachine_test

import (
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"gotest.tools/assert"
)

// var db *sql.DB
var bc *testutils.BlockchainNode

func TestMain(m *testing.M) {
	var err error
	// db, err = storage.New("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	bc, err = testutils.NewBlockchainNode()
	if err != nil {
		log.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(bc.Owner.PublicKey),
	)
	os.Exit(m.Run())
}

func newMarket(t *testing.T) *marketpay.MarketPay {
	market, err := marketpay.New()
	assert.NilError(t, err)
	return market
}
