package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// order: shoot, speed, pass, defence, endurance
type TrainingPerFieldPos struct {
	Shoot     int `json:"Shoot"`
	Speed     int `json:"Speed"`
	Pass      int `json:"Pass"`
	Defence   int `json:"Defence"`
	Endurance int `json:"Endurance"`
}

func (b TrainingPerFieldPos) ToSlice() []uint16 {
	return []uint16{
		uint16(b.Shoot),
		uint16(b.Speed),
		uint16(b.Pass),
		uint16(b.Defence),
		uint16(b.Endurance),
	}
}

// Training represents a row from 'public.trainings'.
type Training struct {
	TeamID             string              `json:"team_id"`              // team_id
	SpecialPlayerShirt int                 `json:"special_player_shirt"` // special_player_shirt
	Goalkeepers        TrainingPerFieldPos `json:"Goalkeepers"`          // goalkeepers_training
	Defenders          TrainingPerFieldPos `json:"Defenders"`            // defenders_training
	Midfielders        TrainingPerFieldPos `json:"Midfielders"`          // midfielders_training
	Attackers          TrainingPerFieldPos `json:"Attackers"`            // attackers_training
	SpecialPlayer      TrainingPerFieldPos `json:"SpecialPlayer"`        // specialPlayer_training
}

func NewTraining() *Training {
	training := Training{}
	training.SpecialPlayerShirt = -1
	return &training
}

// order: shoot, speed, pass, defence, endurance
func CreateDefaultTrainingByTimezone(tx *sql.Tx, timezone uint8) error {
	log.Debugf("[DBMS] Create a default training for each Team in timezone %v", timezone)
	if _, err := tx.Exec(`
		INSERT INTO trainings (
			team_id,
			serialized_training
		) SELECT team_id,"" FROM teams WHERE teams.timezone_idx = $1`, timezone); err != nil {
		return err
	}
	return nil
}

func ResetTrainingsByTimezone(tx *sql.Tx, timezone uint8) error {
	log.Debugf("[DBMS] Create a default training for each Team in timezone %v", timezone)
	if _, err := tx.Exec(`
		UPDATE trainings SET (
			serialized_training
		) = ("") FROM teams WHERE teams.team_id = trainings.team_id AND teams.timezone_idx = $1`, timezone); err != nil {
		return err
	}
	return nil
}

func (b *Training) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create training %v", b)
	_, err := tx.Exec(
		`INSERT INTO trainings (
			team_id,
			serialized_training
		) VALUES (                    
			$1,
			$2,
		);`,
		b.TeamID,
		b.SerializedTraining,
	)
	return err
}

func TrainingsByTimezone(tx *sql.Tx, timezone int) ([]Training, error) {
	var trainings []Training
	rows, err := tx.Query(
		`SELECT 
			trainings.team_id,
			trainings.serialized_training
		FROM trainings LEFT JOIN teams ON trainings.team_id = teams.team_id WHERE teams.timezone_idx = $1;`, timezone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Training
		err = rows.Scan(
			&t.TeamID,
			&t.SerializedTraining,
		)
		if err != nil {
			return nil, err
		}
		trainings = append(trainings, t)
	}
	return trainings, nil
}
