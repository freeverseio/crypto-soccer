package input_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"gotest.tools/assert"
)

func TestSetTeamManagerNameInputHash(t *testing.T) {
	in := input.SetTeamManagerNameInput{}
	hash, err := in.Hash()
	assert.Error(t, err, "Invalid TeamId")

	in.TeamId = "3"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x074b4277787bca36334cf57f0507141ef743a08d7690dba02af123626e6955d0")

	in.Name = "ciao"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x70203f52f1e52e6727239d5197821e1ed77161d7a04a45ad1fb15702c69433d4")
}
