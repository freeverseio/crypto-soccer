package gql

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

type CreateAuctionInput struct {
	Signature  string
	PlayerId   string
	CurrencyId int
	Price      int
	Rnd        int
	ValidUntil string
}

func (b CreateAuctionInput) Hash() ([32]byte, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return [32]byte{}, err
	}
	hash := signer.HashSellMessage(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		validUntil,
		playerId,
	)
	return hash, nil
}

func (b CreateAuctionInput) VerifySignature() (bool, error) {
	hash, err := b.Hash()
	if err != nil {
		return false, err
	}
	return signer.VerifySignature(hash[:], []byte(b.Signature))
}

func (b *Resolver) CreateAuction(args struct{ Input CreateAuctionInput }) (graphql.ID, error) {
	if b.c == nil {
		return graphql.ID(""), errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature()
	if err != nil {
		return graphql.ID(""), err
	}
	if !isValid {
		return graphql.ID(""), errors.New("Invalid signature")
	}

	select {
	case b.c <- args.Input:
	default:
		log.Warning("channel is full")
		return graphql.ID("ciao"), errors.New("channel is full")
	}

	return graphql.ID(args.Input.Signature), nil
}
