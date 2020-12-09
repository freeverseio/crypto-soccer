package regenerateplayernames

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/shurcooL/graphql"
	"gopkg.in/src-d/go-log.v1"
)

type PlayerService struct {
	universeDB     *sql.DB
	universeClient *graphql.Client
	gameClient     *graphql.Client
}

type TeamByTeamId struct {
	TimezoneIdx graphql.String
	CountryIdx  graphql.String
}

type Player struct {
	PlayerId     graphql.String
	Name         graphql.String
	TeamByTeamId TeamByTeamId
}

type AllPlayers struct {
	AllPlayers struct {
		TotalCount graphql.Int
		Nodes      []Player
	}
}

func NewPlayerService(universeDB *sql.DB, universeClient *graphql.Client, gameClient *graphql.Client) *PlayerService {
	return &PlayerService{universeDB, universeClient, gameClient}
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
	for _, player := range players.AllPlayers.Nodes {
		err := b.updatePlayer(player)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b PlayerService) updatePlayer(player Player) error {
	//generate name and race
	generator, err := names.New("../../names/sql/names.db")
	if err != nil {
		return err
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	generation := uint8(0)
	playerId, _ := new(big.Int).SetString(string(player.PlayerId), 10)
	tz, err := strconv.ParseInt(string(player.TeamByTeamId.TimezoneIdx), 10, 8)
	if err != nil {
		return err
	}
	country, err := strconv.ParseInt(string(player.TeamByTeamId.CountryIdx), 10, 64)
	if err != nil {
		return err
	}
	timezone = uint8(tz)
	countryIdxInTZ = uint64(country)
	name, countryISO2, region, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
	if err != nil {
		return err
	}
	fmt.Println(name + " (" + countryISO2 + ", " + region + ") ")
	if len(name) == 0 {
		return errors.New("Expecting non empty player name")
	}

	//update name and race
	tx, err := b.universeDB.Begin()
	if err != nil {
		return err
	}
	log.Debugf("[DBMS] + update player id %v", player.PlayerId)
	if _, err := tx.Exec(`UPDATE players SET 
	name=$1, 
	race=$2
	WHERE player_id=$16;`,
		name,
		countryISO2,
		player.PlayerId,
	); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
