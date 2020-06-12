package complementarydata_test

import (
	"encoding/json"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/complementarydata"
	"gotest.tools/assert"
)

func TestStackMarshall(t *testing.T) {
	data := complementarydata.ComplementaryData{}
	result, err := json.Marshal(data)
	assert.NilError(t, err)
	assert.Equal(t, `{"SetTeamNameEvents":null,"SetTeamManagerNameEvents":null}`, string(result))
}
