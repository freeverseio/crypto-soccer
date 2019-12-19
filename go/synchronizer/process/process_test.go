package process_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"

	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestSyncTeams(t *testing.T) {
	err := universedb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer universedb.Rollback()
	err = relaydb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer relaydb.Rollback()
	// storage, err := storage.NewPostgres("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	namesdb, err := names.New("../../names/sql/names.db")
	if err != nil {
		t.Fatal(err)
	}
	p, err := process.NewEventProcessor(
		bc.Contracts,
		universedb,
		relaydb,
		namesdb,
	)
	if err != nil {
		t.Fatal(err)
	}

	count, err := p.Process(0)
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("processed 0 blocks")
	}

	// the null timezone (0) is only used by the Academy Team
	if count, err := universedb.TimezoneCount(); err != nil {
		t.Fatal(err)
	} else if count != 2 {
		t.Fatalf("Expected 2 time zones at time of creation,  actual %v", count)
	}

	// one country belongs to timezone = 0
	if count, err := universedb.CountryCount(); err != nil {
		t.Fatal(err)
	} else if count != 2 {
		t.Fatalf("Expected 2 countries at time of creation,  actual %v", count)
	}

	// one team (the Academy) belongs to timezone = 0
	if count, err := universedb.TeamCount(); err != nil {
		t.Fatal(err)
	} else if count != (128 + 1) {
		t.Fatalf("Expected 128 actual %v", count)
	}
	if count, err := universedb.PlayerCount(); err != nil {
		t.Fatal(err)
	} else if count != 128*18 {
		t.Fatalf("Expected 128*18=2304 actual %v", count)
	}

	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)

	tx, err := bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(bc.Owner.PublicKey),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = helper.WaitReceipt(bc.Client, tx, 3)
	if err != nil {
		t.Fatal(err)
	}
	var txs []*types.Transaction
	for i := 0; i < 24*4; i++ {
		var root [32]byte
		tx, err := bc.Contracts.Updates.SubmitActionsRoot(
			bind.NewKeyedTransactor(bc.Owner),
			root,
			"cid",
		)
		if err != nil {
			t.Fatal(err)
		}
		txs = append(txs, tx)
	}
	err = helper.WaitReceipts(bc.Client, txs, 3)
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.Process(0)
	if err != nil {
		t.Fatal(err)
	}

	playerIdx := big.NewInt(30)
	playerID, err := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, playerIdx)
	if err != nil {
		t.Fatal(err)
	}
	owner, err := bc.Contracts.Assets.GetOwnerPlayer(&bind.CallOpts{}, playerID)
	if err != nil {
		t.Fatal(err)
	}
	if owner.String() != storage.BotOwner {
		t.Fatalf("Owner is wrong %v", owner.String())
	}

	tx, err = bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(bc.Owner.PublicKey),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = helper.WaitReceipt(bc.Client, tx, 3)
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Process(0)
	if err != nil {
		t.Fatal(err)
	}

	owner, err = bc.Contracts.Assets.GetOwnerPlayer(&bind.CallOpts{}, playerID)
	if err != nil {
		t.Fatal(err)
	}
	if owner != crypto.PubkeyToAddress(bc.Owner.PublicKey) {
		t.Fatalf("Owner is wrong %v", owner.String())
	}
}
