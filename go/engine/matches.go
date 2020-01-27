package engine

import (
	"context"
	"database/sql"
	"runtime"

	"golang.org/x/sync/errgroup"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

type Matches []Match

func (b Matches) Play1stHalf(contracts contracts.Contracts) error {
	for _, match := range b {
		if err := match.Play1stHalf(contracts); err != nil {
			log.Error(match.DumpState())
			return err
		}
	}
	return nil
}

func worker(contracts contracts.Contracts, matchesChannel <-chan Match) error {
	c, err := contracts.Duplicate()
	if err != nil {
		return err
	}
	for match := range matchesChannel {
		if err := match.Play1stHalf(*c); err != nil {
			return err
		}
	}
	return nil
}

func (b Matches) Play1stHalfParallel(ctx context.Context, contracts contracts.Contracts) error {
	numWorkers := runtime.NumCPU()
	log.Debugf("Using %v workers", numWorkers)

	matchesChannel := make(chan Match, len(b))
	g, _ := errgroup.WithContext(ctx)

	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			return worker(contracts, matchesChannel)
		})
	}

	for i := 0; i < len(b); i++ {
		matchesChannel <- b[i]
	}
	close(matchesChannel)
	return g.Wait()
}

func (b Matches) ToStorage(contracts *contracts.Contracts, tx *sql.Tx) error {
	for _, match := range b {
		for _, player := range match.HomeTeam.Players {
			var sPlayer storage.Player
			defence, speed, pass, shoot, endurance, _, _, err := contracts.DecodeSkills(player.Skills())
			if err != nil {
				return err
			}
			sPlayer.Defence = defence.Uint64()
			sPlayer.Speed = speed.Uint64()
			sPlayer.Pass = pass.Uint64()
			sPlayer.Shoot = shoot.Uint64()
			sPlayer.Defence = endurance.Uint64()
			sPlayer.EncodedSkills = player.Skills()
			if err = sPlayer.Update(tx); err != nil {
				return err
			}
		}
	}
	return nil
}
