package process

import (
	"errors"

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

	return nil
}
