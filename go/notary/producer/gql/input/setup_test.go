package input_test

import (
	"database/sql"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

var bc *testutils.BlockchainNode
var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	bc, err = testutils.NewBlockchain()
	if err != nil {
		log.Fatal(err)
	}
	db, err = postgres.New("postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
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
