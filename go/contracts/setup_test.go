package contracts_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	log "github.com/sirupsen/logrus"
)

var bc *contracts.Contracts
var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	bc, err = contracts.NewByTruffle()
	if err != nil {
		log.Fatal(err)
	}
	db, err = postgres.New("postgres://freeverse:freeverse@crypto-soccer_devcontainer_dockerhost_1:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
