package storage

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

type Tactic struct {
	CreatedAt     time.Time `json:"created_at"`      // created_at
	TeamID        string    `json:"team_id"`         // team_id
	TacticID      int       `json:"tactic_id"`       // tactic_id
	Shirt0        int       `json:"shirt_0"`         // shirt_0
	Shirt1        int       `json:"shirt_1"`         // shirt_1
	Shirt2        int       `json:"shirt_2"`         // shirt_2
	Shirt3        int       `json:"shirt_3"`         // shirt_3
	Shirt4        int       `json:"shirt_4"`         // shirt_4
	Shirt5        int       `json:"shirt_5"`         // shirt_5
	Shirt6        int       `json:"shirt_6"`         // shirt_6
	Shirt7        int       `json:"shirt_7"`         // shirt_7
	Shirt8        int       `json:"shirt_8"`         // shirt_8
	Shirt9        int       `json:"shirt_9"`         // shirt_9
	Shirt10       int       `json:"shirt_10"`        // shirt_10
	Shirt11       int       `json:"shirt_11"`        // shirt_11
	Shirt12       int       `json:"shirt_12"`        // shirt_12
	Shirt13       int       `json:"shirt_13"`        // shirt_13
	ExtraAttack1  bool      `json:"extra_attack_1"`  // extra_attack_1
	ExtraAttack2  bool      `json:"extra_attack_2"`  // extra_attack_2
	ExtraAttack3  bool      `json:"extra_attack_3"`  // extra_attack_3
	ExtraAttack4  bool      `json:"extra_attack_4"`  // extra_attack_4
	ExtraAttack5  bool      `json:"extra_attack_5"`  // extra_attack_5
	ExtraAttack6  bool      `json:"extra_attack_6"`  // extra_attack_6
	ExtraAttack7  bool      `json:"extra_attack_7"`  // extra_attack_7
	ExtraAttack8  bool      `json:"extra_attack_8"`  // extra_attack_8
	ExtraAttack9  bool      `json:"extra_attack_9"`  // extra_attack_9
	ExtraAttack10 bool      `json:"extra_attack_10"` // extra_attack_1
}

func (b *Storage) DefaultTactic(teamID string) *Tactic {
	tacticId := uint8(1)
	return &Tactic{time.Now(), teamID, int(tacticId), 0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 25, 25, 25, false, false, true, false, false, true, false, false, false, false}
}
func (b *Storage) TacticCreate(tactic *Tactic) error {
	log.Debugf("[DBMS] Create tactic for TeamID %v", tactic.TeamID)
	_, err := b.tx.Exec(
		`INSERT INTO tactics (
						tactic_id,
						team_id,
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
                        shirt_11,
                        shirt_12,
                        shirt_13,
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
		tactic.TacticID,
		tactic.TeamID,
		tactic.Shirt0,
		tactic.Shirt1,
		tactic.Shirt2,
		tactic.Shirt3,
		tactic.Shirt4,
		tactic.Shirt5,
		tactic.Shirt6,
		tactic.Shirt7,
		tactic.Shirt8,
		tactic.Shirt9,
		tactic.Shirt10,
		tactic.Shirt11,
		tactic.Shirt12,
		tactic.Shirt13,
		tactic.ExtraAttack1,
		tactic.ExtraAttack2,
		tactic.ExtraAttack3,
		tactic.ExtraAttack4,
		tactic.ExtraAttack5,
		tactic.ExtraAttack6,
		tactic.ExtraAttack7,
		tactic.ExtraAttack8,
		tactic.ExtraAttack9,
		tactic.ExtraAttack10,
	)
	return err
}

func (b *Storage) TacticsByVerse(verseNumber int) ([]*Tactic, error) {
	var tactics []*Tactic
	if verseNumber == 0 {
		return tactics, nil
	}
	verse, err := b.GetVerse(verseNumber)
	if err != nil {
		return nil, err
	}
	prevVerse, err := b.GetVerse(verseNumber - 1)
	if err != nil {
		return nil, err
	}
	rows, err := b.tx.Query(
		`SELECT
				created_at,
				team_id,
				tactic_id,
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
                shirt_11,
                shirt_12,
                shirt_13,
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
		FROM tactics WHERE (created_at > $1) AND (created_at <= $2);`, prevVerse.StartAt, verse.StartAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Tactic
		err = rows.Scan(
			&t.CreatedAt,
			&t.TeamID,
			&t.TacticID,
			&t.Shirt0,
			&t.Shirt1,
			&t.Shirt2,
			&t.Shirt3,
			&t.Shirt4,
			&t.Shirt5,
			&t.Shirt6,
			&t.Shirt7,
			&t.Shirt8,
			&t.Shirt9,
			&t.Shirt10,
			&t.Shirt11,
			&t.Shirt12,
			&t.Shirt13,
			&t.ExtraAttack1,
			&t.ExtraAttack2,
			&t.ExtraAttack3,
			&t.ExtraAttack4,
			&t.ExtraAttack5,
			&t.ExtraAttack6,
			&t.ExtraAttack7,
			&t.ExtraAttack8,
			&t.ExtraAttack9,
			&t.ExtraAttack10,
		)
		if err != nil {
			return nil, err
		}
		tactics = append(tactics, &t)
	}
	return tactics, nil
}

func (b *Storage) GetTactic(teamID string, verse uint64) (*Tactic, error) {
	log.Debugf("[DBMS] GetTactic of teamID %v", teamID)
	rows, err := b.tx.Query(
		`SELECT
		created_at,
		team_id,
		tactic_id,
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
                shirt_11,
                shirt_12,
                shirt_13,
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
		FROM tactics WHERE (team_id = $1);`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Unexistent tactic")
	}
	var t Tactic
	err = rows.Scan(
		&t.CreatedAt,
		&t.TeamID,
		&t.TacticID,
		&t.Shirt0,
		&t.Shirt1,
		&t.Shirt2,
		&t.Shirt3,
		&t.Shirt4,
		&t.Shirt5,
		&t.Shirt6,
		&t.Shirt7,
		&t.Shirt8,
		&t.Shirt9,
		&t.Shirt10,
		&t.Shirt11,
		&t.Shirt12,
		&t.Shirt13,
		&t.ExtraAttack1,
		&t.ExtraAttack2,
		&t.ExtraAttack3,
		&t.ExtraAttack4,
		&t.ExtraAttack5,
		&t.ExtraAttack6,
		&t.ExtraAttack7,
		&t.ExtraAttack8,
		&t.ExtraAttack9,
		&t.ExtraAttack10,
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

	rows, err = b.tx.Query("SELECT COUNT(*) FROM tactics;")

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
