package relay_test

import (
	//"math/big"
	"testing"

	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/core/types"
	//"github.com/ethereum/go-ethereum/crypto"

	relay "github.com/freeverseio/crypto-soccer/go/relay/process"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestSubmitActionRoot(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
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

	p, err := relay.NewProcessor(bc.Client, bc.Owner, db, bc.Contracts.Updates, "localhost:5001")
	if err != nil {
		t.Fatal(err)
	}

	err = p.Process(tx)
	if err != nil {
		t.Fatal(err)
	}
}
