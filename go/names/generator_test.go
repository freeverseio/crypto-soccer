package names_test

import (
	"fmt"
	"hash/fnv"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/names"
)

func int_hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func TestGeneratePlayerNamePoland(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	var result string = ""
	generation := uint8(0)
	for i := 0; i < 100; i++ {
		playerId := big.NewInt(int64(i))
		timezone = 8
		countryIdxInTZ = 0
		name, countryISO2, region, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
		}
		fmt.Println(name + " (" + countryISO2 + ", " + region + ") ")
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
		result += name + " (" + countryISO2 + ", " + region + ") "
	}
	if int_hash(result) != uint64(12191899905785622891) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}

func TestGeneratePlayerNameItaly(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	var result string = ""
	generation := uint8(0)
	for i := 0; i < 100; i++ {
		playerId := big.NewInt(int64(i))
		timezone = 11
		countryIdxInTZ = 0
		name, countryISO2, region, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
		}
		fmt.Println(name + " (" + countryISO2 + ", " + region + ") ")
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
		result += name + " (" + countryISO2 + ", " + region + ") "
	}
	if int_hash(result) != uint64(6009804721545048910) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}
func TestGeneratePlayerNameBelgium(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	var result string = ""
	generation := uint8(0)
	for i := 0; i < 100; i++ {
		playerId := big.NewInt(int64(i))
		timezone = 9
		countryIdxInTZ = 1
		name, countryISO2, region, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
		}
		fmt.Println(name + " (" + countryISO2 + ", " + region + ") ")
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
		result += name + " (" + countryISO2 + ", " + region + ") "
	}
	if int_hash(result) != uint64(15455280794236233387) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}

func TestGeneratePlayerNameNetherlands(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	var result string = ""
	generation := uint8(0)
	for i := 0; i < 100; i++ {
		playerId := big.NewInt(int64(i))
		timezone = 9
		countryIdxInTZ = 0
		name, countryISO2, region, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
		}
		fmt.Println(name + " (" + countryISO2 + ", " + region + ") ")
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
		result += name + " (" + countryISO2 + ", " + region + ") "
	}
	if int_hash(result) != uint64(8719546611057531656) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}

func TestGeneratePlayerNameSpain(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	var result string = ""
	generation := uint8(0)
	for i := 0; i < 10; i++ {
		playerId := big.NewInt(int64(i))
		timezone = 10
		countryIdxInTZ = 0
		name, countryISO2, region, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
		}
		fmt.Println(name + " (" + countryISO2 + ", " + region + ") ")
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
		result += name + " (" + countryISO2 + ", " + region + ") "
	}
	if int_hash(result) != uint64(12560545947695793294) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}

func TestGeneratePlayerNameUndefinedCountry(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	var timezone uint8
	var countryIdxInTZ uint64
	var result string = ""
	generation := uint8(0)
	for i := 0; i < 10; i++ {
		playerId := big.NewInt(int64(i))
		timezone = uint8(1 + i)
		countryIdxInTZ = uint64(3*i + 2)
		name, countryISO2, region, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for player: %s", playerId)
		}
		fmt.Println(name)
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
		result += name + " (" + countryISO2 + ", " + region + ") "
	}
	if int_hash(result) != uint64(12560545947695793294) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}

func TestGenerateChildName(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	// if generation < 32, it means that it is an actual son
	generation := uint8(0)
	playerId := big.NewInt(int64(1))
	timezone = 19
	countryIdxInTZ = 0
	name, _, _, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
	if err != nil {
		t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
	}
	fmt.Println(name)
	if len(name) == 0 {
		t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
	}
	generation = uint8(1)
	name2, _, _, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
	if err != nil {
		t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
	}
	fmt.Println(name2)
	if len(name2) == 0 {
		t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
	}
	name = name + " " + name2
	if name != "Bernat Kebede Anartz Kebede Jr." {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(name))
		fmt.Println(name)
		t.Fatal("result of generating names not as expected")
	}
}

func TestGenerateAcademyName(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	// if generation < 32, it means that it is an actual son
	generation := uint8(33)
	playerId := big.NewInt(int64(1))
	timezone = 19
	countryIdxInTZ = 0
	name, _, _, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
	if err != nil {
		t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
	}
	fmt.Println(name)
	if len(name) == 0 {
		t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
	}
	if name != "Rafel Mccornell" {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(name))
		fmt.Println(name)
		t.Fatal("result of generating names not as expected")
	}
}

func TestGenerateTeamName(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	var timezone uint8
	var countryIdxInTZ uint64
	timezone = 19
	countryIdxInTZ = 0
	teamId := big.NewInt(int64(0))
	var name string
	var concatname string
	for i := 0; i < 10; i++ {
		teamId = big.NewInt(int64(41234332 + i))
		name, err = generator.GenerateTeamName(teamId, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for team %s: %s", teamId.String(), err)
		}
		fmt.Println(name)
		if len(name) == 0 {
			t.Fatalf("Expecting non empty team name, but got \"%v\"", name)
		}
		concatname += " " + name
	}
	if concatname != " Violet Jellyfish Water Ox Tequila Sunrise Magenta Indian Astronaut Sunrise Hope Mouse Toucan Magic Twins Blues A. Z. Gaia C. A. Harmony" {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(concatname))
		fmt.Println(concatname)
		t.Fatal("result of generating names not as expected")
	}
}

func TestGeneratePlayerNameForWrongInputs(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	var countryIdxInTZ uint64
	generation := uint8(0)
	playerId := big.NewInt(int64(41234523))
	countryIdxInTZ = 0
	for tz := uint8(0); tz < 30; tz++ {
		_, _, _, err := generator.GeneratePlayerFullName(playerId, generation, tz, countryIdxInTZ)
		shouldFail := !(tz > 0 && tz < 25)
		if shouldFail && err == nil {
			t.Fatalf("test should fail, but it did not, for tz %v: %s", tz, err)
		}
		if !shouldFail && err != nil {
			t.Fatalf("test should be OK, but it failed, for tz %v: %s", tz, err)
		}
	}
	tz := uint8(1)
	for g := uint8(0); g < 70; g++ {
		_, _, _, err := generator.GeneratePlayerFullName(playerId, g, tz, countryIdxInTZ)
		shouldFail := !(g < 64)
		if shouldFail && err == nil {
			t.Fatalf("test should fail, but it did not, for generation %v: %s", g, err)
		}
		if !shouldFail && err != nil {
			t.Fatalf("test should be OK, but it failed, for generation %v: %s", g, err)
		}
	}
}
