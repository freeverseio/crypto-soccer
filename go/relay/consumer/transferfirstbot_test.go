package consumer_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/relay/consumer"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
	"gotest.tools/assert"
)

func TestTransferFirstBot(t *testing.T) {
	auth := bind.NewKeyedTransactor(bc.Owner)
	event := gql.TransferFirstBotToAddrInput{}
	event.Timezone = 10
	event.CountryIdxInTimezone = "0"
	event.Address = "0xeb3ce112d8610382a994646872c4361a96c82cf8"
	c := consumer.NewTransferFirstBot(bc.Client, auth, bc.Contracts.Assets)
	assert.NilError(t, c.Process(event))
}
