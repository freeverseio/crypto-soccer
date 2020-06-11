package complementarydata_test

import (
	"encoding/json"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/complementarydata"
	"gotest.tools/assert"
)

func TestStackMarshall(t *testing.T) {
	stack := complementarydata.ComplementaryData{}
	result, err := json.Marshal(stack)
	assert.NilError(t, err)
	assert.Equal(t, "[]", string(result))

	stack.Push(3)
	result, err = json.Marshal(stack)
	assert.NilError(t, err)
	assert.Equal(t, "[3]", string(result))

	stack.Push(struct{ Ciao int }{5})
	result, err = json.Marshal(stack)
	assert.NilError(t, err)
	assert.Equal(t, `[3,{"Ciao":5}]`, string(result))
}
