package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type TacticHistory struct {
	Tactic
	BlockNumber uint64
}

func NewTacticHistory(tactic Tactic) *TacticHistory {
	tacticHistory := TacticHistory{}
	tacticHistory.Tactic = tactic
	return &tacticHistory
}

func TacticHistoryByTeamID(tx *sql.Tx, teamID string) ([]TacticHistory, error) {
	var tacticHistories []TacticHistory
	rows, err := tx.Query(
		`SELECT
				block_number,
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
				substitution_0_shirt,
				substitution_0_target,
				substitution_0_minute,
                substitution_1_shirt,
				substitution_1_target,
				substitution_1_minute,
                substitution_2_shirt,
				substitution_2_target,
				substitution_2_minute,
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
		FROM tactics_histories WHERE team_id = $1;`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t TacticHistory
		err = rows.Scan(
			&t.BlockNumber,
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
			&t.Substitution0Shirt,
			&t.Substitution0Target,
			&t.Substitution0Minute,
			&t.Substitution1Shirt,
			&t.Substitution1Target,
			&t.Substitution1Minute,
			&t.Substitution2Shirt,
			&t.Substitution2Target,
			&t.Substitution2Minute,
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
		tacticHistories = append(tacticHistories, t)
	}
	return tacticHistories, nil
}

func (b TacticHistory) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create tactic history for TeamID %v", b.TeamID)
	_, err := tx.Exec(
		`INSERT INTO tactics_histories (
						block_number,
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
                        substitution_0_shirt,
						substitution_0_target,
						substitution_0_minute,
                		substitution_1_shirt,
						substitution_1_target,
						substitution_1_minute,
                		substitution_2_shirt,
						substitution_2_target,
						substitution_2_minute,
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
						$27,
						$28,
						$29,
						$30,
						$31,
						$32,
						$33
		);`,
		b.BlockNumber,
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
		b.Substitution0Shirt,
		b.Substitution0Target,
		b.Substitution0Minute,
		b.Substitution1Shirt,
		b.Substitution1Target,
		b.Substitution1Minute,
		b.Substitution2Shirt,
		b.Substitution2Target,
		b.Substitution2Minute,
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
