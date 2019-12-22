package storage

import (
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"
)

type VerseTactic struct {
	Verse  uint64 `json:"verse"`
	Tactic Tactic `json:"verse"`
}

func VerseTacticCount(tx *sql.Tx, verse uint64) (uint64, error) {
	count := uint64(0)
	var rows *sql.Rows
	var err error

	rows, err = tx.Query("SELECT COUNT(*) FROM verse_tactics WHERE (verse = $1);", verse)

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

func (b *VerseTactic) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create tactic for TeamID %v", b.Tactic.TeamID)
	_, err := tx.Exec(
		`INSERT INTO verse_tactics (
						verse,
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
						$26,
						$27
		);`,
		b.Verse,
		b.Tactic.TacticID,
		b.Tactic.TeamID,
		b.Tactic.Shirt0,
		b.Tactic.Shirt1,
		b.Tactic.Shirt2,
		b.Tactic.Shirt3,
		b.Tactic.Shirt4,
		b.Tactic.Shirt5,
		b.Tactic.Shirt6,
		b.Tactic.Shirt7,
		b.Tactic.Shirt8,
		b.Tactic.Shirt9,
		b.Tactic.Shirt10,
		b.Tactic.Shirt11,
		b.Tactic.Shirt12,
		b.Tactic.Shirt13,
		b.Tactic.ExtraAttack1,
		b.Tactic.ExtraAttack2,
		b.Tactic.ExtraAttack3,
		b.Tactic.ExtraAttack4,
		b.Tactic.ExtraAttack5,
		b.Tactic.ExtraAttack6,
		b.Tactic.ExtraAttack7,
		b.Tactic.ExtraAttack8,
		b.Tactic.ExtraAttack9,
		b.Tactic.ExtraAttack10,
	)
	return err
}

func VerseTacticsByVerse(tx *sql.Tx, verse uint64) ([]VerseTactic, error) {
	var tactics []VerseTactic
	rows, err := tx.Query(
		`SELECT
				verse,
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
		FROM verse_tactics WHERE (verse = $1);`, verse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t VerseTactic
		err = rows.Scan(
			&t.Verse,
			&t.Tactic.TeamID,
			&t.Tactic.TacticID,
			&t.Tactic.Shirt0,
			&t.Tactic.Shirt1,
			&t.Tactic.Shirt2,
			&t.Tactic.Shirt3,
			&t.Tactic.Shirt4,
			&t.Tactic.Shirt5,
			&t.Tactic.Shirt6,
			&t.Tactic.Shirt7,
			&t.Tactic.Shirt8,
			&t.Tactic.Shirt9,
			&t.Tactic.Shirt10,
			&t.Tactic.Shirt11,
			&t.Tactic.Shirt12,
			&t.Tactic.Shirt13,
			&t.Tactic.ExtraAttack1,
			&t.Tactic.ExtraAttack2,
			&t.Tactic.ExtraAttack3,
			&t.Tactic.ExtraAttack4,
			&t.Tactic.ExtraAttack5,
			&t.Tactic.ExtraAttack6,
			&t.Tactic.ExtraAttack7,
			&t.Tactic.ExtraAttack8,
			&t.Tactic.ExtraAttack9,
			&t.Tactic.ExtraAttack10,
		)
		if err != nil {
			return nil, err
		}
		tactics = append(tactics, t)
	}
	return tactics, nil
}

func VerseTacticByTeamIDAndVerse(tx *sql.Tx, teamID string, verse uint64) (*VerseTactic, error) {
	log.Debugf("[DBMS] GetTactic by teamID %v an verse %v", teamID, verse)
	rows, err := tx.Query(
		`SELECT
		team_id,
		verse,
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
		FROM verse_tactics WHERE (team_id = $1) AND (verse = $2);`, teamID, verse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Unexistent tactic")
	}
	var t VerseTactic
	err = rows.Scan(
		&t.Tactic.TeamID,
		&t.Verse,
		&t.Tactic.TacticID,
		&t.Tactic.Shirt0,
		&t.Tactic.Shirt1,
		&t.Tactic.Shirt2,
		&t.Tactic.Shirt3,
		&t.Tactic.Shirt4,
		&t.Tactic.Shirt5,
		&t.Tactic.Shirt6,
		&t.Tactic.Shirt7,
		&t.Tactic.Shirt8,
		&t.Tactic.Shirt9,
		&t.Tactic.Shirt10,
		&t.Tactic.Shirt11,
		&t.Tactic.Shirt12,
		&t.Tactic.Shirt13,
		&t.Tactic.ExtraAttack1,
		&t.Tactic.ExtraAttack2,
		&t.Tactic.ExtraAttack3,
		&t.Tactic.ExtraAttack4,
		&t.Tactic.ExtraAttack5,
		&t.Tactic.ExtraAttack6,
		&t.Tactic.ExtraAttack7,
		&t.Tactic.ExtraAttack8,
		&t.Tactic.ExtraAttack9,
		&t.Tactic.ExtraAttack10,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
