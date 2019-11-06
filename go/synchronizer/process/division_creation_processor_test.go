package process_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestDivisionCreationProcess(t *testing.T) {
	db, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	process, err := process.NewDivisionCreationProcessor(db, bc.Assets, bc.Leagues)
	if err != nil {
		t.Fatal(err)
	}

	event := assets.AssetsDivisionCreation{
		Timezone:             uint8(1),
		CountryIdxInTZ:       big.NewInt(0),
		DivisionIdxInCountry: big.NewInt(0),
	}
	err = process.Process(event)
	if err != nil {
		t.Fatal(err)
	}
	player, err := db.GetPlayer(big.NewInt(274877906944))
	if err != nil {
		t.Fatal(err)
	}
	if player.Name == "" {
		t.Fatal("name is empty")
	}
	if player.DayOfBirth == 0 {
		t.Fatal("dayOfBirth is 0")
	}
}
