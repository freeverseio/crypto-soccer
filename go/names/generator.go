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

func New(db_filename string) (*Generator, error) {
	var err error
	generator := Generator{}
	generator.db, err = sql.Open("sqlite3", db_filename)
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

// comparer is either "=" or "!="
func (b *Generator) NamesCount(tableName string, condition string) (uint64, error) {
	count := uint64(0)
	var err error
	var cmd string = `SELECT COUNT(*) FROM ` + tableName + ` ` + condition
	rows, err := b.db.Query(cmd)
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

func (b *Generator) GenerateRnd(seed *big.Int, max_val uint64, nLayers int) uint64 {
	var iterated_seed uint64
	iterated_seed = int_hash(seed.String())
	for i := 1; i < nLayers; i++ {
		iterated_seed = int_hash(big.NewInt(int64(iterated_seed)).String())
	}
	return iterated_seed % max_val
}

func (b *Generator) GenerateName(isSurname bool, playerId *big.Int, country_code int, purity int) (string, error) {
	log.Debugf("[NAMES] GenerateName of playerId %v", playerId)
	nLayers1 := 1
	nLayers2 := 2
	tableName := "names"
	colName := "name"
	if isSurname {
		nLayers1 = 3
		nLayers2 = 4
		tableName = "surnames"
		colName = "surname"
	}

	dice := b.GenerateRnd(playerId, 100, nLayers1)
	var condition string = `WHERE country_code = ` + strconv.Itoa(country_code) + ";"
	if int(dice) > purity {
		condition = `WHERE country_code != ` + strconv.Itoa(country_code) + ";"
	}
	num_names, err := b.NamesCount(tableName, condition)
	if err != nil {
		return "", err
	}
	rows, err := b.db.Query(`SELECT ` + colName + ` FROM ` + tableName + ` ` + condition)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var selected_player uint64 = b.GenerateRnd(playerId, num_names, nLayers2)
	var i uint64
	for i = 0; i <= selected_player; i++ {
		if !rows.Next() {
			return "", errors.New("Rnd choice selected a player too too far in the database")
		}
	}
	var name string
	rows.Scan(&name)
	return name, nil
}

// comparer is either "=" or "!="
func (b *Generator) isCountrySpecified(country_id uint64) (bool, error) {
	var err error
	rows, err := b.db.Query(`SELECT COUNT(*) FROM country_specs WHERE tz_idx = $1;`, strconv.FormatInt(int64(country_id), 10))
	if err != nil {
		return false, err
	}
	defer rows.Close()
	rows.Next()
	count := uint64(0)
	err = rows.Scan(&count)
	if err != nil {
		return false, err
	}
	return (count == 1), nil
}

func (b *Generator) GeneratePlayerFullName(playerId *big.Int, timezone uint8, countryIdxInTZ uint64) (string, error) {
	log.Debugf("[NAMES] GeneratePlayerFullName of playerId %v", playerId)
	var country_id uint64
	country_id = uint64(timezone)*1000000 + countryIdxInTZ
	// if the country is not defined, we use a default country: Spain, at tz = 19
	isSpecified, err := b.isCountrySpecified(country_id)
	if err != nil {
		return "", err
	}
	if !isSpecified {
		country_id = uint64(19)*1000000 + 0
	}
	rows, err := b.db.Query(`SELECT 
		code_name,
		code_surname,
		pure_pure,
		pure_foreign,
		foreign_pure,
		foreign_foreign
	FROM country_specs WHERE tz_idx = $1;`, strconv.FormatInt(int64(country_id), 10))
	if err != nil {
		return "", err
	}
	var code_name int
	var code_surname int
	var pure_pure int
	var pure_foreign int
	var foreign_pure int
	var foreign_foreign int
	defer rows.Close()
	if !rows.Next() {
		return "", errors.New("Rnd choice selected a player too too far in the database")
	}
	err = rows.Scan(&code_name, &code_surname, &pure_pure, &pure_foreign, &foreign_pure, &foreign_foreign)
	if err != nil {
		return "", err
	}
	name, err := b.GenerateName(false, playerId, code_name, pure_pure+pure_foreign)
	if err != nil {
		return "", err
	}
	surname, err := b.GenerateName(true, playerId, code_surname, pure_pure+foreign_pure)
	if err != nil {
		return "", err
	}
	return name + " " + surname, nil
}

func GenerateTeamName(teamId *big.Int) string {
	_ = teamId
	return "s"
	//	return sillyname.GenerateStupidName()
}
