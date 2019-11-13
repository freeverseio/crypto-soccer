package storage

import (
	"database/sql"
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Player struct {
	PlayerId          *big.Int
	Name              string
	PreferredPosition string
	Potential         uint64
	DayOfBirth        uint64
	State             PlayerState
}

type PlayerState struct {
	TeamId        *big.Int
	Defence       uint64
	Speed         uint64
	Pass          uint64
	Shoot         uint64
	Endurance     uint64
	ShirtNumber   uint8
	EncodedSkills *big.Int
	EncodedState  *big.Int
	Frozen        bool
	RedCard       bool
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
		b.Name == player.Name &&
		b.DayOfBirth == player.DayOfBirth
}

func (b *Storage) PlayerCount() (uint64, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM players;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	rows.Scan(&count)
	return count, nil
}

func (b *Storage) PlayerCreate(player Player) error {
	log.Debugf("[DBMS] Create player %v", player)
	_, err := b.db.Exec("INSERT INTO players (player_id, team_id, defence, speed, pass, shoot, endurance, shirt_number, preferred_position, encoded_skills, encoded_state, potential, frozen, name, day_of_birth) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);",
		player.PlayerId.String(),
		player.State.TeamId.String(),
		player.State.Defence,
		player.State.Speed,
		player.State.Pass,
		player.State.Shoot,
		player.State.Endurance,
		player.State.ShirtNumber,
		player.PreferredPosition,
		player.State.EncodedSkills.String(),
		player.State.EncodedState.String(),
		player.Potential,
		player.State.Frozen,
		player.Name,
		player.DayOfBirth,
	)
	if err != nil {
		return err
	}

	// err = b.playerHistoryAdd(player.Id, player.State)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (b *Storage) PlayerUpdate(playerID *big.Int, playerState PlayerState) error {
	log.Debugf("[DBMS] + update player state %v", playerState)
	_, err := b.db.Exec(`UPDATE players SET 
	team_id=$1, 
	defence=$2, 
	speed=$3, 
	pass=$4, 
	shoot=$5,
	endurance=$6,
	shirt_number=$7,
	frozen=$8, 
	encoded_skills=$9,
	red_card=$10
	WHERE player_id=$11;`,
		playerState.TeamId.String(),
		playerState.Defence,
		playerState.Speed,
		playerState.Pass,
		playerState.Shoot,
		playerState.Endurance,
		playerState.ShirtNumber,
		playerState.Frozen,
		playerState.EncodedSkills.String(),
		playerState.RedCard,
		playerID.String(),
	)
	return err
}

func (b *Storage) GetPlayer(playerID *big.Int) (Player, error) {
	player := Player{}
	rows, err := b.db.Query("SELECT team_id, defence, speed, pass, shoot, endurance, shirt_number, preferred_position, encoded_skills, encoded_state, potential, frozen, name, day_of_birth, red_card FROM players WHERE (player_id = $1);", playerID.String())
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
		&player.Name,
		&player.DayOfBirth,
		&player.State.RedCard,
	)
	player.PlayerId = playerID
	player.State.TeamId, _ = new(big.Int).SetString(teamID.String, 10)
	player.State.EncodedSkills, _ = new(big.Int).SetString(encodedSkills.String, 10)
	player.State.EncodedState, _ = new(big.Int).SetString(encodedState.String, 10)
	return player, err
}

func (b *Storage) GetPlayersOfTeam(teamID *big.Int) ([]Player, error) {
	var players []Player
	rows, err := b.db.Query("SELECT player_id FROM players WHERE (team_id = $1);", teamID.String())
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
		player, err := b.GetPlayer(playerID)
		if err != nil {
			return players, err
		}
		players = append(players, player)
	}
	return players, err
}

// func (b *Storage) playerUpdate(id uint64, playerState PlayerState) error {
// 	log.Infof("[DBMS] + update player state %v", playerState)

// 	_, err := b.db.Exec("UPDATE players SET blockNumber=$1, teamId=$2, state=$3, defence=$4, speed=$5, pass=$6, shoot=$7, endurance=$8, inBlockIndex=$9 WHERE id=$10;",
// 		playerState.BlockNumber,
// 		playerState.TeamId,
// 		playerState.State,
// 		playerState.Defence,
// 		playerState.Speed,
// 		playerState.Pass,
// 		playerState.Shoot,
// 		playerState.Endurance,
// 		playerState.InBlockIndex,
// 		id,
// 	)
// 	return err
// }

// func (b *Storage) playerHistoryAdd(id uint64, playerState PlayerState) error {
// 	log.Infof("[DBMS] + add player history %v", playerState)
// 	_, err := b.db.Exec("INSERT INTO players_history (playerId, blockNumber, teamId, state, defence, speed, pass, shoot, endurance, inBlockIndex) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
// 		id,
// 		playerState.BlockNumber,
// 		playerState.TeamId,
// 		playerState.State,
// 		playerState.Defence,
// 		playerState.Speed,
// 		playerState.Pass,
// 		playerState.Shoot,
// 		playerState.Endurance,
// 		playerState.InBlockIndex,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
