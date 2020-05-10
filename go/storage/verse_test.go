package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestVerseInsert(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	verse := storage.Verse{}
	verse.VerseNumber = 2
	verse.Root = "453xe"
	assert.NilError(t, verse.Insert(tx))

	result, err := storage.VerseByNumber(tx, verse.VerseNumber)
	assert.NilError(t, err)
	assert.Equal(t, *result, verse)
}
