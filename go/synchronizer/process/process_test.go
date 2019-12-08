package process_test

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/names"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"

	"github.com/freeverseio/crypto-soccer/go/testutils"
	"github.com/ory/dockertest"
)

var universedb *storage.Storage
var relaydb *relay.Storage
var namesdb *names.Generator
var bc *testutils.BlockchainNode
var processor *process.EventProcessor

func TestMain(m *testing.M) {
	var err error
	namesdb, err = names.New("../../names/sql/names.db")
	if err != nil {
		log.Fatal(err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// Build and run the given Dockerfile
	err = pool.RemoveContainerByName("universe-test-image")
	if err != nil {
		log.Fatal(err)
	}
	universeResource, err := pool.BuildAndRun("universe-test-image", "../../../universe.db/Dockerfile", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	err = pool.RemoveContainerByName("relay-test-image")
	if err != nil {
		log.Fatal(err)
	}
	relayResource, err := pool.BuildAndRun("relay-test-image", "../../../relay.db/Dockerfile", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	var db *sql.DB
	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://freeverse:freeverse@localhost:%s/cryptosoccer?sslmode=disable", universeResource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	universedb, err = storage.New(db)
	if err != nil {
		log.Fatal(err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://freeverse:freeverse@localhost:%s/relay?sslmode=disable", relayResource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	relaydb, err = relay.New(db)
	if err != nil {
		log.Fatal(err)
	}

	bc, err = testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		log.Fatal(err)
	}

	processor, err = process.NewEventProcessor(
		bc.Contracts,
		universedb,
		relaydb,
		namesdb,
	)
	if err != nil {
		log.Fatal(err)
	}
	count, err := processor.Process(0)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		log.Fatal("processed 0 blocks")
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(universeResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(relayResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestSyncTeams(t *testing.T) {
	p, err := process.NewEventProcessor(
		bc.Contracts,
		universedb,
		relaydb,
		namesdb,
	)
	if err != nil {
		t.Fatal(err)
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
