package universe

import (
	"crypto/sha256"
	"database/sql"
	"errors"
)

type Atom struct {
	encodedSkills string
	encodedState  string
}

type Universe struct {
	atoms []Atom
}

func NewFromStorage(tx *sql.Tx, timezone int) (*Universe, error) {
	var err error
	rows, err := tx.Query(`SELECT 
		encoded_skills, encoded_state
		FROM players 
		LEFT JOIN teams 
		ON players.team_id = teams.team_id 
		WHERE teams.timezone_idx = $1  
		ORDER BY player_id ASC;`, timezone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	universe := Universe{}
	for rows.Next() {
		var atom Atom
		if err = rows.Scan(
			&atom.encodedSkills,
			&atom.encodedState,
		); err != nil {
			return nil, err
		}
		universe.atoms = append(universe.atoms, atom)
	}

	return &universe, nil
}

func (b Universe) Hash() ([32]byte, error) {
	var result [32]byte
	h := sha256.New()
	for _, atom := range b.atoms {
		h.Write([]byte(atom.encodedSkills))
		h.Write([]byte(atom.encodedState))
	}
	hash := h.Sum(nil)
	if len(hash) != 32 {
		return result, errors.New("Hash is not 32 byte")
	}
	copy(result[:], hash)
	return result, nil
}

func (b Universe) Size() int {
	return len(b.atoms)
}
