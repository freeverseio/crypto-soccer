package gql

import (
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/worldplayer"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) GetWorldPlayers(args struct{ Input input.GetWorldPlayersInput }) ([]*worldplayer.WorldPlayer, error) {
	log.Debugf("GetWorldPlayers %v", args)

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("not owner of the team")
	}

	return b.createWorldPlayersBatch(b.service, string(args.Input.TeamId))
}

func (b *Resolver) createWorldPlayersBatch(service storage.StorageService, teamId string) ([]*worldplayer.WorldPlayer, error) {
	worldPlayerService := worldplayer.NewWorldPlayerService(b.contracts, b.namesdb)
	players, err := worldPlayerService.CreateBatch(teamId, time.Now().Unix())
	if err != nil {
		return nil, err
	}

	tx, err := b.service.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	sellablePlayers := []*worldplayer.WorldPlayer{}
	for i := range players {
		orders, err := tx.PlayStorePendingOrdersByPlayerId(string(players[i].PlayerId()))
		if err != nil {
			return nil, err
		}
		if len(orders) == 0 {
			sellablePlayers = append(sellablePlayers, players[i])
		}
	}

	return sellablePlayers, nil
}
