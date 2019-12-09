package storage

import (
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Training struct {
	TeamID                 *big.Int
	SpecialPlayerShirt     int
	GoalkeepersDefence     int
	GoalkeepersSpeed       int
	GoalkeepersPass        int
	GoalkeepersShoot       int
	GoalkeepersEndurance   int
	DefendersDefence       int
	DefendersSpeed         int
	DefendersPass          int
	DefendersShoot         int
	DefendersEndurance     int
	MidfieldersDefence     int
	MidfieldersSpeed       int
	MidfieldersPass        int
	MidfieldersShoot       int
	MidfieldersEndurance   int
	AttackersDefence       int
	AttackersSpeed         int
	AttackersPass          int
	AttackersShoot         int
	AttackersEndurance     int
	SpecialPlayerDefence   int
	SpecialPlayerSpeed     int
	SpecialPlayerPass      int
	SpecialPlayerShoot     int
	SpecialPlayerEndurance int
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
		training.TeamID.String(),
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
