package fillgamedb

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
)

type TeamService struct {
	universeClient *graphql.Client
	gameClient     *graphql.Client
}

type Team struct {
	TeamId      graphql.String
	Name        graphql.String
	ManagerName graphql.String
}
type AllTeams struct {
	AllTeams struct {
		TotalCount graphql.Int
		Nodes      []Team
	}
}

func NewTeamService(universeClient *graphql.Client, gameClient *graphql.Client) *TeamService {
	return &TeamService{universeClient, gameClient}
}

func (b TeamService) getAllTeams() (AllTeams, error) {
	var allTeams AllTeams

	err := b.universeClient.Query(context.Background(), &allTeams, nil)

	if err != nil {
		return allTeams, err
	}

	fmt.Println("Teams TotalCount", allTeams.AllTeams.TotalCount)

	return allTeams, nil
}

func (b TeamService) createTeamProp(team Team) error {
	var mutation struct {
		CreateTeamProp struct {
			ClientMutationId graphql.ID
		} `graphql:"createTeamProp(input : { teamProp: { teamId: $teamId, teamName: $teamName, teamManagerName: $teamManagerName }})"`
	}
	variables := map[string]interface{}{
		"teamId":          team.TeamId,
		"teamName":        team.Name,
		"teamManagerName": team.ManagerName,
	}
	err := b.gameClient.Mutate(context.Background(), &mutation, variables)
	if err != nil {
		return err
	}
	fmt.Println("Created team: ", team)

	return nil
}

func (b TeamService) updateTeamProp(team Team) error {
	var mutation struct {
		UpdateTeamPropByTeamId struct {
			ClientMutationId graphql.ID
		} `graphql:"updateTeamPropByTeamId(input : { teamId: $teamId, teamPropPatch: { teamName: $teamName, teamManagerName: $teamManagerName }})"`
	}
	variables := map[string]interface{}{
		"teamId":          team.TeamId,
		"teamName":        team.Name,
		"teamManagerName": team.ManagerName,
	}
	err := b.gameClient.Mutate(context.Background(), &mutation, variables)
	if err != nil {
		return err
	}
	fmt.Println("Updated team: ", team)

	return nil
}

func (b TeamService) upsertAllTeamProps(teams AllTeams) error {
	for _, team := range teams.AllTeams.Nodes {
		err := b.createTeamProp(team)
		if err != nil {
			err = b.updateTeamProp(team)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
