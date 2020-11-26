package names

import (
	"database/sql"
	"fmt"
	"hash/fnv"
	"math/big"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Generator struct {
	db                     *sql.DB
	nonEmptyCountries      []string
	nonEmptyRegions        []string
	nNamesPerCountry       map[string]uint
	nSurnamesPerRegion     map[string]uint
	deployedCountriesSpecs map[uint64]DeployedCountriesSpecs
	nTeamnamesMain         uint
	nTeamnamesPreffix      uint
	nTeamnamesSuffix       uint
}

type DeployedCountriesSpecs struct {
	iso2          string
	region        string
	namePurity    uint
	surnamePurity uint
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
	if err := generator.countEntriesPerArea(false, "iso2", "countries", "names"); err != nil {
		return nil, err
	}

	if err := generator.countEntriesPerArea(true, "region", "regions", "surnames"); err != nil {
		return nil, err
	}

	if err := generator.countTeamsDB(); err != nil {
		return nil, err
	}
	generator.readDeployedCountriesSpecs()
	return &generator, nil
}

func (b *Generator) Close() error {
	return b.db.Close()
}

func serializeTZandCountryIdx(tz uint8, countryIdxInTZ uint64) uint64 {
	return uint64(tz)*1000000 + countryIdxInTZ
}

func (b *Generator) readDeployedCountriesSpecs() {
	m := make(map[uint64]DeployedCountriesSpecs)
	m[serializeTZandCountryIdx(uint8(10), uint64(0))] = DeployedCountriesSpecs{"ES", "Spanish", 75, 60}
	m[serializeTZandCountryIdx(uint8(11), uint64(0))] = DeployedCountriesSpecs{"IT", "ItalySurnames", 75, 60}
	m[serializeTZandCountryIdx(uint8(7), uint64(0))] = DeployedCountriesSpecs{"CN", "Chinese", 75, 70}
	m[serializeTZandCountryIdx(uint8(9), uint64(0))] = DeployedCountriesSpecs{"NL", "NetherlandsSurnames", 75, 60}
	m[serializeTZandCountryIdx(uint8(9), uint64(1))] = DeployedCountriesSpecs{"BE", "BelgiumSurnames", 75, 60}
	m[serializeTZandCountryIdx(uint8(8), uint64(0))] = DeployedCountriesSpecs{"PL", "PolandSurnames", 75, 60}
	b.deployedCountriesSpecs = m
}

func (b *Generator) countEntriesPerArea(isSurname bool, colName string, areasTable string, entryPerAreaTable string) error {
	var err error
	rows, err := b.db.Query(`SELECT ` + colName + ` FROM ` + areasTable)
	if err != nil {
		return err
	}
	defer rows.Close()
	m := make(map[string]uint)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		rows2, err2 := b.db.Query(`SELECT COUNT (*) FROM `+entryPerAreaTable+` WHERE (`+colName+` = $1)`, id)
		if err2 != nil {
			return err2
		}
		defer rows2.Close()
		rows2.Next()
		count := uint(0)
		err = rows2.Scan(&count)
		if err != nil {
			return err
		}
		m[id] = count
		if count > 0 {
			if isSurname {
				b.nonEmptyRegions = append(b.nonEmptyRegions, id)
			}
			if !isSurname {
				b.nonEmptyCountries = append(b.nonEmptyCountries, id)
			}
		}
	}
	if isSurname {
		b.nSurnamesPerRegion = m
	}
	if !isSurname {
		b.nNamesPerCountry = m
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

func (b *Generator) GenerateName(playerId *big.Int, generation uint8, iso2 string, purity uint) (string, string, error) {
	dice := b.GenerateRnd(playerId, "aa", 100)
	if uint(dice) > purity {
		newPick := b.GenerateRnd(playerId, "bb", uint64(len(b.nonEmptyCountries)))
		iso2 = b.nonEmptyCountries[newPick]
	}
	// Make sure we chose different rnds for each generation
	salt := "cc" + strconv.FormatUint(uint64(generation), 10)
	rowRandom := b.GenerateRnd(playerId, salt, uint64(b.nNamesPerCountry[iso2]))
	// LIMIT m OFFSET n will skip the first n entries and return the next m entries
	// Since rowRandom = [0,..., nEntries-1], the upper limit should work
	rows, err := b.db.Query(`SELECT name FROM names WHERE (iso2 = $1) ORDER BY name ASC LIMIT 1 OFFSET $2`, iso2, rowRandom)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	var name string
	if !rows.Next() {
		return "", "", fmt.Errorf("Name: Rnd choice failed: iso2 = %v, rowRandom = %v; input params: playerId = %v, generation = %v, iso2 = %v, purity = %v",
			iso2,
			rowRandom,
			playerId,
			generation,
			iso2,
			purity,
		)
	}
	err = rows.Scan(&name)
	if err != nil {
		return "", "", err
	}
	return name, iso2, nil
}

func (b *Generator) GenerateSurname(playerId *big.Int, generation uint8, region string, purity uint) (string, string, error) {
	dice := b.GenerateRnd(playerId, "dd", 100)
	if uint(dice) > purity {
		newPick := b.GenerateRnd(playerId, "ee", uint64(len(b.nonEmptyRegions)))
		region = b.nonEmptyRegions[newPick]
	}
	// If player is a son, then no matter which generation, it will have the same surname
	// as the primary father (gen=0). Otherwise, random.
	// So, "son" means descendant of the primary father, not descendant of the previous player.
	// All players with gen > 32 are not sons
	var salt string
	if generation < 32 {
		salt = "ff0"
	} else {
		salt = "ff" + strconv.FormatUint(uint64(generation), 10)
	}
	rowRandom := b.GenerateRnd(playerId, salt, uint64(b.nSurnamesPerRegion[region]))
	// LIMIT m OFFSET n will skip the first n entries and return the next m entries
	// Since rowRandom = [0,..., nEntries-1], the upper limit should work
	rows, err := b.db.Query(`SELECT surname FROM surnames WHERE (region = $1) ORDER BY surname ASC LIMIT 1 OFFSET $2`, region, rowRandom)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	var surname string
	if !rows.Next() {
		return "", "", fmt.Errorf("Surname: Rnd choice failed: region = %v, rowRandom = %v; input params: playerId = %v, generation = %v, region = %v, purity = %v",
			region,
			rowRandom,
			playerId,
			generation,
			region,
			purity,
		)
	}
	err = rows.Scan(&surname)
	if err != nil {
		return "", "", err
	}
	isSon := generation > 0 && generation < 32
	if isSon {
		surname += " Jr."
	}
	return surname, region, nil
}

func (b *Generator) GeneratePlayerFullName(playerId *big.Int, generation uint8, tz uint8, countryIdxInTZ uint64) (string, string, string, error) {
	log.Debugf("[NAMES] GeneratePlayerFullName of playerId %v", playerId)
	if tz == 0 || tz > 24 {
		return "", "", "", fmt.Errorf("Timezone should be within [1, 24], but it was %v", tz)
	}
	if generation >= 64 {
		return "", "", "", fmt.Errorf("Generation should be within [0, 63], but it was %v", generation)
	}
	specs, ok := b.deployedCountriesSpecs[serializeTZandCountryIdx(tz, countryIdxInTZ)]
	if !ok {
		// Spain is the default country if you query for one that is not specified
		specs = b.deployedCountriesSpecs[serializeTZandCountryIdx(10, 0)]
	}
	name, countryISO2, err := b.GenerateName(playerId, generation, specs.iso2, specs.namePurity)
	if err != nil {
		return "", "", "", err
	}
	surname, region, err := b.GenerateSurname(playerId, generation, specs.region, specs.surnamePurity)
	if err != nil {
		return "", "", "", err
	}
	return name + " " + surname, countryISO2, region, nil
}

func (b *Generator) GenerateTeamName(teamId *big.Int, tz uint8, countryIdxInTZ uint64) (string, error) {
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
		return "", fmt.Errorf("Rnd choice failed in GenerateTeamName Part1: teamId = %v, tableName = %v, nameIdx = %v, in function with input params: teamId = %v, tz = %v, countryIdxInTZ = %v",
			teamId,
			tableName,
			nameIdx,
			teamId,
			tz,
			countryIdxInTZ,
		)
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
		return "", fmt.Errorf("Rnd choice failed in GenerateTeamName Part2: teamId = %v, tableName = %v, nameIdx = %v, in function with input params: teamId = %v, tz = %v, countryIdxInTZ = %v",
			teamId,
			tableName,
			nameIdx,
			teamId,
			tz,
			countryIdxInTZ,
		)

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
