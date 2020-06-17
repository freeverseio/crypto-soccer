package process

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"runtime"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
	sto "github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
	log "github.com/sirupsen/logrus"
)

type Matches []engine.Match

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
		stoHomeTeam, err := storage.TeamByTeamId(tx, stoMatch.HomeTeamID.String())
		if err != nil {
			return nil, err
		}
		stoVisitorTeam, err := storage.TeamByTeamId(tx, stoMatch.VisitorTeamID.String())
		if err != nil {
			return nil, err
		}
		stoHomePlayers, err := storage.PlayersByTeamId(tx, stoMatch.HomeTeamID.String())
		if err != nil {
			return nil, err
		}
		stoVisitorPlayers, err := storage.PlayersByTeamId(tx, stoMatch.VisitorTeamID.String())
		if err != nil {
			return nil, err
		}
		match := engine.NewMatchFromStorage(
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

func (b *Matches) Play1stHalf(ctx context.Context, contracts contracts.Contracts) error {
	log.Info("Using no paraller Play1stHalf")
	for i := range *b {
		if err := (*b)[i].Play1stHalf(contracts); err != nil {
			return err
		}
	}
	return nil
}

func (b *Matches) Play1stHalfParallel(ctx context.Context, contracts contracts.Contracts) error {
	numWorkers := runtime.NumCPU()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	matchesChannel := make(chan *engine.Match, len(*b))
	g, _ := errgroup.WithContext(ctx)

	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			c, err := contracts.Clone()
			if err != nil {
				return err
			}
			for match := range matchesChannel {
				if err := match.Play1stHalf(*c); err != nil {
					filename := fmt.Sprintf("/tmp/%x", match.Hash()) + ".1st.error.json"
					log.Errorf("play 1st half: %v: saving match state to %v", err.Error(), filename)
					if err := ioutil.WriteFile(filename, match.ToJson(), 0644); err != nil {
						return err
					}
				}
			}
			return nil
		})
	}

	for i := 0; i < len(*b); i++ {
		matchesChannel <- &(*b)[i]
	}
	close(matchesChannel)
	return g.Wait()
}

func (b *Matches) Play2ndHalf(ctx context.Context, contracts contracts.Contracts) error {
	log.Info("Using no paraller Play2ndHalf")
	for i := range *b {
		if err := (*b)[i].Play2ndHalf(contracts); err != nil {
			return err
		}
	}
	return nil
}

func (b *Matches) Play2ndHalfParallel(ctx context.Context, contracts contracts.Contracts) error {
	numWorkers := runtime.NumCPU()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	matchesChannel := make(chan *engine.Match, len(*b))
	g, _ := errgroup.WithContext(ctx)

	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			c, err := contracts.Clone()
			if err != nil {
				return err
			}
			for match := range matchesChannel {
				if err := match.Play2ndHalf(*c); err != nil {
					filename := fmt.Sprintf("/tmp/%x", match.Hash()) + ".2nd.error.json"
					log.Errorf("play 2nd half: %v: saving match state to %v", err.Error(), filename)
					if err := ioutil.WriteFile(filename, match.ToJson(), 0644); err != nil {
						return err
					}
				}
			}
			return nil
		})
	}

	for i := 0; i < len(*b); i++ {
		matchesChannel <- &(*b)[i]
	}
	close(matchesChannel)
	return g.Wait()
}

func (b *Matches) SetSeed(seed [32]byte) {
	for i := range *b {
		(*b)[i].Seed = seed
	}
}

func (b *Matches) SetStartTime(startTime *big.Int) {
	for i := range *b {
		(*b)[i].StartTime = startTime
	}
}

func Minute2Round(minute int) uint8 {
	if minute >= 90 {
		return 11
	}
	if minute > 45 {
		minute -= 45
	}
	mapping := [10]uint8{0, 1, 2, 3, 5, 6, 7, 8, 10, 11}
	idx := int(float32(minute) / 5)
	return uint8(mapping[idx])
}

func (b *Matches) SetTactics(contracts contracts.Contracts, tactics []sto.Tactic) error {
	for _, tactic := range tactics {
		substitutions := [3]uint8{
			uint8(tactic.Substitution0Target),
			uint8(tactic.Substitution1Target),
			uint8(tactic.Substitution2Target),
		}
		substitutionsMinute := [3]uint8{
			Minute2Round(tactic.Substitution0Minute),
			Minute2Round(tactic.Substitution1Minute),
			Minute2Round(tactic.Substitution2Minute),
		}
		formation := [14]uint8{
			uint8(tactic.Shirt0),
			uint8(tactic.Shirt1),
			uint8(tactic.Shirt2),
			uint8(tactic.Shirt3),
			uint8(tactic.Shirt4),
			uint8(tactic.Shirt5),
			uint8(tactic.Shirt6),
			uint8(tactic.Shirt7),
			uint8(tactic.Shirt8),
			uint8(tactic.Shirt9),
			uint8(tactic.Shirt10),
			uint8(tactic.Substitution0Shirt),
			uint8(tactic.Substitution1Shirt),
			uint8(tactic.Substitution2Shirt),
		}
		extraAttack := [10]bool{
			tactic.ExtraAttack1,
			tactic.ExtraAttack2,
			tactic.ExtraAttack3,
			tactic.ExtraAttack4,
			tactic.ExtraAttack5,
			tactic.ExtraAttack6,
			tactic.ExtraAttack7,
			tactic.ExtraAttack8,
			tactic.ExtraAttack9,
			tactic.ExtraAttack10,
		}
		tacticID := uint8(tactic.TacticID)
		encodedTactic, err := contracts.Engine.EncodeTactics(
			&bind.CallOpts{},
			substitutions,
			substitutionsMinute,
			formation,
			extraAttack,
			tacticID,
		)
		if err != nil {
			log.Errorf("%v %v", err.Error(), tactic)
			continue
		}
		for i := range *b {
			if tactic.TeamID == (*b)[i].HomeTeam.TeamID {
				(*b)[i].HomeTeam.Tactic = encodedTactic.String()
			}
			if tactic.TeamID == (*b)[i].VisitorTeam.TeamID {
				(*b)[i].VisitorTeam.Tactic = encodedTactic.String()
			}
		}
	}
	return nil
}

func (b *Matches) SetTrainings(contracts contracts.Contracts, trainings []sto.Training) error {
	for _, training := range trainings {
		for i := range *b {
			if training.TeamID == (*b)[i].HomeTeam.TeamID {
				(*b)[i].HomeTeam.Training = *engine.NewTraining(training)
			}
			if training.TeamID == (*b)[i].VisitorTeam.TeamID {
				(*b)[i].VisitorTeam.Training = *engine.NewTraining(training)
			}
		}
	}
	return nil
}

func (b Matches) ToStorage(contracts contracts.Contracts, tx *sql.Tx, blockNumber uint64) error {
	for _, match := range b {
		if err := match.ToStorage(contracts, tx, blockNumber); err != nil {
			filename := fmt.Sprintf("/tmp/%x", match.Hash()) + ".toStorage.error.json"
			log.Errorf("match to storage: %v: saving match state to %v", err.Error(), filename)
			if err := ioutil.WriteFile(filename, match.ToJson(), 0644); err != nil {
				return err
			}
			return err
		}
	}
	return nil
}
