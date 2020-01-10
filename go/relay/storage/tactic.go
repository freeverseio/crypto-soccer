package storage

import (
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"
)

type Tactic struct {
	Verse         uint64 `json:"verse"`
	TeamID        string `json:"team_id"`         // team_id
	TacticID      int    `json:"tactic_id"`       // tactic_id
	Shirt0        int    `json:"shirt_0"`         // shirt_0
	Shirt1        int    `json:"shirt_1"`         // shirt_1
	Shirt2        int    `json:"shirt_2"`         // shirt_2
	Shirt3        int    `json:"shirt_3"`         // shirt_3
	Shirt4        int    `json:"shirt_4"`         // shirt_4
	Shirt5        int    `json:"shirt_5"`         // shirt_5
	Shirt6        int    `json:"shirt_6"`         // shirt_6
	Shirt7        int    `json:"shirt_7"`         // shirt_7
	Shirt8        int    `json:"shirt_8"`         // shirt_8
	Shirt9        int    `json:"shirt_9"`         // shirt_9
	Shirt10       int    `json:"shirt_10"`        // shirt_10
	Shirt11       int    `json:"shirt_11"`        // shirt_11
	Shirt12       int    `json:"shirt_12"`        // shirt_12
	Shirt13       int    `json:"shirt_13"`        // shirt_13
	ExtraAttack1  bool   `json:"extra_attack_1"`  // extra_attack_1
	ExtraAttack2  bool   `json:"extra_attack_2"`  // extra_attack_2
	ExtraAttack3  bool   `json:"extra_attack_3"`  // extra_attack_3
	ExtraAttack4  bool   `json:"extra_attack_4"`  // extra_attack_4
	ExtraAttack5  bool   `json:"extra_attack_5"`  // extra_attack_5
	ExtraAttack6  bool   `json:"extra_attack_6"`  // extra_attack_6
	ExtraAttack7  bool   `json:"extra_attack_7"`  // extra_attack_7
	ExtraAttack8  bool   `json:"extra_attack_8"`  // extra_attack_8
	ExtraAttack9  bool   `json:"extra_attack_9"`  // extra_attack_9
	ExtraAttack10 bool   `json:"extra_attack_10"` // extra_attack_1
}

func DefaultTactic(teamID string) *Tactic {
	tacticId := 1
	return &Tactic{UpcomingVerse, teamID, tacticId, 0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 25, 25, 25, false, false, true, false, false, true, false, false, false, false}
}

func (b *Tactic) Delete(tx *sql.Tx) error {
	log.Debugf("[DBMS] Delete tactic %v", b)
	_, err := tx.Exec(`DELETE FROM tactics WHERE (verse=$1) AND (team_id=$2);`, b.Verse, b.TeamID)
	return err
}

func TacticsByVerse(tx *sql.Tx, verse uint64) ([]Tactic, error) {
	var tactics []Tactic
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
		FROM tactics WHERE (verse = $1);`, verse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Tactic
		err = rows.Scan(
			&t.Verse,
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
		tactics = append(tactics, t)
	}
	return tactics, nil
}

func UpcomingTactics(tx *sql.Tx) ([]Tactic, error) {
	return TacticsByVerse(tx, UpcomingVerse)
}

func TacticByTeamIDAndVerse(tx *sql.Tx, teamID string, verse uint64) (*Tactic, error) {
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
		FROM tactics WHERE (team_id = $1) AND (verse = $2);`, teamID, verse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Unexistent tactic")
	}
	var t Tactic
	err = rows.Scan(
		&t.TeamID,
		&t.Verse,
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

func TacticCountByVerse(tx *sql.Tx, verse uint64) (uint64, error) {
	count := uint64(0)
	var rows *sql.Rows
	var err error

	rows, err = tx.Query("SELECT COUNT(*) FROM tactics WHERE (verse = $1);", verse)

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

func TacticCount(tx *sql.Tx) (uint64, error) {
	count := uint64(0)
	var rows *sql.Rows
	var err error

	rows, err = tx.Query("SELECT COUNT(*) FROM tactics;")

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

func (b *Tactic) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create tactic for TeamID %v", b.TeamID)
	_, err := tx.Exec(
		`INSERT INTO tactics (
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
		b.TeamID,
		b.TacticID,
		b.Shirt0,
		b.Shirt1,
		b.Shirt2,
		b.Shirt3,
		b.Shirt4,
		b.Shirt5,
		b.Shirt6,
		b.Shirt7,
		b.Shirt8,
		b.Shirt9,
		b.Shirt10,
		b.Shirt11,
		b.Shirt12,
		b.Shirt13,
		b.ExtraAttack1,
		b.ExtraAttack2,
		b.ExtraAttack3,
		b.ExtraAttack4,
		b.ExtraAttack5,
		b.ExtraAttack6,
		b.ExtraAttack7,
		b.ExtraAttack8,
		b.ExtraAttack9,
		b.ExtraAttack10,
	)
	return err
}

func TacticByTeamID(tx *sql.Tx, teamID string) (*Tactic, error) {
	log.Debugf("[DBMS] GetTactic by teamID %v", teamID)
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
		FROM tactics WHERE team_id=$1;`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	var t Tactic
	err = rows.Scan(
		&t.Verse,
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
