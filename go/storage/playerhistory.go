package storage

import (
	"database/sql"
	"fmt"
	"strings"
)

type PlayerHistory struct {
	Player
	BlockNumber uint64
}

func NewPlayerHistory(blockNumber uint64, player Player) *PlayerHistory {
	history := PlayerHistory{}
	history.Player = player
	history.BlockNumber = blockNumber
	return &history
}

func (b PlayerHistory) Insert(tx *sql.Tx) error {
	if _, err := tx.Exec(`INSERT INTO players_histories 
		(block_number, player_id, team_id, defence, speed, pass, shoot, endurance, 
		shirt_number, preferred_position, encoded_skills, 
		encoded_state, potential, day_of_birth, tiredness, country_of_birth, race) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17);`,
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
		b.CountryOfBirth,
		b.Race,
	); err != nil {
		return err
	}
	return nil
}

func PlayersHistoriesBulkInsert(rowsToBeInserted []*PlayerHistory, tx *sql.Tx) error {
	numParams := 17
	var err error = nil
	maxRowsToBeInserted := int(MAX_PARAMS_IN_PG_STMT / numParams)
	x := 0
	for x < len(rowsToBeInserted) {
		newX := x + maxRowsToBeInserted
		if newX > len(rowsToBeInserted) {
			newX = len(rowsToBeInserted)
		}
		currentRowsToBeInserted := rowsToBeInserted[x:newX]
		valueStrings := make([]string, 0, len(currentRowsToBeInserted))
		valueArgs := make([]interface{}, 0, len(currentRowsToBeInserted)*numParams)
		i := 0
		for _, post := range currentRowsToBeInserted {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*numParams+1, i*numParams+2, i*numParams+3, i*numParams+4, i*numParams+5, i*numParams+6, i*numParams+7, i*numParams+8, i*numParams+9, i*numParams+10, i*numParams+11, i*numParams+12, i*numParams+13, i*numParams+14, i*numParams+15, i*numParams+16, i*numParams+17))
			valueArgs = append(valueArgs, post.BlockNumber)
			valueArgs = append(valueArgs, post.PlayerId.String())
			valueArgs = append(valueArgs, post.TeamId)
			valueArgs = append(valueArgs, post.Defence)
			valueArgs = append(valueArgs, post.Speed)
			valueArgs = append(valueArgs, post.Pass)
			valueArgs = append(valueArgs, post.Shoot)
			valueArgs = append(valueArgs, post.Endurance)
			valueArgs = append(valueArgs, post.ShirtNumber)
			valueArgs = append(valueArgs, post.PreferredPosition)
			valueArgs = append(valueArgs, post.EncodedSkills.String())
			valueArgs = append(valueArgs, post.EncodedState.String())
			valueArgs = append(valueArgs, post.Potential)
			valueArgs = append(valueArgs, post.DayOfBirth)
			valueArgs = append(valueArgs, post.Tiredness)
			valueArgs = append(valueArgs, post.CountryOfBirth)
			valueArgs = append(valueArgs, post.Race)
			i++
		}
		stmt := fmt.Sprintf(`INSERT INTO players_histories (
			block_number,
			player_id,
			team_id,
			defence,
			speed,
			pass,
			shoot,
			endurance, 
			shirt_number,
			preferred_position,
			encoded_skills, 
			encoded_state,
			potential,
			day_of_birth,
			tiredness,
			country_of_birth,
			race
			) VALUES %s
			`, strings.Join(valueStrings, ","))
		_, err := tx.Exec(stmt, valueArgs...)
		if err != nil {
			return err
		}
		x = newX
	}
	return err
}
