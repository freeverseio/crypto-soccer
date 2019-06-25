package memory

import "testing"

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
