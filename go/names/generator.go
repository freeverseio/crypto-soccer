package names

import (
	"database/sql"
	"errors"
	"hash/fnv"
	"math/big"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Generator struct {
	db                    *sql.DB
	countryCodes4Names    []uint
	countryCodes4Surnames []uint
	namesInCountry        map[uint]uint
	nTeamnamesMain        uint
	nTeamnamesPreffix     uint
	nTeamnamesSuffix      uint
}


func int_hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func New(db_filename string) (*Generator, error) {
	var err error
	generator := Generator{}
	// PLAYER NAMES DB INIT
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
	if err := generator.countPlayersDB(); err != nil {
		return nil, err
	}
	if err := generator.countTeamsDB(); err != nil {
		return nil, err
	}

	return &generator, nil
}

func (b *Generator) countPlayersDB() error {
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

func (b *Generator) countTeamsDB() error {
	var err error
	// Count main names
	rows, err := b.db.Query(`SELECT COUNT(*) FROM team_mainnames;`)
	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()
	count := uint(0)
	err = rows.Scan(&count)
	if err != nil {
		return err
	}
	b.nTeamnamesMain = count
	// Count prefixes
	rows, err = b.db.Query(`SELECT COUNT(*) FROM team_prefixnames;`)
	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		return err
	}
	b.nTeamnamesPreffix = count
	// Count suffixes
	rows, err = b.db.Query(`SELECT COUNT(*) FROM team_suffixnames;`)
	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		return err
	}
	b.nTeamnamesSuffix = count

	return nil
}

func (b *Generator) GenerateRnd(seed *big.Int, salt string, max_val uint64) uint64 {
	var result uint64 = int_hash(seed.String() + salt)
	if max_val == 0 {
		return result
	}
	return result % max_val
}


func (b *Generator) GenerateSeedAndSalt(isSurname bool, playerId *big.Int, generation uint8, country_code uint, purity int) (*big.Int, string, uint, string, string, error) {
	log.Debugf("[NAMES] GenerateSeedAndSalt of playerId %v", playerId)
	isAcademyPlayer := generation > 31
	if isAcademyPlayer {
		generation = generation - 32
	}
	salt := "a"
	tableName := "names"
	colName := "name"
	codes := b.countryCodes4Names
	// ensure that names are always different for all generations
	seedTemp := b.GenerateRnd(playerId, "aa", 0) + uint64(generation)
	if isSurname {
		salt = "b"
		tableName = "surnames"
		colName = "surname"
		codes = b.countryCodes4Surnames
		seedTemp = b.GenerateRnd(playerId, "bb", 0) + uint64(generation)
		isActualSon := generation > 0 && !isAcademyPlayer
		if isActualSon {
			seedTemp -= 1
		}
	}
	seed := big.NewInt(int64(seedTemp))
	dice := b.GenerateRnd(seed, salt+"cc", 100)
	if int(dice) > purity {
		// pick a different country
		var nCountryCodes = len(codes)
		var rnd_idx int = int(b.GenerateRnd(seed, salt+"dd", uint64(nCountryCodes)))
		if country_code == codes[rnd_idx] {
			country_code = codes[(rnd_idx+1)%nCountryCodes]
		} else {
			country_code = codes[rnd_idx]
		}
	}
	return seed, salt, country_code, colName, tableName, nil
}

