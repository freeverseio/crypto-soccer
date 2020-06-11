package gql

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type TransferFirstBotToAddrInput struct {
	Timezone             int32
	CountryIdxInTimezone string
	Address              string
}

func (b *Resolver) TransferFirstBotToAddr(input TransferFirstBotToAddrInput) (bool, error) {
	if b.ch != nil {
		select {
		case b.ch <- input:
		default:
			log.Warning("TransferFirstBotToAddr: channel is full, discarding value")
			return false, errors.New("channel is full, discarding value")
		}
	}
	return true, nil
}
