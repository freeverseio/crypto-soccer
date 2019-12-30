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

func TestGeneratePlayerName(t *testing.T) {
	generator, err := names.New("./sql/names.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	var result string = ""
	generation := uint8(0)
	// supported (tz, countriesIdxInTz):
	// 		(19, 0): "Spain"
	// 		(16, 0): "China"
	// 		(15, 0): "Japan"
	supported := []uint8{
		19, 0,
		19, 1,
		16, 0,
		18, 0,
		15, 0,
	}
	for place := 0; place < len(supported)/2; place++ {
		for i := 0; i < 10; i++ {
			playerId := big.NewInt(int64(place*1000 + i))
			timezone = supported[place*2]
			countryIdxInTZ = uint64(supported[place*2+1])
			name, final_country_code, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
			if err != nil {
				t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
			}
			fmt.Println(name)
			if len(name) == 0 {
				t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
			}
			result += name
			look, err := generator.GeneratePlayerLook(playerId, generation, final_country_code)
			fmt.Println(look)
		}
		fmt.Println("")
	}

	if int_hash(result) != uint64(3654913364252892805) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}

// func TestGeneratePlayerNameUndefinedCountry(t *testing.T) {
// 	generator, err := names.New("./sql/names.db")
// 	if err != nil {
// 		t.Fatalf("error creating database for player names: %s", err)
// 	}
// 	var timezone uint8
// 	var countryIdxInTZ uint64
// 	var result string = ""
// 	generation := uint8(0)
// 	fmt.Println("TestGeneratePlayerNameUndefinedCountry:")
// 	for i := 0; i < 10; i++ {
// 		playerId := big.NewInt(int64(i))
// 		timezone = uint8(1 + i)
// 		countryIdxInTZ = uint64(3*i + 2)
// 		name, _, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
// 		if err != nil {
// 			t.Fatalf("error generating name for player: %s", playerId)
// 		}
// 		fmt.Println(name)
// 		if len(name) == 0 {
// 			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
// 		}
// 		result += name
// 	}
// 	if int_hash(result) != uint64(14101129100528436475) {
// 		fmt.Println("the just-obtained hash is: ")
// 		fmt.Println(int_hash(result))
// 		t.Fatal("result of generating names not as expected")
// 	}
// }

// func TestGenerateChildName(t *testing.T) {
// 	generator, err := names.New("./sql/names.db")
// 	if err != nil {
// 		t.Fatalf("error creating database for player names: %s", err)
// 	}
// 	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
// 	var timezone uint8
// 	var countryIdxInTZ uint64
// 	// if generation < 32, it means that it is an actual son
// 	generation := uint8(0)
// 	playerId := big.NewInt(int64(1))
// 	timezone = 19
// 	countryIdxInTZ = 0
// 	name, _, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
// 	if err != nil {
// 		t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
// 	}
// 	fmt.Println(name)
// 	if len(name) == 0 {
// 		t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
// 	}
// 	generation = uint8(1)
// 	name2, _, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
// 	if err != nil {
// 		t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
// 	}
// 	fmt.Println(name2)
// 	if len(name2) == 0 {
// 		t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
// 	}
// 	name = name + " " + name2
// 	if name != "Volratas Ramsundar Carles Ramsundar Jr." {
// 		fmt.Println("the just-obtained hash is: ")
// 		fmt.Println(int_hash(name))
// 		fmt.Println(name)
// 		t.Fatal("result of generating names not as expected")
// 	}
// }

// func TestGenerateAcademyName(t *testing.T) {
// 	generator, err := names.New("./sql/names.db")
// 	if err != nil {
// 		t.Fatalf("error creating database for player names: %s", err)
// 	}
// 	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
// 	var timezone uint8
// 	var countryIdxInTZ uint64
// 	// if generation < 32, it means that it is an actual son
// 	generation := uint8(33)
// 	playerId := big.NewInt(int64(1))
// 	timezone = 19
// 	countryIdxInTZ = 0
// 	name, _, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
// 	if err != nil {
// 		t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
// 	}
// 	fmt.Println(name)
// 	if len(name) == 0 {
// 		t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
// 	}
// 	if name != "Carles Mildner" {
// 		fmt.Println("the just-obtained hash is: ")
// 		fmt.Println(int_hash(name))
// 		fmt.Println(name)
// 		t.Fatal("result of generating names not as expected")
// 	}
// }

// func TestGenerateTeamName(t *testing.T) {
// 	generator, err := names.New("./sql/names.db")
// 	if err != nil {
// 		t.Fatalf("error creating database for player names: %s", err)
// 	}
// 	var timezone uint8
// 	var countryIdxInTZ uint64
// 	timezone = 19
// 	countryIdxInTZ = 0
// 	teamId := big.NewInt(int64(0))
// 	var name string
// 	var concatname string
// 	for i := 0; i < 10; i++ {
// 		teamId = big.NewInt(int64(41234332 + i))
// 		name, err = generator.GenerateTeamName(teamId, timezone, countryIdxInTZ)
// 		if err != nil {
// 			t.Fatalf("error generating name for team %s: %s", teamId.String(), err)
// 		}
// 		fmt.Println(name)
// 		if len(name) == 0 {
// 			t.Fatalf("Expecting non empty team name, but got \"%v\"", name)
// 		}
// 		concatname += " " + name
// 	}
// 	if concatname != " Violet Jellyfish Water Ox Tequila Sunrise Magenta Indian Astronaut Sunrise Hope Mouse Toucan Magic Twins Blues A. Z. Gaia C. A. Harmony" {
// 		fmt.Println("the just-obtained hash is: ")
// 		fmt.Println(int_hash(concatname))
// 		fmt.Println(concatname)
// 		t.Fatal("result of generating names not as expected")
// 	}
// }
