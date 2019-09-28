package process

import "github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"

type Calendar struct {
	leagues *leagues.Leagues
}

func NewCalendar(leagues *leagues.Leagues) *Calendar {
	return &Calendar{leagues}
}

func (b *Calendar) Generate(timezoneIdx uint8, countryIdx uint16, leagueIdx uint16) error {
	return nil
}
