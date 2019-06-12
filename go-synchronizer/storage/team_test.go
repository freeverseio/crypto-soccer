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

	if(!rows.Next()){
		t.Error("team not created")
	}

	var id int
	var name string
	err = rows.Scan(&id, &name)
	if err != nil {
		t.Error(err)
	}
	if id != 1 {
		t.Error("wrong id:", id)
	}
	if name != "Barca" {
		t.Error("wrong name:", name)
	}
}