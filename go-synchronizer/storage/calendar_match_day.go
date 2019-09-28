package storage

import (
	"math/big"

	log "github.com/sirupsen/logrus"
)

type CalendarMatchDay struct {
	TimezoneIdx   uint8
	CountryIdx    uint16
	LeagueIdx     uint8
	MatchDayIdx   uint8
	HomeTeamID    *big.Int
	VisitorTeamID *big.Int
	HomeGoals     uint8
	VisitorGoals  uint8
}

func (b *Storage) CalendarMatchDayCreate(day CalendarMatchDay) error {
	log.Debugf("[DBMS] Create Match Day %v", day)
	return nil
}
