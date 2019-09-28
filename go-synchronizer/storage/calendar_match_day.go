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

func (b *Storage) CalendarMatchDayCreate(match CalendarMatchDay) error {
	log.Infof("[DBMS] Create Match Day %v", match)
	_, err := b.db.Exec("INSERT INTO calendar_matches (timezone_idx, country_idx, league_idx, match_day_idx, match_idx) VALUES ($1, $2, $3, $4, $5);",
		match.TimezoneIdx,
		match.CountryIdx,
		match.LeagueIdx,
		match.MatchDayIdx,
		match.MatchIdx,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Storage) GetCalendarMatches(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) (*[]CalendarMatchDay, error) {
	log.Debugf("[DBMS] Get Calendar Matches timezoneIdx %v, countryIdx %v, leagueIdx %v", timezoneIdx, countryIdx, leagueIdx)
	rows, err := b.db.Query("SELECT timezone_idx, country_idx, league_idx, match_day_idx, match_idx FROM calendar_matches WHERE (timezone_idx = $1 AND country_idx = $2 AND league_idx = $1);", timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var matches []CalendarMatchDay
	for rows.Next() {
		var match CalendarMatchDay
		err = rows.Scan(
			&match.TimezoneIdx,
			&match.CountryIdx,
			&match.LeagueIdx,
			&match.MatchDayIdx,
			&match.MatchIdx,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return &matches, nil
}
