package storage_test

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var dump spew.ConfigState

func TestMain(m *testing.M) {
	dump = spew.ConfigState{DisablePointerAddresses: true, Indent: "\t"}
	os.Exit(m.Run())
}
