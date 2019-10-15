package storage

import (
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Tactic struct {
	TeamID          *big.Int
	Defense         uint8
	Center          uint8
	Attack          uint8
	PlayerPositions []uint8
}

func (b *Storage) TacticCreate(t Tactic) error {
	log.Debugf("[DBMS] Create tactic %v", t)
	_, err := b.db.Exec("INSERT INTO tactics (team_id, defense, center, attack) VALUES ($1, $2, $3, $4);",
		t.TeamID.String(),
		t.Defense,
		t.Center,
		t.Attack,
	)
	return err
}
