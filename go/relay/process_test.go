package relay_test

import (
	//"math/big"
	"testing"

	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/core/types"
	//"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/relay"
	"gotest.tools/assert"
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
	conn, err := helper.NewParityBackend("http://localhost:8545")
	auth := conn.Transactor(common.HexToAddress("0xeb3ce112d8610382a994646872c4361a96c82cf8"))
	p, err := relay.NewProcessor(conn.Client, auth, db, bc.Contracts.Updates, "localhost:5001")
	if err != nil {
		t.Fatal(err)
	}

	err = p.Process(tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProcessNextUpdateTime(t *testing.T) {
	conn, err := helper.NewParityBackend("http://localhost:8545")
	auth := conn.Transactor(common.HexToAddress("0xeb3ce112d8610382a994646872c4361a96c82cf8"))
	p, err := relay.NewProcessor(conn.Client, auth, db, bc.Contracts.Updates, "localhost:5001")
	assert.NilError(t, err)
	nextDeadline, err := p.NextUpdateSinceEpochSec()
	assert.NilError(t, err)
	now := relay.NowSinceEpochSec()
	assert.Equal(t, nextDeadline, now)
}
