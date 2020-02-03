package engine

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"runtime"

	"golang.org/x/sync/errgroup"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	log "github.com/sirupsen/logrus"
)

type Matches []Match

func NewMatchesFromTimezoneIdxMatchdayIdx(
	tx *sql.Tx,
	timezoneIdx uint8,
	day uint8,
) (*Matches, error) {
	stoMatches, err := storage.MatchesByTimezoneIdxAndMatchDay(tx, timezoneIdx, day)
	if err != nil {
		return nil, err
	}

	var matches Matches
	for _, stoMatch := range stoMatches {
		stoHomeTeam, err := storage.TeamByTeamId(tx, stoMatch.HomeTeamID)
		if err != nil {
			return nil, err
		}
		stoVisitorTeam, err := storage.TeamByTeamId(tx, stoMatch.VisitorTeamID)
		if err != nil {
			return nil, err
		}
		stoHomePlayers, err := storage.PlayersByTeamId(tx, stoMatch.HomeTeamID)
		if err != nil {
			return nil, err
		}
		stoVisitorPlayers, err := storage.PlayersByTeamId(tx, stoMatch.VisitorTeamID)
		if err != nil {
			return nil, err
		}
		match := NewMatchFromStorage(
			stoMatch,
			stoHomeTeam,
			stoVisitorTeam,
			stoHomePlayers,
			stoVisitorPlayers,
		)
		matches = append(matches, *match)
	}
	return &matches, nil
}

func NewMatchesFromTimezoneIdxCountryIdxLeagueIdxMatchdayIdx(
	tx *sql.Tx,
	timezoneIdx uint8,
	countryIdx uint32,
	leagueIdx uint32,
	day uint8,
) (*Matches, error) {
	stoMatches, err := storage.MatchesByTimezoneIdxCountryIdxLeagueIdxMatchdayIdx(tx, timezoneIdx, countryIdx, leagueIdx, day)
	if err != nil {
		return nil, err
	}

	var matches Matches
	for _, stoMatch := range stoMatches {
		stoHomeTeam, err := storage.TeamByTeamId(tx, stoMatch.HomeTeamID)
		if err != nil {
			return nil, err
		}
		stoVisitorTeam, err := storage.TeamByTeamId(tx, stoMatch.VisitorTeamID)
		if err != nil {
			return nil, err
		}
		stoHomePlayers, err := storage.PlayersByTeamId(tx, stoMatch.HomeTeamID)
		if err != nil {
			return nil, err
		}
		stoVisitorPlayers, err := storage.PlayersByTeamId(tx, stoMatch.VisitorTeamID)
		if err != nil {
			return nil, err
		}
		match := NewMatchFromStorage(
			stoMatch,
			stoHomeTeam,
			stoVisitorTeam,
			stoHomePlayers,
			stoVisitorPlayers,
		)
		matches = append(matches, *match)
	}
	return &matches, nil
}

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
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

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
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

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

func (b Matches) ToStorage(contracts contracts.Contracts, tx *sql.Tx) error {
	for _, match := range b {
		if err := match.ToStorage(contracts, tx); err != nil {
			return err
		}

	}
	return nil
}

func (b Matches) DumpState() string {
	var state string
	for i, match := range b {
		state += fmt.Sprintf("Match: %v\n", i)
		state += match.DumpState()
	}
	return state
}
