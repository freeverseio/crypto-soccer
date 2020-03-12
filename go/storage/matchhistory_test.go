package storage_test

import (
	"testing"

	"gotest.tools/assert"
)

func TestMatchHistory(t *testing.T) {
	t.Parallel()
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)
}
