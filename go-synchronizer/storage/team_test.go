package storage

import (
	"testing"
)

func TestCreateTeam(t *testing.T) {
	url := "postgres://freeverse:freeverse@localhost/cryptosoccer"
	err := Init(url)
	if err != nil {
		t.Error(err)
	}

	err = CreateTeam(1, "Barca")
	if err != nil {
		t.Error(err)
	}

	rows, err := db.Query("SELECT * FROM teams WHERE id= '1';")
	if err != nil {
		t.Error(err)
	}

	if(rows.Next()){
		t.Error("I can get a line of empty db")
	}
}