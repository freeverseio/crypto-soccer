package storage

import (
	"database/sql"
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Player struct {
	PlayerId          *big.Int
	PreferredPosition string
	Potential         uint64
	DayOfBirth        uint64
	State             PlayerState
}

type PlayerState struct {
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
		b.State.TeamId.String() == player.State.TeamId.String() &&
		b.State.Defence == player.State.Defence &&
		b.State.Speed == player.State.Speed &&
		b.State.Pass == player.State.Pass &&
		b.State.Shoot == player.State.Shoot &&
		b.State.Endurance == player.State.Endurance &&
		b.State.ShirtNumber == player.State.ShirtNumber &&
		b.State.EncodedSkills.String() == player.State.EncodedSkills.String() &&
		b.State.EncodedState.String() == player.State.EncodedState.String() &&
		b.State.Frozen == player.State.Frozen &&
		b.State.RedCardMatchesLeft == player.State.RedCardMatchesLeft &&
		b.State.InjuryMatchesLeft == player.State.InjuryMatchesLeft &&
		b.State.Name == player.State.Name &&
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

func (b *Player) PlayerCreate(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create player %v", b)
	_, err := tx.Exec("INSERT INTO players (player_id, team_id, defence, speed, pass, shoot, endurance, shirt_number, preferred_position, encoded_skills, encoded_state, potential, frozen, name, day_of_birth) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);",
		b.PlayerId.String(),
		b.State.TeamId.String(),
		b.State.Defence,
		b.State.Speed,
		b.State.Pass,
		b.State.Shoot,
		b.State.Endurance,
		b.State.ShirtNumber,
		b.PreferredPosition,
		b.State.EncodedSkills.String(),
		b.State.EncodedState.String(),
		b.Potential,
		b.State.Frozen,
		b.State.Name,
		b.DayOfBirth,
	)
	if err != nil {
		return err
	}

	return nil
}

func (b *Player) PlayerUpdate(tx *sql.Tx, playerID *big.Int, playerState PlayerState) error {
	log.Debugf("[DBMS] + update player state %v", playerState)
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
		playerState.TeamId.String(),
		playerState.Defence,
		playerState.Speed,
		playerState.Pass,
		playerState.Shoot,
		playerState.Endurance,
		playerState.ShirtNumber,
		playerState.Frozen,
		playerState.EncodedSkills.String(),
		playerState.RedCardMatchesLeft,
		playerState.InjuryMatchesLeft,
		playerState.Name,
		playerID.String(),
	)
	return err
}

func GetPlayer(tx *sql.Tx, playerID *big.Int) (Player, error) {
	player := Player{}
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
		return player, err
	}
	defer rows.Close()
	if !rows.Next() {
		return player, errors.New("Unexistent player " + playerID.String())
	}
	var teamID sql.NullString
	var encodedSkills sql.NullString
	var encodedState sql.NullString
	err = rows.Scan(
		&teamID,
		&player.State.Defence,
		&player.State.Speed,
		&player.State.Pass,
		&player.State.Shoot,
		&player.State.Endurance,
		&player.State.ShirtNumber,
		&player.PreferredPosition,
		&encodedSkills,
		&encodedState,
		&player.Potential,
		&player.State.Frozen,
		&player.State.Name,
		&player.DayOfBirth,
		&player.State.RedCardMatchesLeft,
		&player.State.InjuryMatchesLeft,
	)
	player.PlayerId = playerID
	player.State.TeamId, _ = new(big.Int).SetString(teamID.String, 10)
	player.State.EncodedSkills, _ = new(big.Int).SetString(encodedSkills.String, 10)
	player.State.EncodedState, _ = new(big.Int).SetString(encodedState.String, 10)
	return player, err
}

func GetPlayersOfTeam(tx *sql.Tx, teamID *big.Int) ([]Player, error) {
	var players []Player
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
		player, err := GetPlayer(tx, playerID)
		if err != nil {
			return players, err
		}
		players = append(players, player)
	}
	return players, err
}
