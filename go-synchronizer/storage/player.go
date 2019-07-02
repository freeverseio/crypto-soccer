package storage

import "math/big"

type Player struct {
	Id    uint64
	State *big.Int
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
	_, err := b.db.Exec("INSERT INTO players (id, state) VALUES ($1, $2);", player.Id, player.State.String())
	if err != nil {
		return err
	}

	return nil
}
