package sqlite3_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage/sqlite3"
)

func TestNew(t *testing.T) {
	_, err := sqlite3.New("../../../postgres/sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTeamCount(t *testing.T) {
	storage, err := sqlite3.New("../../../postgres/sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.TeamCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 result %v", count)
	}
}

func TestTeamAdd(t *testing.T) {
	storage, err := sqlite3.New("../../../postgres/sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	err = storage.TeamAdd(3, "ciao")
	if err != nil {
		t.Fatal(err)
	}
	count, err := storage.TeamCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 result %v", count)
	}
}

func TestGetTeam(t *testing.T) {
	storage, err := sqlite3.New("../../../postgres/sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = storage.GetTeam(1)
	if err == nil {
		t.Fatal("getting team of unexistent team")
	}
	err = storage.TeamAdd(3, "ciao")
	if err != nil {
		t.Fatal(err)
	}
	team, err := storage.GetTeam(3)
	if err != nil {
		t.Fatal(err)
	}
	if team.Name != "ciao" {
		t.Fatalf("Expected ciao result %v", team.Name)
	}
}
