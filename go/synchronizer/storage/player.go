package storage

import (
	"database/sql"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Player struct {
	PlayerId           *big.Int
	PreferredPosition  string
	Potential          uint64
	DayOfBirth         uint64
	TeamId             *big.Int
	Name               string
	Defence            uint64
	Speed              uint64
	Pass               uint64
	Shoot              uint64
	Endurance          uint64
	ShirtNumber        uint8
	EncodedSkills      *big.Int
	EncodedState       *big.Int
	Frozen             bool
	RedCardMatchesLeft uint8
	InjuryMatchesLeft  uint8
}

func (b *Player) Equal(player Player) bool {
	return b.PlayerId.String() == player.PlayerId.String() &&
		b.PreferredPosition == player.PreferredPosition &&
		b.Potential == player.Potential &&
		b.TeamId.String() == player.TeamId.String() &&
		b.Defence == player.Defence &&
		b.Speed == player.Speed &&
		b.Pass == player.Pass &&
		b.Shoot == player.Shoot &&
		b.Endurance == player.Endurance &&
		b.ShirtNumber == player.ShirtNumber &&
		b.EncodedSkills.String() == player.EncodedSkills.String() &&
		b.EncodedState.String() == player.EncodedState.String() &&
		b.Frozen == player.Frozen &&
		b.RedCardMatchesLeft == player.RedCardMatchesLeft &&
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
	_, err := tx.Exec("INSERT INTO players (player_id, team_id, defence, speed, pass, shoot, endurance, shirt_number, preferred_position, encoded_skills, encoded_state, potential, frozen, name, day_of_birth) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);",
		b.PlayerId.String(),
		b.TeamId.String(),
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
		b.Frozen,
		b.Name,
		b.DayOfBirth,
	)
	if err != nil {
		return err
	}

	return nil
}

func (b *Player) Update(tx *sql.Tx) error {
	log.Debugf("[DBMS] + update player id %v", b.PlayerId)
	_, err := tx.Exec(`UPDATE players SET 
	team_id=$1, 
	defence=$2, 
	speed=$3, 
	pass=$4, 
	shoot=$5,
	endurance=$6,
	shirt_number=$7,
	frozen=$8, 
	encoded_skills=$9,
	red_card_matches_left=$10,
	injury_matches_left=$11,
	name=$12
	WHERE player_id=$13;`,
		b.TeamId.String(),
		b.Defence,
		b.Speed,
		b.Pass,
		b.Shoot,
		b.Endurance,
		b.ShirtNumber,
		b.Frozen,
		b.EncodedSkills.String(),
		b.RedCardMatchesLeft,
		b.InjuryMatchesLeft,
		b.Name,
		b.PlayerId.String(),
	)
	return err
}

func PlayerByPlayerId(tx *sql.Tx, playerID *big.Int) (*Player, error) {
	rows, err := tx.Query(`SELECT team_id, 
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
	frozen, 
	name, 
	day_of_birth, 
	red_card_matches_left,
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
	var teamID sql.NullString
	var encodedSkills sql.NullString
	var encodedState sql.NullString
	err = rows.Scan(
		&teamID,
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
		&player.Frozen,
		&player.Name,
		&player.DayOfBirth,
		&player.RedCardMatchesLeft,
		&player.InjuryMatchesLeft,
	)
	player.PlayerId = playerID
	player.TeamId, _ = new(big.Int).SetString(teamID.String, 10)
	player.EncodedSkills, _ = new(big.Int).SetString(encodedSkills.String, 10)
	player.EncodedState, _ = new(big.Int).SetString(encodedState.String, 10)
	return &player, nil
}

func PlayersByTeamId(tx *sql.Tx, teamID *big.Int) ([]*Player, error) {
	var players []*Player
	rows, err := tx.Query("SELECT player_id FROM players WHERE (team_id = $1);", teamID.String())
	if err != nil {
		return players, err
	}
	defer rows.Close()
	var playerIDs []*big.Int
	for rows.Next() {
		var playerID sql.NullString
		err = rows.Scan(
			&playerID,
		)
		result, _ := new(big.Int).SetString(playerID.String, 10)
		playerIDs = append(playerIDs, result)
	}
	rows.Close()
	for i := 0; i < len(playerIDs); i++ {
		playerID := playerIDs[i]
		player, err := PlayerByPlayerId(tx, playerID)
		if err != nil {
			return players, err
		}
		players = append(players, player)
	}
	return players, err
}
