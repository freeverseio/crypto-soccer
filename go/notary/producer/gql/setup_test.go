package gql_test

import (
	"crypto/ecdsa"
	"database/sql"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

var bc *testutils.BlockchainNode
var namesdb *names.Generator
var googleCredentials []byte
var db *sql.DB
var service storage.StorageService

var nHumanTeams *big.Int
var offererTeamIdx int64
var sellerTeamIdx int64
var offererTeamId *big.Int
var offerer *ecdsa.PrivateKey
var seller *ecdsa.PrivateKey
var playerId *big.Int

func TestMain(m *testing.M) {
	var err error
	namesdb, err = names.New("../../../names/sql/names.db")
	if err != nil {
		log.Fatal(err)
	}
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

	// We will here assign the next available team to offerer so she can make an offer for the players of a different team
	// we choose that player as a player very far from the current amount of teams (x2)
	nHumanTeams, _ = bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	offererTeamIdx = nHumanTeams.Int64()
	sellerTeamIdx = offererTeamIdx + 1
	offererTeamId, _ = bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(offererTeamIdx))
	offerer, _ = crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	seller, _ = crypto.HexToECDSA("0A878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ = bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(2 + 18*sellerTeamIdx)))

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(offerer.PublicKey),
	)
	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(seller.PublicKey),
	)

	service = postgres.NewStorageService(db)
	os.Exit(m.Run())
}
