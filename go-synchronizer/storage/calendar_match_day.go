package storage

import (
	"math/big"

	log "github.com/sirupsen/logrus"
)

type CalendarMatchDay struct {
	TimezoneIdx   uint8
	CountryIdx    uint32
	LeagueIdx     uint32
	MatchDayIdx   uint8
	MatchIdx      uint8
	HomeTeamID    *big.Int
	VisitorTeamID *big.Int
	HomeGoals     uint8
	VisitorGoals  uint8
}

func (b *Storage) CalendarMatchDayCreate(day CalendarMatchDay) error {
	log.Infof("[DBMS] Create Match Day %v", day)
	return nil
}
