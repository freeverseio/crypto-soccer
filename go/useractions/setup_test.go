package useractions_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"

	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = storage.New("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
