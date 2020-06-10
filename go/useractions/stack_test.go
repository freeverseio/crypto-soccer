package useractions_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/useractions"
	"gotest.tools/assert"
)

func TestStackMarshall(t *testing.T) {
	stack := useractions.Stack{}
	result, err := stack.Marshal()
	assert.NilError(t, err)
	assert.Equal(t, "[]", string(result))

	stack.Push(3)
	result, err = stack.Marshal()
	assert.NilError(t, err)
	assert.Equal(t, "[3]", string(result))

	stack.Push(struct{ Ciao int }{5})
	result, err = stack.Marshal()
	assert.NilError(t, err)
	assert.Equal(t, `[3,{"Ciao":5}]`, string(result))
}
