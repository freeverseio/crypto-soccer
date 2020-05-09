package universe

import (
	"crypto/sha256"
	"database/sql"
	"errors"
)

type UniversePlayer struct {
	encodedSkills string
	encodedState  string
}

type Universe struct {
	players []UniversePlayer
}

func NewFromStorage(tx *sql.Tx, timezone int) (*Universe, error) {
	var err error
	rows, err := tx.Query(`SELECT 
		encoded_skills, encoded_state
		FROM players 
		LEFT JOIN teams 
		ON players.team_id = teams.team_id 
		WHERE teams.timezone_idx = $1  
		ORDER BY player_id DESC;`, timezone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	universe := Universe{}
	for rows.Next() {
		var player UniversePlayer
		if err = rows.Scan(
			&player.encodedSkills,
			&player.encodedState,
		); err != nil {
			return nil, err
		}
		universe.players = append(universe.players, player)
	}

	return &universe, nil
}

func (b Universe) Hash() ([32]byte, error) {
	var result [32]byte
	h := sha256.New()
	for _, player := range b.players {
		h.Write([]byte(player.encodedSkills))
		h.Write([]byte(player.encodedState))
	}
	hash := h.Sum(nil)
	if len(hash) != 32 {
		return result, errors.New("Hash is not 32 byte")
	}
	copy(result[:], hash)
	return result, nil
}
