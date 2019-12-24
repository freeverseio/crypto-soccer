package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Training represents a row from 'public.trainings'.
type Training struct {
	Verse                  uint64 `json:"verse"`
	TeamID                 string `json:"team_id"`                  // team_id
	SpecialPlayerShirt     int    `json:"special_player_shirt"`     // special_player_shirt
	GoalkeepersDefence     int    `json:"goalkeepers_defence"`      // goalkeepers_defence
	GoalkeepersSpeed       int    `json:"goalkeepers_speed"`        // goalkeepers_speed
	GoalkeepersPass        int    `json:"goalkeepers_pass"`         // goalkeepers_pass
	GoalkeepersShoot       int    `json:"goalkeepers_shoot"`        // goalkeepers_shoot
	GoalkeepersEndurance   int    `json:"goalkeepers_endurance"`    // goalkeepers_endurance
	DefendersDefence       int    `json:"defenders_defence"`        // defenders_defence
	DefendersSpeed         int    `json:"defenders_speed"`          // defenders_speed
	DefendersPass          int    `json:"defenders_pass"`           // defenders_pass
	DefendersShoot         int    `json:"defenders_shoot"`          // defenders_shoot
	DefendersEndurance     int    `json:"defenders_endurance"`      // defenders_endurance
	MidfieldersDefence     int    `json:"midfielders_defence"`      // midfielders_defence
	MidfieldersSpeed       int    `json:"midfielders_speed"`        // midfielders_speed
	MidfieldersPass        int    `json:"midfielders_pass"`         // midfielders_pass
	MidfieldersShoot       int    `json:"midfielders_shoot"`        // midfielders_shoot
	MidfieldersEndurance   int    `json:"midfielders_endurance"`    // midfielders_endurance
	AttackersDefence       int    `json:"attackers_defence"`        // attackers_defence
	AttackersSpeed         int    `json:"attackers_speed"`          // attackers_speed
	AttackersPass          int    `json:"attackers_pass"`           // attackers_pass
	AttackersShoot         int    `json:"attackers_shoot"`          // attackers_shoot
	AttackersEndurance     int    `json:"attackers_endurance"`      // attackers_endurance
	SpecialPlayerDefence   int    `json:"special_player_defence"`   // special_player_defence
	SpecialPlayerSpeed     int    `json:"special_player_speed"`     // special_player_speed
	SpecialPlayerPass      int    `json:"special_player_pass"`      // special_player_pass
	SpecialPlayerShoot     int    `json:"special_player_shoot"`     // special_player_shoot
	SpecialPlayerEndurance int    `json:"special_player_endurance"` // special_player_endurance
}

func CurrentTrainings(tx *sql.Tx) ([]Training, error) {
	return TrainingByVerse(tx, UpcomingVerse)
}

func (b *Training) Delete(tx *sql.Tx) error {
	log.Debugf("[DBMS] Delete training %v", b)
	_, err := tx.Exec(`DELETE FROM trainings WHERE (verse=$1) AND (team_id=$2);`, b.Verse, b.TeamID)
	return err

}

func (b *Training) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create training %v", b)
	_, err := tx.Exec(
		`INSERT INTO trainings (
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
		b.TeamID,
		b.SpecialPlayerShirt,
		b.GoalkeepersDefence,
		b.GoalkeepersSpeed,
		b.GoalkeepersPass,
		b.GoalkeepersShoot,
		b.GoalkeepersEndurance,
		b.DefendersDefence,
		b.DefendersSpeed,
		b.DefendersPass,
		b.DefendersShoot,
		b.DefendersEndurance,
		b.MidfieldersDefence,
		b.MidfieldersSpeed,
		b.MidfieldersPass,
		b.MidfieldersShoot,
		b.MidfieldersEndurance,
		b.AttackersDefence,
		b.AttackersSpeed,
		b.AttackersPass,
		b.AttackersShoot,
		b.AttackersEndurance,
		b.SpecialPlayerDefence,
		b.SpecialPlayerSpeed,
		b.SpecialPlayerPass,
		b.SpecialPlayerShoot,
		b.SpecialPlayerEndurance,
	)
	return err
}

func TrainingByVerse(tx *sql.Tx, verse uint64) ([]Training, error) {
	var trainings []Training
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
		FROM trainings WHERE (verse = $1);`, verse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Training
		err = rows.Scan(
			&t.Verse,
			&t.TeamID,
			&t.SpecialPlayerShirt,
			&t.GoalkeepersDefence,
			&t.GoalkeepersSpeed,
			&t.GoalkeepersPass,
			&t.GoalkeepersShoot,
			&t.GoalkeepersEndurance,
			&t.DefendersDefence,
			&t.DefendersSpeed,
			&t.DefendersPass,
			&t.DefendersShoot,
			&t.DefendersEndurance,
			&t.MidfieldersDefence,
			&t.MidfieldersSpeed,
			&t.MidfieldersPass,
			&t.MidfieldersShoot,
			&t.MidfieldersEndurance,
			&t.AttackersDefence,
			&t.AttackersSpeed,
			&t.AttackersPass,
			&t.AttackersShoot,
			&t.AttackersEndurance,
			&t.SpecialPlayerDefence,
			&t.SpecialPlayerSpeed,
			&t.SpecialPlayerPass,
			&t.SpecialPlayerShoot,
			&t.SpecialPlayerEndurance,
		)
		if err != nil {
			return nil, err
		}
		trainings = append(trainings, t)
	}
	return trainings, nil
}
