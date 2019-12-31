package matchevents_test

import (
	"fmt"
	"hash/fnv"
	"math/big"
	"strconv"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/matchevents"
)

func int_hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func TestMatchEvents(t *testing.T) {
	seed := big.NewInt(12334234543)
	matchLog := [15]uint32{
		0,        //teamSumSkills,
		0,        //winner,
		0,        //nGoals,
		0,        //trainingPoints1stHalf = 0,
		12, 3, 5, //outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],
		4, 14, //yellowCards1[0], yellowCards1[1],
		1, 1, 0, //ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
		0, 0, 0} // halftimesubs
	var events []*big.Int
	events64 := []int64{
		4324,           // seed
		34234,          // starttime
		0, 1, 10, 1, 2, // teamThatAttacks, managesToShoot, shooter, isGoal, assister
		1, 0, 0, 0, 0, // teamThatAttacks, managesToShoot, shooter, isGoal, assister
		1, 1, 8, 0, 0, // teamThatAttacks, managesToShoot, shooter, isGoal, assister
	}
	var NROUNDS = 12
	for i := 0; i < 5*(NROUNDS-3); i++ {
		events64 = append(events64, 0)
	}
	for i := 0; i < len(events64); i++ {
		events = append(events, big.NewInt(events64[i]))
	}

	NO_SUBS := uint8(11)
	lineup := [14]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 19, 12, 21}
	substitutions := [3]uint8{5, 1, NO_SUBS}
	subsRounds := [3]uint8{4, 6, 7}

	is2ndHalf := false
	computedEvents, err := matchevents.GenerateMatchEvents(seed, matchLog, matchLog, events, lineup, lineup, substitutions, substitutions, subsRounds, subsRounds, is2ndHalf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	concat := ""
	for i := 0; i < len(computedEvents); i++ {
		concat += "["
		concat += strconv.Itoa(int(computedEvents[i].Minute))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].Type))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].Team))
		concat += ", "
		concat += strconv.FormatBool(computedEvents[i].ManagesToShoot)
		concat += ", "
		concat += strconv.FormatBool(computedEvents[i].IsGoal)
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].PrimaryPlayer))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].SecondaryPlayer))
		concat += "]"
	}
	expected := "[1, 0, 0, true, true, 10, 2][7, 0, 1, false, false, 7, -1][10, 0, 1, true, false, 8, 0][13, 0, 0, false, false, 9, -1][16, 0, 0, false, false, 4, -1][23, 0, 0, false, false, 3, -1][26, 0, 0, false, false, 6, -1][29, 0, 0, false, false, 5, -1][32, 0, 0, false, false, 9, -1][39, 0, 0, false, false, 8, -1][41, 0, 0, false, false, 7, -1][46, 0, 0, false, false, 8, -1][23, 2, 0, false, false, 12, -1][9, 1, 0, false, false, 4, -1][23, 2, 1, false, false, 12, -1][9, 1, 1, false, false, 4, -1][16, 5, 0, false, false, 5, 19][22, 5, 0, false, false, 1, 12][16, 5, 1, false, false, 5, 19][22, 5, 1, false, false, 1, 12]"
	if concat != expected {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}

}

func TestMatchEvents2ndHalf(t *testing.T) {
	seed := big.NewInt(12334234543)
	matchLog := [15]uint32{
		0,        //teamSumSkills,
		0,        //winner,
		0,        //nGoals,
		0,        //trainingPoints1stHalf = 0,
		12, 3, 5, //outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],
		4, 14, //yellowCards1[0], yellowCards1[1],
		1, 1, 0, //ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
		0, 0, 0} // halftimesubs
	var events []*big.Int
	events64 := []int64{
		4324,           // seed
		34234,          // starttime
		0, 1, 10, 1, 2, // teamThatAttacks, managesToShoot, shooter, isGoal, assister
		1, 0, 0, 0, 0, // teamThatAttacks, managesToShoot, shooter, isGoal, assister
		1, 1, 8, 0, 0, // teamThatAttacks, managesToShoot, shooter, isGoal, assister
	}
	var NROUNDS = 12
	for i := 0; i < 5*(NROUNDS-3); i++ {
		events64 = append(events64, 0)
	}
	for i := 0; i < len(events64); i++ {
		events = append(events, big.NewInt(events64[i]))
	}

	NO_SUBS := uint8(11)
	lineup := [14]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 19, 12, 21}
	substitutions := [3]uint8{5, 1, NO_SUBS}
	subsRounds := [3]uint8{4, 6, 7}

	is2ndHalf := true
	computedEvents, err := matchevents.GenerateMatchEvents(seed, matchLog, matchLog, events, lineup, lineup, substitutions, substitutions, subsRounds, subsRounds, is2ndHalf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	concat := ""
	for i := 0; i < len(computedEvents); i++ {
		concat += "["
		concat += strconv.Itoa(int(computedEvents[i].Minute))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].Type))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].Team))
		concat += ", "
		concat += strconv.FormatBool(computedEvents[i].ManagesToShoot)
		concat += ", "
		concat += strconv.FormatBool(computedEvents[i].IsGoal)
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].PrimaryPlayer))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].SecondaryPlayer))
		concat += "]"
	}
	expected := "[46, 0, 0, true, true, 10, 2][52, 0, 1, false, false, 7, -1][55, 0, 1, true, false, 8, 0][58, 0, 0, false, false, 9, -1][61, 0, 0, false, false, 4, -1][68, 0, 0, false, false, 3, -1][71, 0, 0, false, false, 6, -1][74, 0, 0, false, false, 5, -1][77, 0, 0, false, false, 9, -1][84, 0, 0, false, false, 8, -1][86, 0, 0, false, false, 7, -1][91, 0, 0, false, false, 8, -1][68, 2, 0, false, false, 12, -1][54, 1, 0, false, false, 4, -1][68, 2, 1, false, false, 12, -1][54, 1, 1, false, false, 4, -1][61, 5, 0, false, false, 5, 19][67, 5, 0, false, false, 1, 12][61, 5, 1, false, false, 5, 19][67, 5, 1, false, false, 1, 12]"
	if concat != expected {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}

}
