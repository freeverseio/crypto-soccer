package consumer

import "github.com/freeverseio/crypto-soccer/go/relay/producer/gql"

type TransferFirstBot struct {
}

func NewTransferFirstBot() *TransferFirstBot {
	return &TransferFirstBot{}
}

func (b TransferFirstBot) Process(event gql.TransferFirstBotToAddrInput) error {
	return nil
}
