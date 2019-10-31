package storage

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"hash"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Tactic struct {
	TeamID      *big.Int
	Defense     uint8
	Center      uint8
	Attack      uint8
	Shirts      [11]uint8
	ExtraAttack [10]bool
}

// Hash - computes hash for a Tactic
func (t *Tactic) Hash() ([32]byte, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return [32]byte{}, err
	}
	b := [32]byte{}
	h := computeHash(sha256.New(), data)
	copy(b[:], h[:])
	return b, nil
}

func computeHash(h hash.Hash, data ...[]byte) []byte {
	h.Reset()
	for _, d := range data {
		h.Write(d)
	}
	return h.Sum(nil)
}

func (b *Storage) TacticCreate(t Tactic, verse uint64) error {
	log.Debugf("[DBMS] Create tactic %v", t)
	_, err := b.db.Exec(
		`INSERT INTO tactics (
			team_id,
			verse,
                        defense,
                        center,
                        attack,
                        shirt_0,
                        shirt_1,
                        shirt_2,
                        shirt_3,
                        shirt_4,
                        shirt_5,
                        shirt_6,
                        shirt_7,
                        shirt_8,
                        shirt_9,
                        shirt_10,
                        extra_attack_1,
                        extra_attack_2,
                        extra_attack_3,
                        extra_attack_4,
                        extra_attack_5,
                        extra_attack_6,
                        extra_attack_7,
                        extra_attack_8,
                        extra_attack_9,
                        extra_attack_10
		) VALUES (
			$1,
                        $2,
                        $3,
                        $4,
                        $5,
                        $6,
                        $7,
                        $8,
                        $9,
                        $10,
                        $11,
                        $12,
                        $13,
                        $14,
                        $15,
                        $16,
                        $17,
                        $18,
                        $19,
                        $20,
                        $21,
                        $22,
                        $23,
                        $24,
                        $25,
                        $26
		);`,
		t.TeamID.String(),
		verse,
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
func (b *Storage) GetTactic(teamID *big.Int, verse uint64) (*Tactic, error) {
	log.Debugf("[DBMS] GetTactic of teamID %v", teamID)
	rows, err := b.db.Query(
		`SELECT
		defense,
                center,
                attack,
                shirt_0,
                shirt_1,
                shirt_2,
                shirt_3,
                shirt_4,
                shirt_5,
                shirt_6,
                shirt_7,
                shirt_8,
                shirt_9,
                shirt_10,
                extra_attack_1,
                extra_attack_2,
                extra_attack_3,
                extra_attack_4,
                extra_attack_5,
                extra_attack_6,
                extra_attack_7,
                extra_attack_8,
                extra_attack_9,
                extra_attack_10
		FROM tactics WHERE (team_id = $1) and (verse = $2);`, teamID.String(), verse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Unexistent tactic")
	}
	var t Tactic
	t.TeamID = teamID
	err = rows.Scan(
		&t.Defense,
		&t.Center,
		&t.Attack,
		&t.Shirts[0],
		&t.Shirts[1],
		&t.Shirts[2],
		&t.Shirts[3],
		&t.Shirts[4],
		&t.Shirts[5],
		&t.Shirts[6],
		&t.Shirts[7],
		&t.Shirts[8],
		&t.Shirts[9],
		&t.Shirts[10],
		&t.ExtraAttack[0],
		&t.ExtraAttack[1],
		&t.ExtraAttack[2],
		&t.ExtraAttack[3],
		&t.ExtraAttack[4],
		&t.ExtraAttack[5],
		&t.ExtraAttack[6],
		&t.ExtraAttack[7],
		&t.ExtraAttack[8],
		&t.ExtraAttack[9],
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func (b *Storage) TacticCount(verse *uint64) (uint64, error) {
	count := uint64(0)
	var rows *sql.Rows
	var err error

	if verse == nil {
		rows, err = b.db.Query("SELECT COUNT(*) FROM tactics;")
	} else {
		rows, err = b.db.Query("SELECT COUNT(*) FROM tactics WHERE (verse = $1);", *verse)
	}

	if err != nil {
		return 0, err
	}

	defer rows.Close()
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
