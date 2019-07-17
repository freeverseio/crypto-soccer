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
	Id        uint64
	TeamId    uint64
	State     string
	Defence   uint64
	Speed     uint64
	Pass      uint64
	Shoot     uint64
	Endurance uint64
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
	log.Infof("(DBMS) Adding player state %v", player)
	_, err := b.db.Exec("INSERT INTO players (id, monthOfBirthInUnixTime) VALUES ($1, $2);",
		player.Id,
		player.MonthOfBirthInUnixTime)
	if err != nil {
		return err
	}

	err = b.PlayerStateAdd(player.State)

	return nil
}

func (b *Storage) PlayerStateAdd(playerState PlayerState) error {
	log.Infof("(DBMS) Adding player state %v", playerState)
	_, err := b.db.Exec("INSERT INTO players_history (playerId, teamId, state, defence, speed, pass, shoot, endurance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
		playerState.Id,
		playerState.TeamId,
		playerState.State,
		playerState.Defence,
		playerState.Speed,
		playerState.Shoot,
		playerState.Pass,
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
	return player, nil
}