func (b *Generator) GenerateName(isSurname bool, playerId *big.Int, generation uint8, country_code uint, purity int) (string, uint, error) {
	log.Debugf("[NAMES] GenerateName of playerId %v", playerId)
	seed, salt, final_country_code, colName, tableName, err := b.GenerateSeedAndSalt(isSurname, playerId, generation, country_code, purity)
	
	var namesInCountry uint = b.namesInCountry[final_country_code]
	var idxInCountry uint64 = b.GenerateRnd(seed, salt+"ee", uint64(namesInCountry))
	rows, err := b.db.Query(`SELECT `+colName+` FROM `+tableName+` WHERE (country_code = $1 AND idx_in_country = $2)`, final_country_code, idxInCountry)
	if err != nil {
		return "", country_code, err
	}
	defer rows.Close()
	if !rows.Next() {
		var str string = "Rnd choice failed, final_country_code = " + strconv.FormatInt(int64(final_country_code), 10) +
			", idxInCountry = " + strconv.FormatInt(int64(idxInCountry), 10) +
			", tableName = " + tableName
		return "", country_code, errors.New(str)
	}
	var name string
	rows.Scan(&name)
	return name, country_code, nil
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

func encodeCountry(timezone uint8, countryIdxInTZ uint64) (uint64) {
	return uint64(timezone)*1000000 + countryIdxInTZ
}

func (b *Generator) GeneratePlayerFullName(playerId *big.Int, generation uint8, timezone uint8, countryIdxInTZ uint64) (string, uint, error) {
	log.Debugf("[NAMES] GeneratePlayerFullName of playerId %v", playerId)
	// country_id is an encoding of (tz, countryIdx):
	country_id := encodeCountry(timezone, countryIdxInTZ)
	// if the country is not defined, we use a default country: Spain, at tz = 19
	isSpecified, err := b.isCountrySpecified(country_id)
	if err != nil {
		return "", 0, err
	}
	// Spain is the default country if you query for one that is not specified
	if !isSpecified {
		country_id = encodeCountry(19, 0)
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
		return "", 0, err
	}
	var code_name uint
	var code_surname uint
	var pure_pure int
	var pure_foreign int
	var foreign_pure int
	var foreign_foreign int
	defer rows.Close()
	if !rows.Next() {
		return "", 0, errors.New("Cannot find specs for country_id = %s" + strconv.FormatInt(int64(country_id), 10))
	}
	err = rows.Scan(&code_name, &code_surname, &pure_pure, &pure_foreign, &foreign_pure, &foreign_foreign)
	if err != nil {
		return "", 0, err
	}
	name, _, err := b.GenerateName(false, playerId, generation, code_name, pure_pure+pure_foreign)
	if err != nil {
		return "", 0, err
	}
	
	surname, final_country_code, err := b.GenerateName(true, playerId, generation, code_surname, pure_pure+foreign_pure)
	if err != nil {
		return "", 0, err
	}
	isActualSon := generation > 0 && generation < 32
	if isActualSon {
		surname += " Jr."
	}
	return name + " " + surname, final_country_code, nil
}

func (b *Generator) GenerateTeamName(teamId *big.Int, timezone uint8, countryIdxInTZ uint64) (string, error) {
	// For the time being, we don't use the country information. At some point, we may.
	salt := teamId.String() + "ff"
	// MAIN NAME
	tableName := "team_mainnames"
	nameIdx := b.GenerateRnd(teamId, salt, uint64(b.nTeamnamesMain))
	rows, err := b.db.Query(`SELECT name FROM `+tableName+` WHERE (idx = $1)`, nameIdx)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if !rows.Next() {
		var str string = "Rnd choice failed, teamId = " + teamId.String() +
			", nameIdx = " + strconv.FormatInt(int64(nameIdx), 10) +
			", tableName = " + tableName
		return "", errors.New(str)
	}
	var name string
	rows.Scan(&name)

	// return if has space in name
	if strings.Contains(name, " ") {
		return name, nil
	}

	// ADD PREFFIX OR SUFFIX
	salt += "gg"
	dice := b.GenerateRnd(teamId, salt, uint64(b.nTeamnamesPreffix+b.nTeamnamesSuffix))
	var nNames uint
	addPrefix := uint(dice) < b.nTeamnamesPreffix
	if addPrefix {
		tableName = "team_prefixnames"
		nNames = b.nTeamnamesPreffix
	} else {
		tableName = "team_suffixnames"
		nNames = b.nTeamnamesSuffix
	}
	salt += "hh"
	nameIdx = b.GenerateRnd(teamId, salt, uint64(nNames))
	rows, err = b.db.Query(`SELECT name FROM `+tableName+` WHERE (idx = $1)`, nameIdx)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if !rows.Next() {
		var str string = "Rnd choice failed, teamId = " + teamId.String() +
			", nameIdx = " + strconv.FormatInt(int64(nameIdx), 10) +
			", tableName = " + tableName
		return "", errors.New(str)
	}
	var extraname string
	rows.Scan(&extraname)

	name = strings.Title(strings.ToLower(name))
	extraname = strings.Title(strings.ToLower(extraname))
	if addPrefix {
		return extraname + " " + name, nil
	} else {
		return name + " " + extraname, nil
	}
}
