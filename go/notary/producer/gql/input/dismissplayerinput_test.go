package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	"gotest.tools/assert"
)

func TestDismissPlayerInputHash(t *testing.T) {
	msg := input.DismissPlayerInput{}
	msg.PlayerId = "123455"
	msg.ValidUntil = "5646456"
	msg.ReturnToAcademy = true

	hash, err := msg.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x26a63dd7a77ba6da621296c5433d235fa802b0eed629457ff3237b321f6db462")

	hash = helper.PrefixedHash(hash)
	assert.Equal(t, hash.Hex(), "0xa345906cc0144e72ba04ea426d34bd486000e51de093b4b1a106deafa21c3244")
}

func TestDismissPlayerSignerAddress(t *testing.T) {
	msg := input.DismissPlayerInput{}
	msg.PlayerId = "123455"
	msg.ValidUntil = "5646456"
	msg.ReturnToAcademy = true
	msg.Signature = "2148732eeca5265898a5fe8dd3ba1c1af5b3d5b815fb23d9d6e383b376a2c91c694170ebd18b64b122905f82d3d6961a78a784b8966fcb350d51c6c5e7917d2d1b"

	address, err := msg.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0xb8CE9ab6943e0eCED004cDe8e3bBed6568B2Fa01")

	r, s, v, err := helper.RSV(msg.Signature)
	assert.Equal(t, hex.EncodeToString(r[:]), "2148732eeca5265898a5fe8dd3ba1c1af5b3d5b815fb23d9d6e383b376a2c91c")
	assert.Equal(t, hex.EncodeToString(s[:]), "694170ebd18b64b122905f82d3d6961a78a784b8966fcb350d51c6c5e7917d2d")
	assert.Equal(t, v, uint8(0x1b))
}

func TestCheckPlayerOnSale(t *testing.T) {
	alice, _ := crypto.HexToECDSA("0B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	inAuction := storage.Auction{
		ID:     "123abc",
		Seller: crypto.PubkeyToAddress(alice.PublicKey).Hex(),
	}

	mock := mockup.Tx{
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return []storage.Auction{inAuction}, nil },
		RollbackFunc:           func() error { return nil },
		CommitFunc:             func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}
	tx, err := service.Begin()
	defer tx.Rollback()
	assert.NilError(t, err)

	in := input.DismissPlayerInput{}
	in.PlayerId = "123455"
	in.ValidUntil = "5646456"
	in.ReturnToAcademy = true
	in.Signature = "2148732eeca5265898a5fe8dd3ba1c1af5b3d5b815fb23d9d6e383b376a2c91c694170ebd18b64b122905f82d3d6961a78a784b8966fcb350d51c6c5e7917d2d1b"

	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	isPlayerOnSale, auctionID, err := in.IsPlayerOnSale(tx)
	assert.NilError(t, err)
	assert.Equal(t, isPlayerOnSale, true)
	assert.Equal(t, auctionID, "123abc")

}

func TestCheckPlayerNotOnSale(t *testing.T) {
	alice, _ := crypto.HexToECDSA("0B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	mock := mockup.Tx{
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return nil, nil },
		RollbackFunc:           func() error { return nil },
		CommitFunc:             func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}
	tx, err := service.Begin()
	defer tx.Rollback()
	assert.NilError(t, err)

	in := input.DismissPlayerInput{}
	in.PlayerId = "123455"
	in.ValidUntil = "5646456"
	in.ReturnToAcademy = true
	in.Signature = "2148732eeca5265898a5fe8dd3ba1c1af5b3d5b815fb23d9d6e383b376a2c91c694170ebd18b64b122905f82d3d6961a78a784b8966fcb350d51c6c5e7917d2d1b"

	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	isPlayerOnSale, auctionID, err := in.IsPlayerOnSale(tx)
	assert.NilError(t, err)
	assert.Equal(t, isPlayerOnSale, false)
	assert.Equal(t, auctionID, "")

}
