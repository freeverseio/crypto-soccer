package regenerateplayernames

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/shurcooL/graphql"
)

type PlayerService struct {
	universeDB     *sql.DB
	universeClient *graphql.Client
}

type TeamByTeamId struct {
	TimezoneIdx graphql.Int
	CountryIdx  graphql.Int
}

type Player struct {
	PlayerId          graphql.String
	Name              graphql.String
	TeamId            graphql.String
	Defence           graphql.Int
	Speed             graphql.Int
	Pass              graphql.Int
	Shoot             graphql.Int
	Endurance         graphql.Int
	ShirtNumber       graphql.Int
	EncodedSkills     graphql.String
	RedCard           graphql.Boolean
	InjuryMatchesLeft graphql.Int
	Tiredness         graphql.Int
	CountryOfBirth    graphql.String
	Race              graphql.String
	YellowCard1stHalf graphql.Boolean
	PreferredPosition graphql.String
	Potential         graphql.Int
	DayOfBirth        graphql.Int
	EncodedState      graphql.String
	TeamByTeamId      TeamByTeamId
}
type AllPlayers struct {
	AllPlayers struct {
		TotalCount graphql.Int
		Nodes      []Player
	}
}

func NewPlayerService(universeDB *sql.DB, universeClient *graphql.Client) *PlayerService {
	return &PlayerService{universeDB, universeClient}
}

func (b PlayerService) getAllPlayers() (AllPlayers, error) {
	var allPlayers AllPlayers

	err := b.universeClient.Query(context.Background(), &allPlayers, nil)

	if err != nil {
		return allPlayers, err
	}

	fmt.Println("Players TotalCount", allPlayers.AllPlayers.TotalCount)

	return allPlayers, nil
}

func (b PlayerService) updatetAllPlayerNamesAndRaces(players AllPlayers) error {
	var tz9Players []Player
	var tz7Players []Player
	var tz8Players []Player
	var tz11Players []Player
	for _, player := range players.AllPlayers.Nodes {
		if player.TeamByTeamId.TimezoneIdx == 9 {
			tz9Players = append(tz9Players, player)
		}
		if player.TeamByTeamId.TimezoneIdx == 7 {
			tz7Players = append(tz7Players, player)
		}
		if player.TeamByTeamId.TimezoneIdx == 8 {
			tz8Players = append(tz8Players, player)
		}
		if player.TeamByTeamId.TimezoneIdx == 11 {
			tz11Players = append(tz11Players, player)
		}
	}
	fmt.Printf("%v players for tz7\n", len(tz7Players))
	fmt.Printf("%v players for tz8\n", len(tz8Players))
	fmt.Printf("%v players for tz9\n", len(tz9Players))
	fmt.Printf("%v players for tz11\n", len(tz11Players))

	err := b.updatePlayers(tz11Players)
	if err != nil {
		return err
	}
	err = b.updatePlayers(tz9Players)
	if err != nil {
		return err
	}
	err = b.updatePlayers(tz7Players)
	if err != nil {
		return err
	}
	err = b.updatePlayers(tz8Players)
	if err != nil {
		return err
	}

	return nil
}

func (b PlayerService) updatePlayers(players []Player) error {
	generator, err := names.New("../../names/sql/names.db")
	if err != nil {
		return err
	}
	var playersToUpdate []Player
	for _, player := range players {
		var timezone uint8
		var countryIdxInTZ uint64
		generation := uint8(0)

		playerID, _ := new(big.Int).SetString(string(player.PlayerId), 10)
		timezone = uint8(player.TeamByTeamId.TimezoneIdx)
		countryIdxInTZ = uint64(player.TeamByTeamId.CountryIdx)
		name, _, region, err := generator.GeneratePlayerFullName(playerID, generation, timezone, countryIdxInTZ)
		if err != nil {
			return err
		}
		if len(name) == 0 {
			return errors.New("Expecting non empty player name")
		}
		player.Name = graphql.String(name)
		player.Race = graphql.String(region)
		playersToUpdate = append(playersToUpdate, player)
	}
	fmt.Println("Generated names")
	tx, err := b.universeDB.Begin()
	if err != nil {
		return err
	}
	err = b.PlayersBulkInsertUpdateNameRace(playersToUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (b PlayerService) updatePlayer(player Player) error {

	tx, err := b.universeDB.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE players SET
			name=$1,
			race=$2
			WHERE player_id=$3;`,
		player.Name,
		player.Race,
		string(player.PlayerId),
	); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (b PlayerService) PlayersBulkInsertUpdateNameRace(rowsToBeInserted []Player, tx *sql.Tx) error {
	numParams := 20
	var err error = nil
	maxRowsToBeInserted := int(35000 / numParams)
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
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*numParams+1, i*numParams+2, i*numParams+3, i*numParams+4, i*numParams+5, i*numParams+6, i*numParams+7, i*numParams+8, i*numParams+9, i*numParams+10, i*numParams+11, i*numParams+12, i*numParams+13, i*numParams+14, i*numParams+15, i*numParams+16, i*numParams+17, i*numParams+18, i*numParams+19, i*numParams+20))
			valueArgs = append(valueArgs, post.PlayerId)
			valueArgs = append(valueArgs, post.TeamId)
			valueArgs = append(valueArgs, post.Defence)
			valueArgs = append(valueArgs, post.Speed)
			valueArgs = append(valueArgs, post.Pass)
			valueArgs = append(valueArgs, post.Shoot)
			valueArgs = append(valueArgs, post.Endurance)
			valueArgs = append(valueArgs, post.ShirtNumber)
			valueArgs = append(valueArgs, post.EncodedSkills)
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
			valueArgs = append(valueArgs, post.EncodedState)
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
			name = excluded.name,
			race = excluded.race
			`, strings.Join(valueStrings, ","))
		_, err := tx.Exec(stmt, valueArgs...)
		if err != nil {
			return err
		}
		x = newX
	}
	return err
}
