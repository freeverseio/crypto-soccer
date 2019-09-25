package storage

import (
	"math/big"

	log "github.com/sirupsen/logrus"
)

type Player struct {
	PlayerId          *big.Int
	PreferredPosition string
	State             PlayerState
}

type PlayerState struct {
	TeamId      *big.Int
	Defence     uint64
	Speed       uint64
	Pass        uint64
	Shoot       uint64
	Endurance   uint64
	ShirtNumber uint8
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
	_, err := b.db.Exec("INSERT INTO players (player_id, team_id, defence, speed, pass, shoot, endurance, shirt_number, preferred_position) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);",
		player.PlayerId.String(),
		player.State.TeamId.String(),
		player.State.Defence,
		player.State.Speed,
		player.State.Pass,
		player.State.Shoot,
		player.State.Endurance,
		player.State.ShirtNumber,
		player.PreferredPosition,
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

// func (b *Storage) PlayerStateUpdate(id uint64, playerState PlayerState) error {
// 	log.Infof("[DBMS] Adding player state %v", playerState)

// 	err := b.playerUpdate(id, playerState)
// 	if err != nil {
// 		return err
// 	}
// 	err = b.playerHistoryAdd(id, playerState)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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

// func (b *Storage) GetPlayer(id uint64) (Player, error) {
// 	player := Player{}
// 	rows, err := b.db.Query("SELECT id, monthOfBirthInUnixTime, blockNumber, teamId, state, defence, speed, pass, shoot, endurance, inBlockIndex FROM players WHERE (id = $1);", id)
// 	if err != nil {
// 		return player, err
// 	}
// 	defer rows.Close()
// 	if !rows.Next() {
// 		return player, errors.New("Unexistent player " + strconv.FormatUint(id, 10))
// 	}
// 	err = rows.Scan(&player.Id, &player.MonthOfBirthInUnixTime, &player.State.BlockNumber, &player.State.TeamId, &player.State.State, &player.State.Defence, &player.State.Speed, &player.State.Pass, &player.State.Shoot, &player.State.Endurance, &player.State.InBlockIndex)
// 	return player, err
// }
