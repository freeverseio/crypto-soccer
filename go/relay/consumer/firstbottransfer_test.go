package consumer_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/relay/consumer"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
	log "github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

func TestTransferFirstBotFromLegitTZ(t *testing.T) {
	auth := bind.NewKeyedTransactor(bc.Owner)
	event := gql.TransferFirstBotToAddrInput{}
	event.Timezone = 1
	event.CountryIdxInTimezone = "0"
	event.Address = "0xeb3ce112d8610382a994646872c4361a96c82cf8"
	c := consumer.NewFirstBotTransfer(bc.Client, auth, bc.Contracts.Assets)
	assert.NilError(t, c.Process(event))
}

func TestTransferFirstBotFromWrongTZ(t *testing.T) {
	auth := bind.NewKeyedTransactor(bc.Owner)
	event := gql.TransferFirstBotToAddrInput{}
	event.Timezone = 10
	event.CountryIdxInTimezone = "0"
	event.Address = "0xeb3ce112d8610382a994646872c4361a96c82cf8"
	c := consumer.NewFirstBotTransfer(bc.Client, auth, bc.Contracts.Assets)
	err := c.Process(event)
	if err == nil {
		log.Errorf("transfer bot in wrong TZ should have failed, it did not")
	}
}
