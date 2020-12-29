package gql_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	log "github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

func TestDismissPlayer(t *testing.T) {

	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	in := input.DismissPlayerInput{}

	_, err := r.DismissPlayer(struct{ Input input.DismissPlayerInput }{in})
	assert.Error(t, err, "invalid validUntil")

	in.ValidUntil = "32323"
	_, err = r.DismissPlayer(struct{ Input input.DismissPlayerInput }{in})
	assert.Error(t, err, "invalid playerId")

	in.PlayerId = "32323"
	_, err = r.DismissPlayer(struct{ Input input.DismissPlayerInput }{in})
	assert.Error(t, err, "signature must be 65 bytes long")

	in.Signature = "0f13e4028d911bbf7e305267d593c6b67888030032e73f94a5cf8af204567ab629848e9290568aa5d19c1b7a4761a20ed4059072aacd79bde56e1b52c17a21311b"
	_, err = r.DismissPlayer(struct{ Input input.DismissPlayerInput }{in})
	assert.Error(t, err, "not player owner")
}

func TestDismissPlayerWhenPlayerHasAuction(t *testing.T) {
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	aliceTeamIdx := nHumanTeams.Int64()
	alice, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(7 + 18*aliceTeamIdx)))

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(alice.PublicKey),
	)

	inAuction := storage.Auction{
		ID:     "123abc",
		Seller: crypto.PubkeyToAddress(alice.PublicKey).Hex(),
	}

	mock := mockup.Tx{
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return []storage.Auction{inAuction}, nil },
		RollbackFunc:           func() error { return nil },
		AuctionCancelFunc:      func(ID string) error { return nil },
		CommitFunc:             func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	in := input.DismissPlayerInput{}
	in.ValidUntil = "32323"
	in.PlayerId = playerId.String()

	hash, err := in.Hash()
	assert.NilError(t, err)
	hash = helper.PrefixedHash(hash)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	log.SetLevel(log.DebugLevel)

	_, err = r.DismissPlayer(struct{ Input input.DismissPlayerInput }{in})
	assert.NilError(t, err)
}
