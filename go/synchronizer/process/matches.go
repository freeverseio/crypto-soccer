package process

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"runtime"
	"time"

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
		stoHomePlayers, err := storage.ActivePlayersByTeamId(tx, stoMatch.HomeTeamID.String())
		if err != nil {
			return nil, err
		}
		stoVisitorPlayers, err := storage.ActivePlayersByTeamId(tx, stoMatch.VisitorTeamID.String())
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

func (b *Matches) Play1stHalfParallel(ctx context.Context, contracts contracts.Contracts) error {
	numWorkers := runtime.NumCPU()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	matchesChannel := make(chan *engine.Match, len(*b))
	g, _ := errgroup.WithContext(ctx)

	start := time.Now()
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
	if err := g.Wait(); err != nil {
		return err
	}

	elapsed := time.Now().Sub(start)
	log.Infof("[precessor|1stHalfParallelProcess] %v workers, took %v secs", numWorkers, elapsed.Seconds())

	return nil
}

func (b *Matches) Play2ndHalfParallel(ctx context.Context, contracts contracts.Contracts) error {
	numWorkers := runtime.NumCPU()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	matchesChannel := make(chan *engine.Match, len(*b))
	g, _ := errgroup.WithContext(ctx)

	start := time.Now()
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
	if err := g.Wait(); err != nil {
		return err
	}

	elapsed := time.Now().Sub(start)
	log.Infof("[precessor|2ndHalfParallelProcess] %v workers, took %v secs", numWorkers, elapsed.Seconds())

	return nil
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

func (b Matches) ToStorage(contracts contracts.Contracts, tx *sql.Tx, blockNumber uint64, ctx context.Context) error {

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

func (b Matches) ToStorageMaxParallelTx(contracts contracts.Contracts, blockNumber uint64, ctx context.Context, universedb *sql.DB) error {

	errs, _ := errgroup.WithContext(ctx)

	for _, match := range b {
		matchRoutine := match
		errs.Go(func() error {
			log.Infof("Processing match %x", matchRoutine.Hash())
			tx, err := universedb.Begin()
			if err != nil {
				return err
			}
			err = matchRoutine.ToStorage(contracts, tx, blockNumber)
			if err != nil {
				filename := fmt.Sprintf("/tmp/%x", matchRoutine.Hash()) + ".toStorage.error.json"
				log.Errorf("match to storage: %v: saving match state to %v", err.Error(), filename)
				if err := ioutil.WriteFile(filename, matchRoutine.ToJson(), 0644); err != nil {
					tx.Rollback()
					return err
				}
				tx.Rollback()
				return err
			}
			return tx.Commit()
		})
	}
	return errs.Wait()
}

func (b Matches) ToStorageLimitedParallelTx(contracts contracts.Contracts, blockNumber uint64, ctx context.Context, universedb *sql.DB) error {

	errs, _ := errgroup.WithContext(ctx)
	maxClients := 32
	for v := 0; v < len(b)/maxClients; v++ {
		matchV := v
		for i := 0; i < maxClients; i++ {
			matchI := i
			log.Infof("Processing matchroutine %v", matchV*maxClients+matchI)
			matchRoutine := b[matchV*maxClients+matchI]
			errs.Go(func() error {
				tx, err := universedb.Begin()
				if err != nil {
					return err
				}
				err = matchRoutine.ToStorage(contracts, tx, blockNumber)
				if err != nil {
					filename := fmt.Sprintf("/tmp/%x", matchRoutine.Hash()) + ".toStorage.error.json"
					log.Errorf("match to storage: %v: saving match state to %v", err.Error(), filename)
					if err := ioutil.WriteFile(filename, matchRoutine.ToJson(), 0644); err != nil {
						tx.Rollback()
						return err
					}
					tx.Rollback()
					return err
				}
				return tx.Commit()
			})
			errs.Wait()
		}
		// time.Sleep(2 * time.Second)

	}
	return nil
}

func (b Matches) ToStorageBulk(contracts contracts.Contracts, tx *sql.Tx, blockNumber uint64) error {

	var teamsToUpdate []storage.Team
	var teamsHistoriesToInsert []*storage.TeamHistory
	var playersToUpdate []storage.Player
	var playersHistoriesToInsert []*storage.PlayerHistory
	var matchesToUpdate []storage.Match
	var matchesHistoriesToInsert []*storage.MatchHistory
	var eventsToUpdate []*storage.MatchEvent

	for _, match := range b {
		teamsToUpdate = append(teamsToUpdate, match.HomeTeam.Team)
		teamsToUpdate = append(teamsToUpdate, match.VisitorTeam.Team)

		teamsHistoriesToInsert = append(teamsHistoriesToInsert, storage.NewTeamHistory(blockNumber, match.HomeTeam.Team))
		teamsHistoriesToInsert = append(teamsHistoriesToInsert, storage.NewTeamHistory(blockNumber, match.VisitorTeam.Team))

		for _, player := range match.HomeTeam.Players {
			if player.IsNil() {
				continue
			}

			playersToUpdate = append(playersToUpdate, player.Player)
			playersHistoriesToInsert = append(playersHistoriesToInsert, storage.NewPlayerHistory(blockNumber, player.Player))
		}
		for _, player := range match.VisitorTeam.Players {
			if player.IsNil() {
				continue
			}

			playersToUpdate = append(playersToUpdate, player.Player)
			playersHistoriesToInsert = append(playersHistoriesToInsert, storage.NewPlayerHistory(blockNumber, player.Player))
		}
		matchesToUpdate = append(matchesToUpdate, match.Match)
		matchesHistoriesToInsert = append(matchesHistoriesToInsert, storage.NewMatchHistory(blockNumber, match.Match))
		newEventsToUpdate, err := match.ToStorageBulkReturn(contracts, tx, blockNumber)
		if err != nil {
			filename := fmt.Sprintf("/tmp/%x", match.Hash()) + ".toStorage.error.json"
			log.Errorf("match to storage: %v: saving match state to %v", err.Error(), filename)
			if err := ioutil.WriteFile(filename, match.ToJson(), 0644); err != nil {
				return err
			}
			return err
		}
		eventsToUpdate = append(eventsToUpdate, newEventsToUpdate[:]...)
	}

	log.Infof("Num teams to update %v", len(teamsToUpdate))
	// for _, team := range teamsToUpdate {
	// 	log.Infof("Team id %v", team.TeamID)
	// }
	log.Infof("Num teams histories to insert %v", len(teamsHistoriesToInsert))
	// for _, team := range teamsHistoriesToInsert {
	// 	log.Infof("Team id %v", team.TeamID)
	// }
	log.Infof("Num players to update %v", len(playersToUpdate))
	// for _, player := range playersToUpdate {
	// 	log.Infof("player id %v", player.PlayerId)
	// }
	log.Infof("Num players histories to insert %v", len(playersHistoriesToInsert))
	// for _, player := range playersHistoriesToInsert {
	// 	log.Infof("player id %v", player.PlayerId)
	// }
	log.Infof("Num matches to update %v", len(matchesToUpdate))
	log.Infof("Num events to insert %v", len(eventsToUpdate))

	err := storage.TeamsBulkInsertUpdate(teamsToUpdate, tx)
	log.Infof("Updated all teams")
	if err != nil {
		return err
	}

	i := 0
	for i < len(playersToUpdate) {
		// log.Infof("Updating from %v to %v", i, i+3000)
		newI := i + 3000
		if newI > len(playersToUpdate) {
			newI = len(playersToUpdate)
		}
		err := storage.PlayersBulkInsertUpdate(playersToUpdate[i:newI], tx)
		if err != nil {
			return err
		}
		i = newI

	}
	log.Infof("Updated all players")

	err = storage.MatchesBulkInsertUpdate(matchesToUpdate, tx)
	log.Infof("Updated all matches")
	if err != nil {
		return err
	}

	if len(eventsToUpdate) > 0 {
		v := 0
		for v < len(eventsToUpdate) {
			newV := v + 3000
			if newV > len(eventsToUpdate) {
				newV = len(eventsToUpdate)
			}
			if err := storage.MatchEventsBulkInsert(eventsToUpdate[v:newV], tx); err != nil {
				return err
			}
			v = newV

		}
		log.Infof("Updated all events")
	}

	err = storage.TeamsHistoriesBulkInsert(teamsHistoriesToInsert, tx)
	if err != nil {
		return err
	}
	log.Infof("Inserted all teams histories")
	x := 0
	for x < len(playersHistoriesToInsert) {
		newX := x + 3000
		if newX > len(playersHistoriesToInsert) {
			newX = len(playersHistoriesToInsert)
		}
		err = storage.PlayersHistoriesBulkInsert(playersHistoriesToInsert[x:newX], tx)
		if err != nil {
			return err
		}
		x = newX
	}
	log.Infof("Inserted all players histories")

	w := 0
	for w < len(matchesHistoriesToInsert) {
		newW := w + 3000
		if newW > len(matchesHistoriesToInsert) {
			newW = len(matchesHistoriesToInsert)
		}
		err = storage.MatchesHistoriesBulkInsert(matchesHistoriesToInsert[w:newW], tx)
		if err != nil {
			return err
		}
		w = newW
	}
	log.Infof("Inserted all matches histories")
	return nil
}
