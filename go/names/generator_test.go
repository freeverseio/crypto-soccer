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
	generator, err := names.New("./sql/names.db", "./sql/teamnames.db")
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
		timezone = 19
		countryIdxInTZ = 0
		name, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
		}
		fmt.Println(name)
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
		result += name
	}
	if int_hash(result) != uint64(11127233765381183977) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}

func TestGeneratePlayerNameUndefinedCountry(t *testing.T) {
	generator, err := names.New("./sql/names.db", "./sql/teamnames.db")
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
		name, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for player: %s", playerId)
		}
		fmt.Println(name)
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
		result += name
	}
	if int_hash(result) != uint64(11127233765381183977) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}

func TestGenerateChildName(t *testing.T) {
	generator, err := names.New("./sql/names.db", "./sql/teamnames.db")
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
	name, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
	if err != nil {
		t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
	}
	fmt.Println(name)
	if len(name) == 0 {
		t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
	}
	generation = uint8(1)
	name2, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
	if err != nil {
		t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
	}
	fmt.Println(name2)
	if len(name2) == 0 {
		t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
	}
	name = name + " " + name2
	if name != "Vicent Jessmer Edvin Jessmer Jr." {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(name))
		fmt.Println(name)
		t.Fatal("result of generating names not as expected")
	}
}

func TestGenerateAcademyName(t *testing.T) {
	generator, err := names.New("./sql/names.db", "./sql/teamnames.db")
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
	name, err := generator.GeneratePlayerFullName(playerId, generation, timezone, countryIdxInTZ)
	if err != nil {
		t.Fatalf("error generating name for player %s: %s", playerId.String(), err)
	}
	fmt.Println(name)
	if len(name) == 0 {
		t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
	}
	if name != "Edvin Blasco" {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(name))
		fmt.Println(name)
		t.Fatal("result of generating names not as expected")
	}
}

func TestGenerateTeamName(t *testing.T) {
	generator, err := names.New("./sql/names.db", "./sql/teamnames.db")
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
	for i := 0; i < 100; i++ {
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
	if concatname != " Scarlet Butterfly A. Z. Tortoise Blue Cyclopes Snow Hurricane F. Z. Paradise C. Z. Firefly Sea Jackal Kallisto Plus Tequila Sunset Pegasos Island" {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(concatname))
		fmt.Println(concatname)
		t.Fatal("result of generating names not as expected")
	}
}
