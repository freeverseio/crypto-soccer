package storage

import (
	log "github.com/sirupsen/logrus"
)

type Player struct {
	Id                     uint64
	MonthOfBirthInUnixTime string
	State                  PlayerState
}

type PlayerState struct {
	TeamId      uint64
	BlockNumber uint64
	State       string
	Defence     uint64
	Speed       uint64
	Pass        uint64
	Shoot       uint64
	Endurance   uint64
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

func (b *Storage) PlayerAdd(player Player) error {
	log.Infof("(DBMS) Adding player %v", player)
	_, err := b.db.Exec("INSERT INTO players (id, monthOfBirthInUnixTime, blockNumber, teamId, state, defence, speed, pass, shoot, endurance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
		player.Id,
		player.MonthOfBirthInUnixTime,
		player.State.BlockNumber,
		player.State.TeamId,
		player.State.State,
		player.State.Defence,
		player.State.Speed,
		player.State.Pass,
		player.State.Shoot,
		player.State.Endurance,
	)
	if err != nil {
		return err
	}

	err = b.PlayerStateAdd(player.Id, player.State)
	if err != nil {
		return err
	}

	return nil
}

func (b *Storage) PlayerStateAdd(id uint64, playerState PlayerState) error {
	log.Infof("(DBMS) Adding player state %v", playerState)
	_, err := b.db.Exec("INSERT INTO players_history (playerId, blockNumber, teamId, state, defence, speed, pass, shoot, endurance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);",
		id,
		playerState.BlockNumber,
		playerState.TeamId,
		playerState.State,
		playerState.Defence,
		playerState.Speed,
		playerState.Pass,
		playerState.Shoot,
		playerState.Endurance)
	if err != nil {
		return err
	}

	return nil
}

func (b *Storage) GetPlayer(id uint64) (Player, error) {
	player := Player{}
	rows, err := b.db.Query("SELECT id, monthOfBirthInUnixTime FROM players WHERE (id = $1);", id)
	if err != nil {
		return player, err
	}
	defer rows.Close()
	if !rows.Next() {
		return player, nil
	}
	rows.Scan(&player.Id, &player.MonthOfBirthInUnixTime)
	rows.Close()
	player.State, err = b.GetPlayerState(id)
	if err != nil {
		return player, err
	}
	return player, nil
}

func (b *Storage) GetPlayerState(id uint64) (PlayerState, error) {
	playerState := PlayerState{}
	rows, err := b.db.Query("SELECT blockNumber, teamId, state, defence, speed, pass, shoot, endurance FROM players_history WHERE (playerId = $1) ORDER BY blockNumber DESC LIMIT 1 ;", id)
	if err != nil {
		return playerState, err
	}
	defer rows.Close()
	if !rows.Next() {
		return playerState, nil
	}
	rows.Scan(&playerState.BlockNumber, &playerState.TeamId, &playerState.State, &playerState.Defence, &playerState.Speed, &playerState.Pass, &playerState.Shoot, &playerState.Endurance)

	return playerState, nil
}
