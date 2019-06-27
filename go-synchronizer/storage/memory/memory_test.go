package memory

import (
	"math/big"
	"testing"
)

func TestTeamAdd(t *testing.T) {
	storage := New()
	err := storage.TeamAdd(1, "ciao")
	if err != nil {
		t.Fatal(err)
	}
	team, err := storage.GetTeam(1)
	if err != nil {
		t.Fatal(err)
	}
	if team.Id != 1 {
		t.Fatalf("Expected 0 result %v", team.Id)
	}
	if team.Name != "ciao" {
		t.Fatalf("Expected ciao result %v", team.Name)
	}
}

func TestGetUnexistentTeam(t *testing.T) {
	storage := New()
	_, err := storage.GetTeam(0)
	if err == nil {
		t.Fatal("No error on get unexistent team")
	}
}

func TestBlockNumber(t *testing.T) {
	storage := New()
	blockNumber, err := storage.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	if blockNumber != nil {
		t.Fatalf("Expected nil result %v", blockNumber)
	}

	err = storage.SetBlockNumber(big.NewInt(3))
	if err != nil {
		t.Fatal(err)
	}

	blockNumber, err = storage.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	if blockNumber.String() != "3" {
		t.Fatalf("Expected 3 result %v", blockNumber)
	}
}
