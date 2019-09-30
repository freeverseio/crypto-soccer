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
		log.Warnf("[LeagueProcessor] Skipping %v", event)
		return nil
	}

	countryCount, err := b.storage.CountryInTimezoneCount(timezoneIdx)
	if err != nil {
		return err
	}
	for countryIdx := uint32(0); countryIdx < countryCount; countryIdx++ {
		leagueCount, err := b.storage.LeagueInCountryCount(timezoneIdx, countryIdx)
		if err != nil {
			return err
		}
		for leagueIdx := uint32(0); leagueIdx < leagueCount; leagueIdx++ {
			matches, err := b.storage.GetMatches(timezoneIdx, countryIdx, leagueIdx)
			if err != nil {
				return err
			}
			for matchIdx := 0; matchIdx < len(*matches); matchIdx++ {
				matchSeed := big.NewInt(0) // TODO ??? what's this
				var states [2][25]*big.Int
				var tactics [2]*big.Int
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
					return err
				}
				log.Info(result)
			}
		}
	}

	return nil
}
