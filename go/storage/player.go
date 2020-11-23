package storage

import (
	"database/sql"
	"fmt"
	"math/big"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Player struct {
	PlayerId          *big.Int
	PreferredPosition string
	Potential         uint64
	DayOfBirth        uint64
	TeamId            string
	Name              string
	Defence           uint64
	Speed             uint64
	Pass              uint64
	Shoot             uint64
	Endurance         uint64
	ShirtNumber       uint8
	EncodedSkills     *big.Int
	EncodedState      *big.Int
	RedCard           bool
	InjuryMatchesLeft uint8
	Tiredness         int
	CountryOfBirth    string
	Race              string
	YellowCard1stHalf bool
}

func (b *Player) Equal(player Player) bool {
	return b.PlayerId.String() == player.PlayerId.String() &&
		b.PreferredPosition == player.PreferredPosition &&
		b.Potential == player.Potential &&
		b.TeamId == player.TeamId &&
		b.Defence == player.Defence &&
		b.Speed == player.Speed &&
		b.Pass == player.Pass &&
		b.Shoot == player.Shoot &&
		b.Endurance == player.Endurance &&
		b.ShirtNumber == player.ShirtNumber &&
		b.EncodedSkills.String() == player.EncodedSkills.String() &&
		b.EncodedState.String() == player.EncodedState.String() &&
		b.RedCard == player.RedCard &&
		b.InjuryMatchesLeft == player.InjuryMatchesLeft &&
		b.Name == player.Name &&
		b.DayOfBirth == player.DayOfBirth &&
		b.CountryOfBirth == player.CountryOfBirth &&
		b.Race == player.Race &&
		b.YellowCard1stHalf == player.YellowCard1stHalf
}

func PlayerCount(tx *sql.Tx) (uint64, error) {
	rows, err := tx.Query("SELECT COUNT(*) FROM players;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	rows.Scan(&count)
	return count, nil
}

func (b Player) Insert(tx *sql.Tx, blockNumber uint64) error {
	log.Debugf("[DBMS] Create player %v", b)
	if _, err := tx.Exec(`INSERT INTO players 
		(name, player_id, team_id, defence, speed,
		pass, shoot, endurance, shirt_number, preferred_position, 
		encoded_skills, encoded_state, potential, day_of_birth, tiredness, country_of_birth, race) 
		VALUES ($1, $2,$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17);`,
		b.Name,
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
	history := NewPlayerHistory(blockNumber, b)
	if err := history.Insert(tx); err != nil {
		return err
	}
	return nil
}

func (b Player) Update(tx *sql.Tx, blockNumber uint64) error {
	log.Debugf("[DBMS] + update player id %v", b.PlayerId)
	if _, err := tx.Exec(`UPDATE players SET 
	team_id=$1, 
	defence=$2, 
	speed=$3, 
	pass=$4, 
	shoot=$5,
	endurance=$6,
	shirt_number=$7,
	encoded_skills=$8,
	red_card=$9,
	injury_matches_left=$10,
	name=$11,
	tiredness=$12,
	country_of_birth=$13,
	race=$14,
	yellow_card_1st_half=$15
	WHERE player_id=$16;`,
		b.TeamId,
		b.Defence,
		b.Speed,
		b.Pass,
		b.Shoot,
		b.Endurance,
		b.ShirtNumber,
		b.EncodedSkills.String(),
		b.RedCard,
		b.InjuryMatchesLeft,
		b.Name,
		b.Tiredness,
		b.CountryOfBirth,
		b.Race,
		b.YellowCard1stHalf,
		b.PlayerId.String(),
	); err != nil {
		return err
	}
	history := NewPlayerHistory(blockNumber, b)
	if err := history.Insert(tx); err != nil {
		return err
	}
	return nil
}

func PlayerByPlayerId(tx *sql.Tx, playerID *big.Int) (*Player, error) {
	rows, err := tx.Query(`SELECT 
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
	name, 
	day_of_birth, 
	red_card,
	injury_matches_left,
	tiredness,
	country_of_birth,
	race,
    yellow_card_1st_half
	FROM players WHERE (player_id = $1);`, playerID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	player := Player{}
	var encodedSkills sql.NullString
	var encodedState sql.NullString
	err = rows.Scan(
		&player.TeamId,
		&player.Defence,
		&player.Speed,
		&player.Pass,
		&player.Shoot,
		&player.Endurance,
		&player.ShirtNumber,
		&player.PreferredPosition,
		&encodedSkills,
		&encodedState,
		&player.Potential,
		&player.Name,
		&player.DayOfBirth,
		&player.RedCard,
		&player.InjuryMatchesLeft,
		&player.Tiredness,
		&player.CountryOfBirth,
		&player.Race,
		&player.YellowCard1stHalf,
	)
	player.PlayerId = playerID
	player.EncodedSkills, _ = new(big.Int).SetString(encodedSkills.String, 10)
	player.EncodedState, _ = new(big.Int).SetString(encodedState.String, 10)
	return &player, nil
}

func ActivePlayersByTeamId(tx *sql.Tx, teamID string) ([]*Player, error) {
	rows, err := tx.Query(`SELECT 
	player_id, 
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
	name, 
	day_of_birth, 
	red_card,
	injury_matches_left,
	tiredness,
	country_of_birth,
	race,
	yellow_card_1st_half
	FROM players WHERE (team_id = $1 AND shirt_number < 25);`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*Player
	for rows.Next() {
		player := Player{}
		var encodedSkills sql.NullString
		var encodedState sql.NullString
		var playerID sql.NullString
		err = rows.Scan(
			&playerID,
			&player.Defence,
			&player.Speed,
			&player.Pass,
			&player.Shoot,
			&player.Endurance,
			&player.ShirtNumber,
			&player.PreferredPosition,
			&encodedSkills,
			&encodedState,
			&player.Potential,
			&player.Name,
			&player.DayOfBirth,
			&player.RedCard,
			&player.InjuryMatchesLeft,
			&player.Tiredness,
			&player.CountryOfBirth,
			&player.Race,
			&player.YellowCard1stHalf,
		)
		player.TeamId = teamID
		player.EncodedSkills, _ = new(big.Int).SetString(encodedSkills.String, 10)
		player.EncodedState, _ = new(big.Int).SetString(encodedState.String, 10)
		player.PlayerId, _ = new(big.Int).SetString(playerID.String, 10)
		players = append(players, &player)
	}
	return players, err
}

func PlayersBulkInsertUpdate(rowsToBeInserted []Player, tx *sql.Tx) error {
	numParams := 20
	valueStrings := make([]string, 0, len(rowsToBeInserted))
	valueArgs := make([]interface{}, 0, len(rowsToBeInserted)*numParams)
	i := 0
	for _, post := range rowsToBeInserted {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*numParams+1, i*numParams+2, i*numParams+3, i*numParams+4, i*numParams+5, i*numParams+6, i*numParams+7, i*numParams+8, i*numParams+9, i*numParams+10, i*numParams+11, i*numParams+12, i*numParams+13, i*numParams+14, i*numParams+15, i*numParams+16, i*numParams+17, i*numParams+18, i*numParams+19, i*numParams+20))
		valueArgs = append(valueArgs, post.PlayerId.String())
		valueArgs = append(valueArgs, post.TeamId)
		valueArgs = append(valueArgs, post.Defence)
		valueArgs = append(valueArgs, post.Speed)
		valueArgs = append(valueArgs, post.Pass)
		valueArgs = append(valueArgs, post.Shoot)
		valueArgs = append(valueArgs, post.Endurance)
		valueArgs = append(valueArgs, post.ShirtNumber)
		valueArgs = append(valueArgs, post.EncodedSkills.String())
		valueArgs = append(valueArgs, post.RedCard)
		valueArgs = append(valueArgs, post.InjuryMatchesLeft)
		valueArgs = append(valueArgs, post.Name)
		valueArgs = append(valueArgs, post.Tiredness)
		valueArgs = append(valueArgs, post.CountryOfBirth)
		valueArgs = append(valueArgs, post.Race)
		valueArgs = append(valueArgs, post.YellowCard1stHalf)
		valueArgs = append(valueArgs, post.PreferredPosition)
		valueArgs = append(valueArgs, post.Potential)
		valueArgs = append(valueArgs, post.DayOfBirth)
		valueArgs = append(valueArgs, post.EncodedState.String())
		i++
	}
	stmt := fmt.Sprintf(`INSERT INTO players (
		player_id,
		team_id, 
		defence, 
		speed, 
		pass, 
		shoot,
		endurance,
		shirt_number,
		encoded_skills,
		red_card,
		injury_matches_left,
		name,
		tiredness,
		country_of_birth,
		race,
		yellow_card_1st_half,
		preferred_position,
		potential,
		day_of_birth,
		encoded_state
		) VALUES %s
		ON CONFLICT(player_id) DO UPDATE SET
		team_id = excluded.team_id, 
		defence = excluded.defence, 
		speed = excluded.speed, 
		pass = excluded.pass, 
		shoot = excluded.shoot,
		endurance = excluded.endurance,
		shirt_number = excluded.shirt_number,
		encoded_skills = excluded.encoded_skills,
		red_card = excluded.red_card,
		injury_matches_left = excluded.injury_matches_left,
		name = excluded.name,
		tiredness = excluded.tiredness,
		country_of_birth = excluded.country_of_birth,
		race = excluded.race,
		yellow_card_1st_half = excluded.yellow_card_1st_half,
		preferred_position = excluded.preferred_position,
		potential = excluded.potential,
		day_of_birth = excluded.day_of_birth,
		encoded_state = excluded.encoded_state
		`, strings.Join(valueStrings, ","))
	_, err := tx.Exec(stmt, valueArgs...)
	return err
}
