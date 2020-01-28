package engine

import (
	"context"
	"database/sql"
	"runtime"

	"golang.org/x/sync/errgroup"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	log "github.com/sirupsen/logrus"
)

type Matches []Match

func (b Matches) Play1stHalf(ctx context.Context, contracts contracts.Contracts) error {
	for _, match := range b {
		if err := match.Play1stHalf(contracts); err != nil {
			log.Error(match.DumpState())
			return err
		}
	}
	return nil
}

func (b Matches) Play2ndHalf(ctx context.Context, contracts contracts.Contracts) error {
	for _, match := range b {
		if err := match.Play2ndHalf(contracts); err != nil {
			log.Error(match.DumpState())
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
			c, err := contracts.Clone()
			if err != nil {
				return err
			}
			for match := range matchesChannel {
				if err := match.Play1stHalf(*c); err != nil {
					return err
				}
			}
			return nil
		})
	}

	for i := 0; i < len(b); i++ {
		matchesChannel <- b[i]
	}
	close(matchesChannel)
	return g.Wait()
}

func (b Matches) Play2ndHalfParallel(ctx context.Context, contracts contracts.Contracts) error {
	numWorkers := runtime.NumCPU()
	log.Debugf("Using %v workers", numWorkers)

	matchesChannel := make(chan Match, len(b))
	g, _ := errgroup.WithContext(ctx)

	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			c, err := contracts.Clone()
			if err != nil {
				return err
			}
			for match := range matchesChannel {
				if err := match.Play2ndHalf(*c); err != nil {
					return err
				}
			}
			return nil
		})
	}

	for i := 0; i < len(b); i++ {
		matchesChannel <- b[i]
	}
	close(matchesChannel)
	return g.Wait()
}

func (b Matches) ToStorage(tx *sql.Tx) error {
	for _, match := range b {
		if err := match.ToStorage(tx); err != nil {
			return err
		}

	}
	return nil
}
