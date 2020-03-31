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
	if b.c != nil {
		select {
		case b.c <- input:
		default:
			log.Warning("TransferFirstBotToAddr: channel is full, discarding value")
			return false, errors.New("channel is full, discarding value")
		}
	}
	return true, nil
}
