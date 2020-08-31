package gql_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	"gotest.tools/assert"
)

func TestCreateAuctionCallRollbackOnError(t *testing.T) {
	rollbackCounter := 0

	mock := mockup.Tx{
		AuctionInsertFunc: func(auction storage.Auction) error { return errors.New("error") },
		RollbackFunc:      func() error { rollbackCounter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CreateAuctionInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = "274877906948"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	playerId, _ := new(big.Int).SetString(in.PlayerId, 10)
	validUntil, err := strconv.ParseInt(in.ValidUntil, 10, 64)
	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntil,
		playerId,
	)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CreateAuction(struct{ Input input.CreateAuctionInput }{in})
	assert.Error(t, err, "error")
	assert.Equal(t, rollbackCounter, 1)
}

func TestCreateAuctionCallCommit(t *testing.T) {
	counter := 0

	mock := mockup.Tx{
		AuctionInsertFunc: func(auction storage.Auction) error { return nil },
		CommitFunc:        func() error { counter++; return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	in := input.CreateAuctionInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = "274877906948"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	playerId, _ := new(big.Int).SetString(in.PlayerId, 10)
	validUntil, err := strconv.ParseInt(in.ValidUntil, 10, 64)
	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntil,
		playerId,
	)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CreateAuction(struct{ Input input.CreateAuctionInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, counter, 1)
}

func TestCreateAuctionReturnIDOfTheAuction(t *testing.T) {
	mock := mockup.Tx{
		AuctionInsertFunc: func(auction storage.Auction) error { return nil },
		CommitFunc:        func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	r := gql.NewResolver(make(chan interface{}, 10), *bc.Contracts, namesdb, googleCredentials, service)

	in := input.CreateAuctionInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = "274877906948"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

	playerId, _ := new(big.Int).SetString(in.PlayerId, 10)
	validUntil, err := strconv.ParseInt(in.ValidUntil, 10, 64)
	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntil,
		playerId,
	)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	id, err := in.ID()
	assert.NilError(t, err)
	result, err := r.CreateAuction(struct{ Input input.CreateAuctionInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, id, result)
}
