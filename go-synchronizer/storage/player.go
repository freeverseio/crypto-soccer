package storage

import (
	log "github.com/sirupsen/logrus"
)

type Player struct {
	Id        uint64
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

func (b *Storage) PlayerAdd(player *Player) error {
	log.Infof("(DBMS) Adding player %v %v", player.Id, player.State)
	_, err := b.db.Exec("INSERT INTO players (id, state, defence, speed, pass, shoot, endurance) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		player.Id,
		player.State,
		player.Defence,
		player.Speed,
		player.Shoot,
		player.Pass,
		player.Endurance)
	if err != nil {
		return err
	}

	return nil
}

func (b *Storage) GetPlayer(id uint64) (*Player, error) {
	player := Player{}
	rows, err := b.db.Query("SELECT id, state FROM players WHERE (id = $1);", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	rows.Scan(&player.Id, &player.State)
	return &player, nil
}
