package storage

import "database/sql"

type PlayerHistory struct {
	Player
	BlockNumber uint64
}

func (b PlayerHistory) InsertHistory(tx *sql.Tx) error {
	if _, err := tx.Exec(`INSERT INTO players_histories 
		(block_number, player_id, team_id, defence, speed, pass, shoot, endurance, 
		shirt_number, preferred_position, encoded_skills, 
		encoded_state, potential, day_of_birth, tiredness) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`,
		b.BlockNumber,
		b.PlayerId.String(),
		b.TeamId,
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
		b.DayOfBirth,
		b.Tiredness,
	); err != nil {
		return err
	}
	return nil
}
