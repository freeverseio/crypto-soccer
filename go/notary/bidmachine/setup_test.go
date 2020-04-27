package bidmachine_test

import (
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/testutils"
)

// var db *sql.DB
var bc *testutils.BlockchainNode

func TestMain(m *testing.M) {
	var err error
	// db, err = storage.New("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	bc, err = testutils.NewBlockchainNode()
	if err != nil {
		log.Fatal(err)
	}
	bc.DeployContracts(bc.Owner)
	os.Exit(m.Run())
}
