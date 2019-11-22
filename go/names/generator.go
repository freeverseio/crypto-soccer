package names

import (
	"database/sql"
	"errors"
	"hash/fnv"
	"math/big"
	"strconv"

	"github.com/Pallinder/sillyname-go"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Generator struct {
	db                    *sql.DB
	countryCodes4Names    []uint
	countryCodes4Surnames []uint
	namesInCountry        map[uint]uint
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
	if err := generator.countDB(); err != nil {
		return nil, err
	}
	return &generator, nil
}

func (b *Generator) countDB() error {
	var err error
	rows, err := b.db.Query(`SELECT country_code, num_names FROM countries`)
	if err != nil {
		return err
	}
	defer rows.Close()
	var country_code uint
	var num_names uint
	b.namesInCountry = make(map[uint]uint)
	for rows.Next() {
		err = rows.Scan(&country_code, &num_names)
		if err != nil {
			return err
		}
		// all country_codes for names have code < 1000, all surnames are > 1000
		if country_code < 1000 {
			b.countryCodes4Names = append(b.countryCodes4Names, country_code)
		} else {
			b.countryCodes4Surnames = append(b.countryCodes4Surnames, country_code)
		}
		b.namesInCountry[country_code] = num_names
	}
	return nil
}

func (b *Generator) GenerateRnd(seed *big.Int, max_val uint64, nLayers int) uint64 {
	var iterated_seed uint64
	iterated_seed = int_hash(seed.String())
	for i := 1; i < nLayers; i++ {
		iterated_seed = int_hash(big.NewInt(int64(iterated_seed)).String())
	}
	return iterated_seed % max_val
}

func (b *Generator) GenerateName(isSurname bool, playerId *big.Int, country_code uint, purity int) (string, error) {
	log.Debugf("[NAMES] GenerateName of playerId %v", playerId)
	nLayers1 := 1
	nLayers2 := 2
	nLayers3 := 1
	tableName := "names"
	colName := "name"
	codes := b.countryCodes4Names
	if isSurname {
		nLayers1 = 3
		nLayers2 = 4
		nLayers3 = 2
		tableName = "surnames"
		colName = "surname"
		codes = b.countryCodes4Surnames
	}
	dice := b.GenerateRnd(playerId, 100, nLayers1)
	if int(dice) > purity {
		// pick a different country
		var nCountryCodes = len(codes)
		var rnd_idx int = int(b.GenerateRnd(playerId, uint64(nCountryCodes), nLayers3))
		if country_code == codes[rnd_idx] {
			country_code = codes[(rnd_idx+1)%nCountryCodes]
		} else {
			country_code = codes[rnd_idx]
		}
	}
	var namesInCountry uint = b.namesInCountry[country_code]
	var idxInCountry uint64 = b.GenerateRnd(playerId, uint64(namesInCountry), nLayers2)
	rows, err := b.db.Query(`SELECT `+colName+` FROM `+tableName+` WHERE (country_code = $1 AND idx_in_country = $2)`, country_code, idxInCountry)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if !rows.Next() {
		var str string = "Rnd choice failed, country_code = " + strconv.FormatInt(int64(country_code), 10) +
			", idxInCountry = " + strconv.FormatInt(int64(idxInCountry), 10) +
			", tableName = " + tableName
		return "", errors.New(str)
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
	// country_id is an encoding of (tz, countryIdx):
	country_id = uint64(timezone)*1000000 + countryIdxInTZ
	// if the country is not defined, we use a default country: Spain, at tz = 19
	isSpecified, err := b.isCountrySpecified(country_id)
	if err != nil {
		return "", err
	}
	// Spain is the default country if you query for one that is not specified
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
	var code_name uint
	var code_surname uint
	var pure_pure int
	var pure_foreign int
	var foreign_pure int
	var foreign_foreign int
	defer rows.Close()
	if !rows.Next() {
		return "", errors.New("Cannot find specs for country_id = %s" + strconv.FormatInt(int64(country_id), 10))
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
	return sillyname.GenerateStupidName()
}
