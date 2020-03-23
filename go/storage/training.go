package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// order: shoot, speed, pass, defence, endurance
type TrainingPerFieldPos struct {
	Shoot     int `json:"shoot"`
	Speed     int `json:"speed"`
	Pass      int `json:"pass"`
	Defence   int `json:"defence"`
	Endurance int `json:"endurance"`
}

// Training represents a row from 'public.trainings'.
type Training struct {
	TeamID             string              `json:"team_id"`                // team_id
	SpecialPlayerShirt int                 `json:"special_player_shirt"`   // special_player_shirt
	Goalkeepers        TrainingPerFieldPos `json:"goalkeepers_training"`   // goalkeepers_training
	Defenders          TrainingPerFieldPos `json:"defenders_training"`     // defenders_training
	Midfielders        TrainingPerFieldPos `json:"midfielders_training"`   // midfielders_training
	Attackers          TrainingPerFieldPos `json:"attackers_training"`     // attackers_training
	SpecialPlayer      TrainingPerFieldPos `json:"specialPlayer_training"` // specialPlayer_training
}

func NewTraining() *Training {
	training := Training{}
	training.SpecialPlayerShirt = -1
	return &training
}

func DeleteTrainingsByTimezone(tx *sql.Tx, timezone uint8) error {
	log.Debugf("[DBMS] Delete trainings by Timezone %v", timezone)
	if _, err := tx.Exec(`DELETE FROM trainings USING teams WHERE trainings.team_id = teams.team_id AND teams.timezone_idx = $1`, timezone); err != nil {
		return err
	}
	return nil
}

// order: shoot, speed, pass, defence, endurance
func CreateDefaultTrainingByTimezone(tx *sql.Tx, timezone uint8) error {
	log.Debugf("[DBMS] Create a default training for each Team in timezone %v", timezone)
	if _, err := tx.Exec(`
		INSERT INTO trainings (
			team_id,
    		special_player_shirt,
    		goalkeepers_shoot,
    		goalkeepers_speed,
    		goalkeepers_pass,
			goalkeepers_defence,
    		goalkeepers_endurance,
    		defenders_shoot,
    		defenders_speed,
    		defenders_pass,
    		defenders_defence,
    		defenders_endurance,
    		midfielders_shoot,
    		midfielders_speed,
    		midfielders_pass,
    		midfielders_defence,
    		midfielders_endurance,
    		attackers_shoot,
    		attackers_speed,
    		attackers_pass,
    		attackers_defence,
    		attackers_endurance,
    		special_player_shoot,
    		special_player_speed,
    		special_player_pass,
    		special_player_defence,
			special_player_endurance
		) SELECT team_id,-1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0 FROM teams WHERE teams.timezone_idx = $1`, timezone); err != nil {
		return err
	}
	return nil
}

func (b *Training) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create training %v", b)
	_, err := tx.Exec(
		`INSERT INTO trainings (
			team_id,
    		special_player_shirt,
    		goalkeepers_shoot,
    		goalkeepers_speed,
    		goalkeepers_pass,
			goalkeepers_defence,
    		goalkeepers_endurance,
    		defenders_shoot,
    		defenders_speed,
    		defenders_pass,
    		defenders_defence,
    		defenders_endurance,
    		midfielders_shoot,
    		midfielders_speed,
    		midfielders_pass,
    		midfielders_defence,
    		midfielders_endurance,
    		attackers_shoot,
    		attackers_speed,
    		attackers_pass,
    		attackers_defence,
    		attackers_endurance,
    		special_player_shoot,
    		special_player_speed,
    		special_player_pass,
    		special_player_defence,
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
			$27
		);`,
		b.TeamID,
		b.SpecialPlayerShirt,
		b.Goalkeepers.Shoot,
		b.Goalkeepers.Speed,
		b.Goalkeepers.Pass,
		b.Goalkeepers.Defence,
		b.Goalkeepers.Endurance,
		b.Defenders.Shoot,
		b.Defenders.Speed,
		b.Defenders.Pass,
		b.Defenders.Defence,
		b.Defenders.Endurance,
		b.Midfielders.Shoot,
		b.Midfielders.Speed,
		b.Midfielders.Pass,
		b.Midfielders.Defence,
		b.Midfielders.Endurance,
		b.Attackers.Shoot,
		b.Attackers.Speed,
		b.Attackers.Pass,
		b.Attackers.Defence,
		b.Attackers.Endurance,
		b.SpecialPlayer.Shoot,
		b.SpecialPlayer.Speed,
		b.SpecialPlayer.Pass,
		b.SpecialPlayer.Defence,
		b.SpecialPlayer.Endurance,
	)
	return err
}

func TrainingsByTimezone(tx *sql.Tx, timezone int) ([]Training, error) {
	var trainings []Training
	rows, err := tx.Query(
		`SELECT 
			trainings.team_id,
    		special_player_shirt,
    		goalkeepers_shoot,
    		goalkeepers_speed,
    		goalkeepers_pass,
			goalkeepers_defence,
    		goalkeepers_endurance,
    		defenders_shoot,
    		defenders_speed,
    		defenders_pass,
    		defenders_defence,
    		defenders_endurance,
    		midfielders_shoot,
    		midfielders_speed,
    		midfielders_pass,
    		midfielders_defence,
    		midfielders_endurance,
    		attackers_shoot,
    		attackers_speed,
    		attackers_pass,
    		attackers_defence,
    		attackers_endurance,
    		special_player_shoot,
    		special_player_speed,
    		special_player_pass,
    		special_player_defence,
			special_player_endurance
		FROM trainings LEFT JOIN teams ON trainings.team_id = teams.team_id WHERE teams.timezone_idx = $1;`, timezone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Training
		err = rows.Scan(
			&t.TeamID,
			&t.SpecialPlayerShirt,
			&t.Goalkeepers.Shoot,
			&t.Goalkeepers.Speed,
			&t.Goalkeepers.Pass,
			&t.Goalkeepers.Defence,
			&t.Goalkeepers.Endurance,
			&t.Defenders.Shoot,
			&t.Defenders.Speed,
			&t.Defenders.Pass,
			&t.Defenders.Defence,
			&t.Defenders.Endurance,
			&t.Midfielders.Shoot,
			&t.Midfielders.Speed,
			&t.Midfielders.Pass,
			&t.Midfielders.Defence,
			&t.Midfielders.Endurance,
			&t.Attackers.Shoot,
			&t.Attackers.Speed,
			&t.Attackers.Pass,
			&t.Attackers.Defence,
			&t.Attackers.Endurance,
			&t.SpecialPlayer.Shoot,
			&t.SpecialPlayer.Speed,
			&t.SpecialPlayer.Pass,
			&t.SpecialPlayer.Defence,
			&t.SpecialPlayer.Endurance,
		)
		if err != nil {
			return nil, err
		}
		trainings = append(trainings, t)
	}
	return trainings, nil
}
