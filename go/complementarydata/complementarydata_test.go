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
	assert.Equal(t, "[]", string(result))

	data.Push(3)
	result, err = json.Marshal(data)
	assert.NilError(t, err)
	assert.Equal(t, "[3]", string(result))

	data.Push(struct{ Ciao int }{5})
	result, err = json.Marshal(data)
	assert.NilError(t, err)
	assert.Equal(t, `[3,{"Ciao":5}]`, string(result))
}
