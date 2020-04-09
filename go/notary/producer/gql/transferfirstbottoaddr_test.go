package gql_test

import (
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"gotest.tools/assert"
)

func TestTransferFirstBot(t *testing.T) {
	t.Parallel()
	input := gql.TransferFirstBotToAddrInput{}
	resolver := gql.NewResolver(nil)
	result, err := resolver.TransferFirstBotToAddr(input)
	assert.NilError(t, err)
	assert.Equal(t, result, true)
}

func TestTransferFirstBotChannel(t *testing.T) {
	t.Parallel()
	c := make(chan interface{})
	resolver := gql.NewResolver(c)

	input := gql.TransferFirstBotToAddrInput{}
	input.Timezone = 23
	input.CountryIdxInTimezone = "4"
	input.Address = "sdfsgsd"
	go resolver.TransferFirstBotToAddr(input)

	select {
	case result := <-c:
		assert.Equal(t, input, result)
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}
