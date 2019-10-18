package process_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
	"github.com/freeverseio/crypto-soccer/relay/process"
	"github.com/freeverseio/crypto-soccer/relay/storage"
)

func TestSyncTeams(t *testing.T) {
	storage, err := storage.NewSqlite3("../db/00_schema.sql")
	// storage, err := storage.NewPostgres("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}

	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	timezoneIdx := uint8(1)
	err = bc.InitOneTimezone(timezoneIdx)
	if err != nil {
		t.Fatal(err)
	}

	p, err := process.NewProcessor(bc.Client, bc.Owner, bc.Updates)
	if err != nil {
		t.Fatal(err)
	}

	err := p.Process()
	if err != nil {
		t.Fatal(err)
	}
}
