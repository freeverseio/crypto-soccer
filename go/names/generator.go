package names

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"math/big"
	"os"
)

type Generator struct {
	db *sql.DB
}

func New() (*Generator, error) {
	var err error
	generator := Generator{}
	generator.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	if err := generator.db.Ping(); err != nil {
		return nil, err
	}
	_, err = generator.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}
	file, err := os.Open("./sql/00_goalRev.sql")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	script, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	_, err = generator.db.Exec(string(script))
	if err != nil {
		return nil, err
	}
	return &generator, nil
}

func GeneratePlayerName(playerId *big.Int) string {
	_ = playerId
	return "s"
	//	return sillyname.GenerateStupidName()
}

func GenerateTeamName(teamId *big.Int) string {
	_ = teamId
	return "s"
	//	return sillyname.GenerateStupidName()
}
