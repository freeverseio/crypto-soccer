package gql

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/notary/worldplayer"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) GetWorldPlayers(args struct{ Input input.GetWorldPlayersInput }) ([]*worldplayer.WorldPlayer, error) {
	log.Debugf("GetWorldPlayers %v", args)

	hash, err := args.Input.Hash()
	if err != nil {
		return nil, err
	}
	sign, err := hex.DecodeString(args.Input.Signature)
	if err != nil {
		return nil, err
	}

	isValid, err := helper.VerifySignature(hash, sign)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, errors.New("Invalid signature")
	}

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("not owner of the team")
	}

	return b.createWorldPlayersBatch(string(args.Input.TeamId))
}

func (b *Resolver) createWorldPlayersBatch(teamId string) ([]*worldplayer.WorldPlayer, error) {
	worldPlayerService := worldplayer.NewWorldPlayerService(b.contracts, b.namesdb)
	players, err := worldPlayerService.CreateBatch(teamId, time.Now().Unix())
	if err != nil {
		return nil, err
	}

	tx, err := b.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	storageService := postgres.NewPlaystoreOrderService(tx)

	sellablePlayers := []*worldplayer.WorldPlayer{}
	for i := range players {
		orders, err := storageService.PendingOrdersByPlayerId(string(players[i].PlayerId()))
		if err != nil {
			return nil, err
		}
		if len(orders) == 0 {
			sellablePlayers = append(sellablePlayers, players[i])
		}
	}

	return sellablePlayers, nil
}
