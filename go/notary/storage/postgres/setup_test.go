package postgres_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = postgres.New("postgres://freeverse:freeverse@localhost:5432/market?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
