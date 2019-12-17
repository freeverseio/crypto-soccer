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
		0, 0, 0, //ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
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
	is2ndHalf := false

	var lineup [14]uint8
	var substitutions [3]uint8
	var subsRounds [3]uint8

	computedEvents, err := matchevents.GenerateMatchEvents(seed, matchLog, events, lineup, substitutions, subsRounds, is2ndHalf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	concat := ""
	nRows := len(computedEvents[0])
	for i := 0; i < len(computedEvents); i++ {
		concat += "["
		for j := 0; j < nRows; j++ {
			concat += strconv.Itoa(int(computedEvents[i][j]))
			if j < nRows-1 {
				concat += ", "
			}
		}
		concat += "]"
	}
	expected := "[1, 0, 1, 1, 10, 2][4, 1, 0, 0, 7, -1][8, 1, 1, 0, 8, 0][9, 0, 0, 0, 9, -1][14, 0, 0, 0, 4, -1][15, 0, 0, 0, 3, -1][19, 0, 0, 0, 6, -1][23, 0, 0, 0, 5, -1][26, 0, 0, 0, 9, -1][27, 0, 0, 0, 8, -1][30, 0, 0, 0, 7, -1][35, 0, 0, 0, 8, -1][15, 3, -1, -1, 3, -1][9, 2, -1, -1, 4, -1]"
	if concat != expected {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}

}
