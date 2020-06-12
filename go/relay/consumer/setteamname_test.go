package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/consumer"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/storage/memory"
	"github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestSetTeamNameOfUnexistentTeam(t *testing.T) {
	teamStorageService := memory.NewTeamStorageService()
	event := input.SetTeamNameInput{}
	assert.Error(t, consumer.SetTeamName(teamStorageService, event), "unexistent team")
}

func TestSetTeamName(t *testing.T) {
	teamStorageService := memory.NewTeamStorageService()
	team := storage.NewTeam()
	team.TeamID = "3"
	team.Name = "pippo"
	assert.NilError(t, teamStorageService.Insert(*team))

	event := input.SetTeamNameInput{}
	event.TeamId = graphql.ID(team.TeamID)
	event.Name = "eccolo"
	assert.NilError(t, consumer.SetTeamName(teamStorageService, event))

	result, err := teamStorageService.Team(team.TeamID)
	assert.NilError(t, err)
	assert.Equal(t, result.Name, event.Name)
}
