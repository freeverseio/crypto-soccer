package process

import (
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"

	log "github.com/sirupsen/logrus"
)

type LeagueProcessor struct {
}

func NewLeagueProcessor() *LeagueProcessor {
	return &LeagueProcessor{}
}

func (b *LeagueProcessor) Process(event updates.UpdatesActionsSubmission) error {
	day := event.Day
	turnInDay := event.TurnInDay
	if (turnInDay > 1) ||
		(turnInDay == 1 && day != 1) ||
		(turnInDay == 0 && (day < 2 || day > 14)) {
		log.Warnf("[LeagueProcessor] Skipping %v", event)
		return nil
	}

	return nil
}
