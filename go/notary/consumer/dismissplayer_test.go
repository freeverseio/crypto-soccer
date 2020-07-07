package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestDismissPlayer(t *testing.T) {
	in := input.DismissPlayerInput{}
	in.PlayerId = "123455"
	in.ValidUntil = "5646456"
	in.ReturnToAcademy = true
	in.Signature = "0f13e4028d911bbf7e305267d593c6b67888030032e73f94a5cf8af204567ab629848e9290568aa5d19c1b7a4761a20ed4059072aacd79bde56e1b52c17a21311b"

	assert.Error(t, consumer.DismissPlayer(*bc.Contracts, bc.Owner, in), "failed to estimate gas needed: The execution failed due to an exception.")
}
