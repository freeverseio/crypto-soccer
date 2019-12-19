package process_test

import (
	"log"
	"os"
	"testing"

	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

var universedb *storage.Storage
var relaydb *relay.Storage

func TestMain(m *testing.M) {
	var err error
	universedb, err = storage.NewPostgres("postgres://freeverse:freeverse@localhost:15432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	relaydb, err = relay.NewPostgres("postgres://freeverse:freeverse@localhost:15433/relay?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
