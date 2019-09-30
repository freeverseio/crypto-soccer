package process

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"

	log "github.com/sirupsen/logrus"
)

type LeagueProcessor struct {
	engine  *engine.Engine
	storage *storage.Storage
}

func NewLeagueProcessor(engine *engine.Engine, storage *storage.Storage) *LeagueProcessor {
	return &LeagueProcessor{engine, storage}
}

func (b *LeagueProcessor) Process(event updates.UpdatesActionsSubmission) error {
	day := event.Day
	turnInDay := event.TurnInDay
	timezoneIdx := event.TimeZone
	if timezoneIdx < 1 || timezoneIdx > 24 {
		return errors.New("Wront timezone " + string(timezoneIdx))
	}
	if (turnInDay > 1) ||
		(turnInDay == 1 && day != 1) ||
		(turnInDay == 0 && (day < 2 || day > 14)) {
		log.Warnf("[LeagueProcessor] Skipping timezone %v, day %v, turnInDay %v", timezoneIdx, day, turnInDay)
		return nil
	}
	log.Infof("[LeagueProcessor] Processing timezone %v, day %v, turnInDay %v", timezoneIdx, day, turnInDay)

	countryCount, err := b.storage.CountryInTimezoneCount(timezoneIdx)
	if err != nil {
		return err
	}
	log.Infof("countryCount %v", countryCount)
	for countryIdx := uint32(0); countryIdx < countryCount; countryIdx++ {
		leagueCount, err := b.storage.LeagueInCountryCount(timezoneIdx, countryIdx)
		if err != nil {
			return err
		}
		log.Infof("leagueCount %v", leagueCount)
		for leagueIdx := uint32(0); leagueIdx < leagueCount; leagueIdx++ {
			matches, err := b.storage.GetMatchesInDay(timezoneIdx, countryIdx, leagueIdx, day-1)
			log.Infof("matches count %v", len(matches))
			if err != nil {
				return err
			}
			for matchIdx := 0; matchIdx < len(matches); matchIdx++ {
				matchSeed := big.NewInt(0) // TODO ??? what's this
				var states [2][25]*big.Int
				for i := 0; i < 25; i++ {
					states[0][i] = big.NewInt(0)
					states[1][i] = big.NewInt(0)
				}
				var tactics [2]*big.Int
				tactics[0] = big.NewInt(0)
				tactics[1] = big.NewInt(0)
				is2ndHalf := false
				isHomeStadium := true
				result, err := b.engine.PlayMatch(
					&bind.CallOpts{},
					matchSeed,
					states,
					tactics,
					is2ndHalf,
					isHomeStadium,
				)
				if err != nil {
					log.Fatal(err)
					return err
				}
				log.Infof("result %v - %v", result[0], result[1])
			}
		}
	}

	return nil
}
