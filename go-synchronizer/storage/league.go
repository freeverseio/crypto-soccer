package storage

import (
	"math/big"

	log "github.com/sirupsen/logrus"
)

type League struct {
	TimezoneIdx uint8
	CountryIdx  uint16
	LeagueIdx   uint8
	TeamID      *big.Int
	Points      uint8
}

func (b *Storage) LeagueCount() (uint64, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM leagues;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *Storage) LeagueCreate(league League) error {
	log.Debugf("[DBMS] Create league %v", league)
	_, err := b.db.Exec("INSERT INTO leagues (timezone_idx, country_idx, league_idx, team_id, points) VALUES ($1, $2, $3, $4, $5);",
		league.TimezoneIdx,
		league.CountryIdx,
		league.LeagueIdx,
		league.TeamID.String(),
		league.Points,
	)
	if err != nil {
		return err
	}
	return nil
}

// func (b *Storage) GetLeague(id uint64) (League, error) {
// 	league := League{}
// 	rows, err := b.db.Query("SELECT id FROM leagues WHERE (id = $1);", id)
// 	if err != nil {
// 		return league, err
// 	}
// 	defer rows.Close()
// 	if !rows.Next() {
// 		return league, errors.New("Unexistent league")
// 	}
// 	err = rows.Scan(
// 		&league.Id,
// 	)
// 	if err != nil {
// 		return league, err
// 	}
// 	return league, nil
// }
