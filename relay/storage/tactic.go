package storage

import (
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Tactic struct {
	TeamID      *big.Int
	Defense     uint8
	Center      uint8
	Attack      uint8
	Shirts      [11]uint8
	ExtraAttack [10]uint8
}

func (b *Storage) TacticCreate(t Tactic) error {
	log.Debugf("[DBMS] Create tactic %v", t)
	_, err := b.db.Exec("INSERT INTO tactics (team_id, defense, center, attack, shirt_0, shirt_1, shirt_2, shirt_3, shirt_4, shirt_5, shirt_6, shirt_7, shirt_8, shirt_9, shirt_10, extra_attack_1, extra_attack_2, extra_attack_3, extra_attack_4, extra_attack_5, extra_attack_6, extra_attack_7, extra_attack_8, extra_attack_9, extra_attack_10) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25);",
		t.TeamID.String(),
		t.Defense,
		t.Center,
		t.Attack,
		t.Shirts[0],
		t.Shirts[1],
		t.Shirts[2],
		t.Shirts[3],
		t.Shirts[4],
		t.Shirts[5],
		t.Shirts[6],
		t.Shirts[7],
		t.Shirts[8],
		t.Shirts[9],
		t.Shirts[10],
		t.ExtraAttack[0],
		t.ExtraAttack[1],
		t.ExtraAttack[2],
		t.ExtraAttack[3],
		t.ExtraAttack[4],
		t.ExtraAttack[5],
		t.ExtraAttack[6],
		t.ExtraAttack[7],
		t.ExtraAttack[8],
		t.ExtraAttack[9],
	)
	return err
}
