package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Training represents a row from 'public.trainings'.
type VerseTraining struct {
	Verse    uint64   `json:"verse"`
	Training Training `json:"training"`
}

func (b *VerseTraining) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create training %v", b)
	_, err := tx.Exec(
		`INSERT INTO verse_trainings (
			verse,
			team_id,
    		special_player_shirt,
			goalkeepers_defence,
    		goalkeepers_speed,
    		goalkeepers_pass,
    		goalkeepers_shoot,
    		goalkeepers_endurance,
    		defenders_defence,
    		defenders_speed,
    		defenders_pass,
    		defenders_shoot,
    		defenders_endurance,
    		midfielders_defence,
    		midfielders_speed,
    		midfielders_pass,
    		midfielders_shoot,
    		midfielders_endurance,
    		attackers_defence,
    		attackers_speed,
    		attackers_pass,
    		attackers_shoot,
    		attackers_endurance,
    		special_player_defence,
    		special_player_speed,
    		special_player_pass,
    		special_player_shoot,
			special_player_endurance
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
			$28
		);`,
		b.Verse,
		b.Training.TeamID,
		b.Training.SpecialPlayerShirt,
		b.Training.GoalkeepersDefence,
		b.Training.GoalkeepersSpeed,
		b.Training.GoalkeepersPass,
		b.Training.GoalkeepersShoot,
		b.Training.GoalkeepersEndurance,
		b.Training.DefendersDefence,
		b.Training.DefendersSpeed,
		b.Training.DefendersPass,
		b.Training.DefendersShoot,
		b.Training.DefendersEndurance,
		b.Training.MidfieldersDefence,
		b.Training.MidfieldersSpeed,
		b.Training.MidfieldersPass,
		b.Training.MidfieldersShoot,
		b.Training.MidfieldersEndurance,
		b.Training.AttackersDefence,
		b.Training.AttackersSpeed,
		b.Training.AttackersPass,
		b.Training.AttackersShoot,
		b.Training.AttackersEndurance,
		b.Training.SpecialPlayerDefence,
		b.Training.SpecialPlayerSpeed,
		b.Training.SpecialPlayerPass,
		b.Training.SpecialPlayerShoot,
		b.Training.SpecialPlayerEndurance,
	)
	return err
}

func VerseTrainingByVerse(tx *sql.Tx, verse uint64) ([]VerseTraining, error) {
	var trainings []VerseTraining
	rows, err := tx.Query(
		`SELECT
			verse,
			team_id,
    		special_player_shirt,
			goalkeepers_defence,
    		goalkeepers_speed,
    		goalkeepers_pass,
    		goalkeepers_shoot,
    		goalkeepers_endurance,
    		defenders_defence,
    		defenders_speed,
    		defenders_pass,
    		defenders_shoot,
    		defenders_endurance,
    		midfielders_defence,
    		midfielders_speed,
    		midfielders_pass,
    		midfielders_shoot,
    		midfielders_endurance,
    		attackers_defence,
    		attackers_speed,
    		attackers_pass,
    		attackers_shoot,
    		attackers_endurance,
    		special_player_defence,
    		special_player_speed,
    		special_player_pass,
    		special_player_shoot,
			special_player_endurance
		FROM verse_trainings WHERE (verse = $1);`, verse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t VerseTraining
		err = rows.Scan(
			&t.Verse,
			&t.Training.TeamID,
			&t.Training.SpecialPlayerShirt,
			&t.Training.GoalkeepersDefence,
			&t.Training.GoalkeepersSpeed,
			&t.Training.GoalkeepersPass,
			&t.Training.GoalkeepersShoot,
			&t.Training.GoalkeepersEndurance,
			&t.Training.DefendersDefence,
			&t.Training.DefendersSpeed,
			&t.Training.DefendersPass,
			&t.Training.DefendersShoot,
			&t.Training.DefendersEndurance,
			&t.Training.MidfieldersDefence,
			&t.Training.MidfieldersSpeed,
			&t.Training.MidfieldersPass,
			&t.Training.MidfieldersShoot,
			&t.Training.MidfieldersEndurance,
			&t.Training.AttackersDefence,
			&t.Training.AttackersSpeed,
			&t.Training.AttackersPass,
			&t.Training.AttackersShoot,
			&t.Training.AttackersEndurance,
			&t.Training.SpecialPlayerDefence,
			&t.Training.SpecialPlayerSpeed,
			&t.Training.SpecialPlayerPass,
			&t.Training.SpecialPlayerShoot,
			&t.Training.SpecialPlayerEndurance,
		)
		if err != nil {
			return nil, err
		}
		trainings = append(trainings, t)
	}
	return trainings, nil
}
