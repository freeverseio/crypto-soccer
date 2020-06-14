package complementarydata

import (
	"encoding/json"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
)

type event struct {
	Name string
	Data json.RawMessage
}

type ComplementaryData []event

type ComplementaryDataPublishService interface {
	Publish(data ComplementaryData) (string, error)
	Retrive(id string) (*ComplementaryData, error)
}

func (b *ComplementaryData) Push(i interface{}) error {
	e := event{}
	switch v := i.(type) {
	case input.SetTeamNameInput:
		e.Name = "SetTeamNameInput"
	case input.SetTeamManagerNameInput:
		e.Name = "SetTeamManagerNameInput"
	default:
		return fmt.Errorf("unknown type %T", v)
	}

	var err error
	e.Data, err = json.Marshal(i)
	if err != nil {
		return err
	}
	*b = append(*b, e)
	return nil
}
