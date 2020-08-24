package fillgamedb

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
)

type PlayerService struct {
	universeClient *graphql.Client
	gameClient     *graphql.Client
}

type Player struct {
	PlayerId graphql.String
	Name     graphql.String
}
type AllPlayers struct {
	AllPlayers struct {
		TotalCount graphql.Int
		Nodes      []Player
	}
}

func NewPlayerService(universeClient *graphql.Client, gameClient *graphql.Client) *PlayerService {
	return &PlayerService{universeClient, gameClient}
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

func (b PlayerService) createPlayerProp(player Player) error {
	var mutation struct {
		CreatePlayerProp struct {
			ClientMutationId graphql.ID
		} `graphql:"createPlayerProp(input : { playerProp: { playerId: $playerId, playerName: $playerName }})"`
	}
	variables := map[string]interface{}{
		"playerId":   player.PlayerId,
		"playerName": player.Name,
	}
	err := b.gameClient.Mutate(context.Background(), &mutation, variables)
	if err != nil {
		return err
	}
	fmt.Println("Created player: ", player)

	return nil
}

func (b PlayerService) updatePlayerProp(player Player) error {
	var mutation struct {
		UpdatePlayerPropByPlayerId struct {
			ClientMutationId graphql.ID
		} `graphql:"updatePlayerPropByPlayerId(input : { playerId: $playerId, playerPropPatch: { playerName: $playerName }})"`
	}
	variables := map[string]interface{}{
		"playerId":   player.PlayerId,
		"playerName": player.Name,
	}
	err := b.gameClient.Mutate(context.Background(), &mutation, variables)
	if err != nil {
		return err
	}
	fmt.Println("Updated player: ", player)

	return nil
}

func (b PlayerService) upsertAllPlayerProps(players AllPlayers) error {
	for _, player := range players.AllPlayers.Nodes {
		err := b.createPlayerProp(player)
		if err != nil {
			err = b.updatePlayerProp(player)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
