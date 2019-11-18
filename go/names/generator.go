package names

import (
	"database/sql"
	"errors"
	"hash/fnv"
	"math/big"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Generator struct {
	db *sql.DB
}

func int_hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func New() (*Generator, error) {
	var err error
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
	rows, err := b.db.Query(`SELECT COUNT(*) FROM names WHERE country_code = $1;`, teamId.String())
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

func (b *Generator) GenerateRnd(seed *big.Int, max_val uint64, layers int) uint64 {
	var iterated_seed uint64
	iterated_seed = int_hash(seed.String())
	for i := 1; i < layers; i++ {
		iterated_seed = int_hash(big.NewInt(int64(iterated_seed)).String())
	}
	return iterated_seed % max_val
}

func (b *Generator) GeneratePlayerName(playerId *big.Int, countryId *big.Int) (string, error) {
	_ = playerId
	log.Debugf("[NAMES] GeneratePlayerName of playerId %v", playerId)
	num_names, err := b.NamesCount(countryId)
	if err != nil {
		return "", err
	}
	rows, err := b.db.Query(`SELECT name FROM names WHERE country_code = $1;`, countryId.String())
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var selected_player uint64 = b.GenerateRnd(playerId, num_names, 1)
	var i uint64
	for i = 0; i < selected_player; i++ {
		if !rows.Next() {
			return "", errors.New("Rnd choice selected a player too too far in the database")
		}
	}
	var name string
	rows.Scan(&name)
	return name, nil
}

func (b *Generator) GeneratePlayerFullName(playerId *big.Int, timezone uint8, countryIdxInTZ uint64) (string, error) {
	log.Debugf("[NAMES] GeneratePlayerFullName of playerId %v", playerId)
	var country_id uint64
	country_id = uint64(timezone)*1000000 + countryIdxInTZ
	rows, err := b.db.Query(`SELECT 
		code_name,
		code_surname,
		pure_pure,
		pure_foreign,
		foreign_pure,
		foreign_foreign
	FROM country_specs WHERE tz_idx = $1;`, strconv.FormatInt(int64(country_id), 10))
	_ = rows
	if err != nil {
		return "", err
	}

	var name string
	// name, err := b.GeneratePlayerName(playerId, nameCountryIdName)
	// if err != nil {
	// 	return "", err
	// }
	return name, nil
}

func GenerateTeamName(teamId *big.Int) string {
	_ = teamId
	return "s"
	//	return sillyname.GenerateStupidName()
}
