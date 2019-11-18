package names

import (
	"database/sql"
	"errors"
	"math/big"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Generator struct {
	db *sql.DB
}

func New() (*Generator, error) {
	var err error
	log.Info("new")
	generator := Generator{}
	generator.db, err = sql.Open("sqlite3", "./sql/00_goalRev.db")
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
	return &generator, nil
}

func (b *Generator) NamesCount(teamId *big.Int) (uint64, error) {
	count := uint64(0)
	var err error
	rows, err := b.db.Query(`SELECT COUNT(*) FROM names WHERE country_id = $1;`, teamId.String())
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *Generator) GeneratePlayerName(playerId *big.Int, teamId *big.Int) (string, error) {
	_ = playerId
	log.Debugf("[NAMES] GeneratePlayerName of playerId %v", playerId)
	_, err := b.NamesCount(teamId)
	if err != nil {
		return "", err
	}
	rows, err := b.db.Query(`SELECT name FROM names WHERE country_id = $1;`, teamId.String())
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if !rows.Next() {
		return "", errors.New("Unexistent playerId")
	}
	var name string
	rows.Scan(&name)
	return name, nil
}

func GenerateTeamName(teamId *big.Int) string {
	_ = teamId
	return "s"
	//	return sillyname.GenerateStupidName()
}
