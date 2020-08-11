package storage

import (
	"database/sql"
)

type TeamProps struct {
	Team
}

func NewTeamProps(team *Team) *TeamProps {
	props := TeamProps{}
	props.Team = *team
	return &props
}

func (b TeamProps) Insert(tx *sql.Tx) error {
	if _, err := tx.Exec(`INSERT INTO team_props 
		(player_id, player_name) 
		VALUES ($1, $2);`,
		b.TeamID,
		b.Name,
	); err != nil {
		return err
	}
	return nil
}

func (b TeamProps) Update(tx *sql.Tx) error {
	if _, err := tx.Exec(`UPDATE team_props SET
		player_name = $1
		WHERE player_id = $2`,
		b.Name,
		b.TeamID,
	); err != nil {
		return err
	}
	return nil
}
