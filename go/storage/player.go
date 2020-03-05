package storage

import (
	"database/sql"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Player struct {
	PlayerId          *big.Int
	PreferredPosition string
	Potential         uint64
	DayOfBirth        uint64
	TeamId            string
	Name              string
	Defence           uint64
	Speed             uint64
	Pass              uint64
	Shoot             uint64
	Endurance         uint64
	ShirtNumber       uint8
	EncodedSkills     *big.Int
	EncodedState      *big.Int
	RedCard           bool
	InjuryMatchesLeft uint8
	BlockNumber       uint64
}

func (b *Player) Equal(player Player) bool {
	return b.PlayerId.String() == player.PlayerId.String() &&
		b.PreferredPosition == player.PreferredPosition &&
		b.Potential == player.Potential &&
		b.TeamId == player.TeamId &&
		b.Defence == player.Defence &&
		b.Speed == player.Speed &&
		b.Pass == player.Pass &&
		b.Shoot == player.Shoot &&
		b.Endurance == player.Endurance &&
		b.ShirtNumber == player.ShirtNumber &&
		b.EncodedSkills.String() == player.EncodedSkills.String() &&
		b.EncodedState.String() == player.EncodedState.String() &&
		b.RedCard == player.RedCard &&
		b.InjuryMatchesLeft == player.InjuryMatchesLeft &&
		b.Name == player.Name &&
		b.DayOfBirth == player.DayOfBirth
}

func PlayerCount(tx *sql.Tx) (uint64, error) {
	rows, err := tx.Query("SELECT COUNT(*) FROM players;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	rows.Scan(&count)
	return count, nil
}

func (b *Player) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create player %v", b)
	if _, err := tx.Exec("INSERT INTO players (name, block_number, player_id, team_id, defence, speed, pass, shoot, endurance, shirt_number, preferred_position, encoded_skills, encoded_state, potential, day_of_birth) VALUES ($1, $2,$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);",
		b.Name,
		b.BlockNumber,
		b.PlayerId.String(),
		b.TeamId,
		b.Defence,
		b.Speed,
		b.Pass,
		b.Shoot,
		b.Endurance,
		b.ShirtNumber,
		b.PreferredPosition,
		b.EncodedSkills.String(),
		b.EncodedState.String(),
		b.Potential,
		b.DayOfBirth,
	); err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO players_states (block_number, player_id, team_id, defence, speed, pass, shoot, endurance, shirt_number, preferred_position, encoded_skills, encoded_state, potential, day_of_birth) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);",
		b.BlockNumber,
		b.PlayerId.String(),
		b.TeamId,
		b.Defence,
		b.Speed,
		b.Pass,
		b.Shoot,
		b.Endurance,
		b.ShirtNumber,
		b.PreferredPosition,
		b.EncodedSkills.String(),
		b.EncodedState.String(),
		b.Potential,
		b.DayOfBirth,
	); err != nil {
		return err
	}
	return nil
}

func (b *Player) Update(tx *sql.Tx) error {
	log.Debugf("[DBMS] + update player id %v", b.PlayerId)
	if _, err := tx.Exec(`UPDATE players SET 
	team_id=$1, 
	defence=$2, 
	speed=$3, 
	pass=$4, 
	shoot=$5,
	endurance=$6,
	shirt_number=$7,
	encoded_skills=$8,
	red_card=$9,
	injury_matches_left=$10,
	name=$11,
	block_number=$12
	WHERE player_id=$13;`,
		b.TeamId,
		b.Defence,
		b.Speed,
		b.Pass,
		b.Shoot,
		b.Endurance,
		b.ShirtNumber,
		b.EncodedSkills.String(),
		b.RedCard,
		b.InjuryMatchesLeft,
		b.Name,
		b.BlockNumber,
		b.PlayerId.String(),
	); err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT INTO players_states (block_number, player_id, team_id, defence, speed, pass, shoot, endurance, shirt_number, preferred_position, encoded_skills, encoded_state, potential, day_of_birth, red_card,
	injury_matches_left) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`,
		b.BlockNumber,
		b.PlayerId.String(),
		b.TeamId,
		b.Defence,
		b.Speed,
		b.Pass,
		b.Shoot,
		b.Endurance,
		b.ShirtNumber,
		b.PreferredPosition,
		b.EncodedSkills.String(),
		b.EncodedState.String(),
		b.Potential,
		b.DayOfBirth,
		b.RedCard,
		b.InjuryMatchesLeft,
	); err != nil {
		return err
	}
	return nil
}

func PlayerByPlayerId(tx *sql.Tx, playerID *big.Int) (*Player, error) {
	rows, err := tx.Query(`SELECT 
	block_number,
	team_id, 
	defence,
	speed,
	pass, 
	shoot, 
	endurance, 
	shirt_number, 
	preferred_position, 
	encoded_skills, 
	encoded_state, 
	potential, 
	name, 
	day_of_birth, 
	red_card,
	injury_matches_left
	FROM players WHERE (player_id = $1);`, playerID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	player := Player{}
	var encodedSkills sql.NullString
	var encodedState sql.NullString
	err = rows.Scan(
		&player.BlockNumber,
		&player.TeamId,
		&player.Defence,
		&player.Speed,
		&player.Pass,
		&player.Shoot,
		&player.Endurance,
		&player.ShirtNumber,
		&player.PreferredPosition,
		&encodedSkills,
		&encodedState,
		&player.Potential,
		&player.Name,
		&player.DayOfBirth,
		&player.RedCard,
		&player.InjuryMatchesLeft,
	)
	player.PlayerId = playerID
	player.EncodedSkills, _ = new(big.Int).SetString(encodedSkills.String, 10)
	player.EncodedState, _ = new(big.Int).SetString(encodedState.String, 10)
	return &player, nil
}

func PlayersByTeamId(tx *sql.Tx, teamID string) ([]*Player, error) {
	rows, err := tx.Query(`SELECT 
	block_number,
	player_id, 
	defence,
	speed,
	pass, 
	shoot, 
	endurance, 
	shirt_number, 
	preferred_position, 
	encoded_skills, 
	encoded_state, 
	potential, 
	name, 
	day_of_birth, 
	red_card,
	injury_matches_left
	FROM players WHERE (team_id = $1);`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*Player
	for rows.Next() {
		player := Player{}
		var encodedSkills sql.NullString
		var encodedState sql.NullString
		var playerID sql.NullString
		err = rows.Scan(
			&player.BlockNumber,
			&playerID,
			&player.Defence,
			&player.Speed,
			&player.Pass,
			&player.Shoot,
			&player.Endurance,
			&player.ShirtNumber,
			&player.PreferredPosition,
			&encodedSkills,
			&encodedState,
			&player.Potential,
			&player.Name,
			&player.DayOfBirth,
			&player.RedCard,
			&player.InjuryMatchesLeft,
		)
		player.TeamId = teamID
		player.EncodedSkills, _ = new(big.Int).SetString(encodedSkills.String, 10)
		player.EncodedState, _ = new(big.Int).SetString(encodedState.String, 10)
		player.PlayerId, _ = new(big.Int).SetString(playerID.String, 10)
		players = append(players, &player)
	}
	return players, err
}
