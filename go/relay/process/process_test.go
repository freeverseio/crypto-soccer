package relay_test

import (
	//"math/big"
	"testing"

	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/core/types"
	//"github.com/ethereum/go-ethereum/crypto"

	"github.com/freeverseio/crypto-soccer/go/relay/process"
	"github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestSubmitActionRoot(t *testing.T) {
	db, err := storage.NewSqlite3("../db/00_schema.sql")
	// db, err := storage.NewPostgres("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
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

	p, err := relay.NewProcessor(bc.Client, bc.Owner, db, bc.Updates)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Process()
	if err != nil {
		t.Fatal(err)
	}
}
