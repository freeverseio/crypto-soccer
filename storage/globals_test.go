package storage

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTestDB(t *testing.T) *Storage {
	tmp, err := ioutil.TempDir("", "dbtest")
	assert.Nil(t, err)
	s, err := New(tmp)
	assert.Nil(t, err)
	err = s.SetGlobals(GlobalsEntry{
		CurrentQuota: 0,
	})
	assert.Nil(t, err)
	return s
}

func TestSetGlobals(t *testing.T) {
	s := CreateTestDB(t)

	err := s.SetGlobals(GlobalsEntry{})
	assert.Nil(t, err)

	g, err := s.Globals()
	assert.Nil(t, err)
}
