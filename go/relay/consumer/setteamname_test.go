package consumer_test

import (
	"errors"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/consumer"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/storage/mock"
	"gotest.tools/assert"
)

func TestSetTeamNameOfUnexistentTeam(t *testing.T) {
	teamStorageService := mock.TeamStorageService{
		UpdateNameFunc: func(teamId string, name string) error {
			return errors.New("unexistent team")
		},
	}
	event := input.SetTeamNameInput{}
	assert.Error(t, consumer.SetTeamName(teamStorageService, event), "unexistent team")
}

func TestSetTeamNameOfExistentTeam(t *testing.T) {
	var calledTeamId string
	var calledName string
	teamStorageService := mock.TeamStorageService{
		UpdateNameFunc: func(teamId string, name string) error {
			calledTeamId = teamId
			calledName = name
			return nil
		},
	}

	event := input.SetTeamNameInput{
		TeamId: "434234245345",
		Name:   "ippo",
	}
	assert.NilError(t, consumer.SetTeamName(teamStorageService, event))
	assert.Equal(t, calledTeamId, string(event.TeamId))
	assert.Equal(t, calledName, event.Name)
}
