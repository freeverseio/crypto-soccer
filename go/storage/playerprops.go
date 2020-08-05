package storage

import (
	"database/sql"
)

type PlayerProps struct {
	Player
}

func NewPlayerProps(player Player) *PlayerProps {
	props := PlayerProps{}
	props.Player = player
	return &props
}

func (b PlayerProps) Insert(tx *sql.Tx) error {
	if _, err := tx.Exec(`INSERT INTO player_props 
		(player_id, player_name) 
		VALUES ($1, $2);`,
		b.PlayerId.String(),
		b.Name,
	); err != nil {
		return err
	}
	return nil
}

func (b PlayerProps) Update(tx *sql.Tx) error {
	if _, err := tx.Exec(`UPDATE player_props SET
		player_name = $1
		WHERE player_id = $2`,
		b.Name,
		b.PlayerId.String(),
	); err != nil {
		return err
	}
	return nil
}
