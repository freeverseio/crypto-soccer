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
		concat += strconv.Itoa(int(computedEvents[i].ManagesToShoot))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].IsGoal))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].PrimaryPlayer))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].SecondaryPlayer))
		concat += "]"
	}
	expected := "[1, 0, 0, 1, 1, 10, 2][4, 0, 1, 0, 0, 7, -1][8, 0, 1, 1, 0, 8, 0][9, 0, 0, 0, 0, 9, -1][14, 0, 0, 0, 0, 4, -1][15, 0, 0, 0, 0, 3, -1][19, 0, 0, 0, 0, 6, -1][23, 0, 0, 0, 0, 5, -1][26, 0, 0, 0, 0, 9, -1][27, 0, 0, 0, 0, 8, -1][30, 0, 0, 0, 0, 7, -1][35, 0, 0, 0, 0, 8, -1][15, 2, 0, -1, -1, 12, -1][9, 1, 0, -1, -1, 4, -1][15, 2, 1, -1, -1, 12, -1][9, 1, 1, -1, -1, 4, -1][14, 5, 0, -1, -1, 5, 19][14, 5, 0, -1, -1, 1, 12][14, 5, 1, -1, -1, 5, 19][14, 5, 1, -1, -1, 1, 12]"
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
		concat += strconv.Itoa(int(computedEvents[i].ManagesToShoot))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].IsGoal))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].PrimaryPlayer))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].SecondaryPlayer))
		concat += "]"
	}
	expected := "[46, 0, 0, 1, 1, 10, 2][49, 0, 1, 0, 0, 7, -1][53, 0, 1, 1, 0, 8, 0][54, 0, 0, 0, 0, 9, -1][59, 0, 0, 0, 0, 4, -1][60, 0, 0, 0, 0, 3, -1][64, 0, 0, 0, 0, 6, -1][68, 0, 0, 0, 0, 5, -1][71, 0, 0, 0, 0, 9, -1][72, 0, 0, 0, 0, 8, -1][75, 0, 0, 0, 0, 7, -1][80, 0, 0, 0, 0, 8, -1][60, 2, 0, -1, -1, 12, -1][54, 1, 0, -1, -1, 4, -1][60, 2, 1, -1, -1, 12, -1][54, 1, 1, -1, -1, 4, -1][59, 5, 0, -1, -1, 5, 19][59, 5, 0, -1, -1, 1, 12][59, 5, 1, -1, -1, 5, 19][59, 5, 1, -1, -1, 1, 12]"
	if concat != expected {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}

}