package storage

import "math/big"

type Calendar struct {
	TimezoneIdx uint8
	CountryIdx  uint16
	LeagueIdx   uint8
}

type MatchDay struct {
	MatchDayIdx uint8
}

type Match struct {
	HomeTeamID    *big.Int
	VisitorTeamID *big.Int
	HomeGoals     uint8
	VisitorGoals  uint8
}
