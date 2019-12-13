package storage

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

// Training represents a row from 'public.trainings'.
type Training struct {
	CreatedAt              time.Time `json:"created_at"`               // created_at
	TeamID                 string    `json:"team_id"`                  // team_id
	SpecialPlayerShirt     int       `json:"special_player_shirt"`     // special_player_shirt
	GoalkeepersDefence     int       `json:"goalkeepers_defence"`      // goalkeepers_defence
	GoalkeepersSpeed       int       `json:"goalkeepers_speed"`        // goalkeepers_speed
	GoalkeepersPass        int       `json:"goalkeepers_pass"`         // goalkeepers_pass
	GoalkeepersShoot       int       `json:"goalkeepers_shoot"`        // goalkeepers_shoot
	GoalkeepersEndurance   int       `json:"goalkeepers_endurance"`    // goalkeepers_endurance
	DefendersDefence       int       `json:"defenders_defence"`        // defenders_defence
	DefendersSpeed         int       `json:"defenders_speed"`          // defenders_speed
	DefendersPass          int       `json:"defenders_pass"`           // defenders_pass
	DefendersShoot         int       `json:"defenders_shoot"`          // defenders_shoot
	DefendersEndurance     int       `json:"defenders_endurance"`      // defenders_endurance
	MidfieldersDefence     int       `json:"midfielders_defence"`      // midfielders_defence
	MidfieldersSpeed       int       `json:"midfielders_speed"`        // midfielders_speed
	MidfieldersPass        int       `json:"midfielders_pass"`         // midfielders_pass
	MidfieldersShoot       int       `json:"midfielders_shoot"`        // midfielders_shoot
	MidfieldersEndurance   int       `json:"midfielders_endurance"`    // midfielders_endurance
	AttackersDefence       int       `json:"attackers_defence"`        // attackers_defence
	AttackersSpeed         int       `json:"attackers_speed"`          // attackers_speed
	AttackersPass          int       `json:"attackers_pass"`           // attackers_pass
	AttackersShoot         int       `json:"attackers_shoot"`          // attackers_shoot
	AttackersEndurance     int       `json:"attackers_endurance"`      // attackers_endurance
	SpecialPlayerDefence   int       `json:"special_player_defence"`   // special_player_defence
	SpecialPlayerSpeed     int       `json:"special_player_speed"`     // special_player_speed
	SpecialPlayerPass      int       `json:"special_player_pass"`      // special_player_pass
	SpecialPlayerShoot     int       `json:"special_player_shoot"`     // special_player_shoot
	SpecialPlayerEndurance int       `json:"special_player_endurance"` // special_player_endurance
}

func (b *Storage) CreateTraining(training Training) error {
	log.Debugf("[DBMS] Create training %v", training)
	_, err := b.tx.Exec(
		`INSERT INTO trainings (
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
			$27
        );`,
		training.TeamID,
		training.SpecialPlayerShirt,
		training.GoalkeepersDefence,
		training.GoalkeepersSpeed,
		training.GoalkeepersPass,
		training.GoalkeepersShoot,
		training.GoalkeepersEndurance,
		training.DefendersDefence,
		training.DefendersSpeed,
		training.DefendersPass,
		training.DefendersShoot,
		training.DefendersEndurance,
		training.MidfieldersDefence,
		training.MidfieldersSpeed,
		training.MidfieldersPass,
		training.MidfieldersShoot,
		training.MidfieldersEndurance,
		training.AttackersDefence,
		training.AttackersSpeed,
		training.AttackersPass,
		training.AttackersShoot,
		training.AttackersEndurance,
		training.SpecialPlayerDefence,
		training.SpecialPlayerSpeed,
		training.SpecialPlayerPass,
		training.SpecialPlayerShoot,
		training.SpecialPlayerEndurance,
	)
	return err
}

func (b *Storage) GetRowsTrainingsRange(start *Verse, end *Verse) (*sql.Rows, error) {
	log.Debugf("[DBMS] GetRowsTrainingsRange from verse %v to verse  %v", start, end)
	return b.tx.Query(
		`SELECT * FROM trainings WHERE (created_at > $1) AND (created_at <= $2);`,
		start.StartAt,
		end.StartAt,
	)
}

// func (b *Storage) UpdateTraining(training Training) error {
// 	log.Debugf("[DBMS] Create training %v", training)
// 	_, err := b.tx.Exec(
// 		`INSERT INTO trainings (
// 			team_id,
//     		special_player_shirt,
// 			goalkeepers_defence,
//     		goalkeepers_speed,
//     		goalkeepers_pass,
//     		goalkeepers_shoot,
//     		goalkeepers_endurance,
//     		defenders_defence,
//     		defenders_speed,
//     		defenders_pass,
//     		defenders_shoot,
//     		defenders_endurance,
//     		midfielders_defence,
//     		midfielders_speed,
//     		midfielders_pass,
//     		midfielders_shoot,
//     		midfielders_endurance,
//     		attackers_defence,
//     		attackers_speed,
//     		attackers_pass,
//     		attackers_shoot,
//     		attackers_endurance,
//     		special_player_defence,
//     		special_player_speed,
//     		special_player_pass,
//     		special_player_shoot,
//     		special_player_endurance,
// 		);`,
// 		training.TeamID.String(),
// 		training.SpecialPlayerShirt,
// 		training.GoalkeepersDefence,
// 		training.GoalkeepersSpeed,
// 		training.GoalkeepersPass,
// 		training.GoalkeepersShoot,
// 		training.GoalkeepersEndurance,
// 		training.DefendersDefence,
// 		training.DefendersSpeed,
// 		training.DefendersPass,
// 		training.DefendersShoot,
// 		training.DefendersEndurance,
// 		training.MidfieldersDefence,
// 		training.MidfieldersSpeed,
// 		training.MidfieldersPass,
// 		training.MidfieldersShoot,
// 		training.MidfieldersEndurance,
// 		training.AttackersDefence,
// 		training.AttackersSpeed,
// 		training.AttackersPass,
// 		training.AttackersShoot,
// 		training.AttackersEndurance,
// 		training.SpecialPlayerDefence,
// 		training.SpecialPlayerSpeed,
// 		training.SpecialPlayerPass,
// 		training.SpecialPlayerShoot,
// 		training.SpecialPlayerEndurance,
// 	)
// 	return err
// }
