package complementarydata

import (
	"encoding/json"

	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
)

type event struct {
	Name string
	Data json.RawMessage
}

type ComplementaryData []event

func (b *ComplementaryData) PushSetTeamNameInput(in input.SetTeamNameInput) error {
	ev := event{}
	ev.Name = "SetTeamNameInput"
	var err error
	ev.Data, err = json.Marshal(in)
	if err != nil {
		return err
	}
	*b = append(*b, ev)
	return nil
}

func (b *ComplementaryData) PushSetTeamManagerNameInput(in input.SetTeamManagerNameInput) error {
	ev := event{}
	ev.Name = "SetTeamManagerNameInput"
	var err error
	ev.Data, err = json.Marshal(in)
	if err != nil {
		return err
	}
	*b = append(*b, ev)
	return nil
}
