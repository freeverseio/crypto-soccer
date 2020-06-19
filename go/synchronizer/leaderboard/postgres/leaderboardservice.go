package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/leaderboard"
)

type LeaderboardService struct {
	tx *sql.Tx
}

func NewLeaderboardService(tx *sql.Tx) *LeaderboardService {
	return &LeaderboardService{
		tx: tx,
	}
}

func (b LeaderboardService) Compute(timezone int, country int, league int) (*leaderboard.Leaderboard, error) {
	return nil, nil
}
