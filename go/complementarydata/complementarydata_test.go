package complementarydata_test

import (
	"encoding/json"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/complementarydata"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql/input"
	"gotest.tools/assert"
)

func TestComplementaryDataMarshall(t *testing.T) {
	data := complementarydata.ComplementaryData{}
	result, err := json.Marshal(data)
	assert.NilError(t, err)
	assert.Equal(t, `[]`, string(result))

	assert.NilError(t, data.Push(input.SetTeamNameInput{}))
	result, err = json.Marshal(data)
	assert.NilError(t, err)
	assert.Equal(t, `[{"Name":"SetTeamNameInput","Data":{"Signature":"","TeamId":"","Name":""}}]`, string(result))

	assert.NilError(t, data.Push(input.SetTeamManagerNameInput{}))
	result, err = json.Marshal(data)
	assert.NilError(t, err)
	assert.Equal(t, `[{"Name":"SetTeamNameInput","Data":{"Signature":"","TeamId":"","Name":""}},{"Name":"SetTeamManagerNameInput","Data":{"Signature":"","TeamId":"","Name":""}}]`, string(result))
}
