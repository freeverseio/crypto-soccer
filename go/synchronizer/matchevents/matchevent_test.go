package matchevents_test

import (
	"fmt"
	"hash/fnv"
	"math/big"
	"strconv"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/matchevents"
	"gotest.tools/assert"
)

func int_hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func TestMatchEventsTwoYellows(t *testing.T) {
	verseSeed := [32]byte{0x2, 0x1}
	teamId0 := "1"
	teamId1 := "2"
	matchLog := [15]uint32{
		0,       //teamSumSkills,
		0,       //winner,
		0,       //nGoals,
		0,       //trainingPoints1stHalf = 0,
		3, 3, 5, //outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],
		3, 3, //yellowCards1[0], yellowCards1[1],
		1, 1, 0, //ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
		0, 0, 0} // halftimesubs: 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
	var events []*big.Int
	events64 := []int64{
		0,              // log0 (not used by the algorithm)
		0,              // log1 (not used by the algorithm)
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

	computedEvents, err := matchevents.Generate(
		verseSeed,
		teamId0,
		teamId1,
		matchLog,
		matchLog,
		events,
		lineup,
		lineup,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
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
	expected := "[1, 0, 0, true, true, 10, 2][6, 0, 1, false, false, 4, -1][11, 0, 1, true, false, 8, 0][12, 0, 0, false, false, 2, -1][17, 0, 0, false, false, 8, -1][22, 0, 0, false, false, 9, -1][27, 0, 0, false, false, 6, -1][28, 0, 0, false, false, 7, -1][33, 0, 0, false, false, 4, -1][38, 0, 0, false, false, 5, -1][42, 0, 0, false, false, 7, -1][45, 0, 0, false, false, 6, -1][22, 2, 0, false, false, 3, -1][15, 1, 0, false, false, 3, -1][22, 1, 0, false, false, 3, -1][22, 2, 1, false, false, 3, -1][15, 1, 1, false, false, 3, -1][22, 1, 1, false, false, 3, -1][17, 5, 0, false, false, 5, 19][27, 5, 0, false, false, 1, 12][17, 5, 1, false, false, 5, 19][27, 5, 1, false, false, 1, 12]"
	if concat != expected {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}
}

func TestMatchEvent2Events(t *testing.T) {
	verseSeed := [32]byte{0x2, 0x1}
	teamId0 := "1"
	teamId1 := "2"
	matchLog := [15]uint32{
		0,       //teamSumSkills,
		0,       //winner,
		0,       //nGoals,
		0,       //trainingPoints1stHalf = 0,
		3, 3, 5, //outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],
		3, 3, //yellowCards1[0], yellowCards1[1],
		1, 1, 0, //ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
		0, 0, 0} // halftimesubs: 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)

	events := []*big.Int{}
	events[0] = big.NewInt(0)
	events[1] = big.NewInt(0)
	assert.Equal(t, len(events), 2)

	NO_SUBS := uint8(11)
	lineup := [14]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 19, 12, 21}
	substitutions := [3]uint8{5, 1, NO_SUBS}
	subsRounds := [3]uint8{4, 6, 7}

	is2ndHalf := false

	_, err := matchevents.Generate(
		verseSeed,
		teamId0,
		teamId1,
		matchLog,
		matchLog,
		events,
		lineup,
		lineup,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
	assert.NilError(t, err)
}

func TestMatchEvents(t *testing.T) {
	verseSeed := [32]byte{0x2, 0x1}
	teamId0 := "1"
	teamId1 := "2"
	matchLog := [15]uint32{
		0,        //teamSumSkills,
		0,        //winner,
		0,        //nGoals,
		0,        //trainingPoints1stHalf = 0,
		12, 3, 5, //outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],
		4, 14, //yellowCards1[0], yellowCards1[1],
		1, 1, 0, //ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
		0, 0, 0} // halftimesubs: 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
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

	computedEvents, err := matchevents.Generate(
		verseSeed,
		teamId0,
		teamId1,
		matchLog,
		matchLog,
		events,
		lineup,
		lineup,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
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
	expected := "[1, 0, 0, true, true, 10, 2][6, 0, 1, false, false, 4, -1][11, 0, 1, true, false, 8, 0][12, 0, 0, false, false, 2, -1][17, 0, 0, false, false, 8, -1][22, 0, 0, false, false, 9, -1][27, 0, 0, false, false, 6, -1][28, 0, 0, false, false, 7, -1][33, 0, 0, false, false, 4, -1][38, 0, 0, false, false, 5, -1][42, 0, 0, false, false, 7, -1][45, 0, 0, false, false, 6, -1][22, 2, 0, false, false, 12, -1][30, 1, 0, false, false, 4, -1][22, 2, 1, false, false, 12, -1][30, 1, 1, false, false, 4, -1][17, 5, 0, false, false, 5, 19][21, 5, 0, false, false, 1, 12][17, 5, 1, false, false, 5, 19][21, 5, 1, false, false, 1, 12]"
	if concat != expected {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}

}

func TestMatchEvents2ndHalf(t *testing.T) {
	verseSeed := [32]byte{0x2, 0x1}
	teamId0 := "1"
	teamId1 := "2"
	matchLog := [15]uint32{
		0,        //teamSumSkills,
		0,        //winner,
		0,        //nGoals,
		0,        //trainingPoints1stHalf = 0,
		12, 3, 5, //outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],
		4, 14, //yellowCards1[0], yellowCards1[1],
		1, 1, 0, //ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
		0, 0, 0} // halftimesubs: 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
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
	computedEvents, err := matchevents.Generate(
		verseSeed,
		teamId0,
		teamId1,
		matchLog,
		matchLog,
		events,
		lineup,
		lineup,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
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
	expected := "[46, 0, 0, true, true, 10, 2][51, 0, 1, false, false, 4, -1][56, 0, 1, true, false, 8, 0][57, 0, 0, false, false, 2, -1][62, 0, 0, false, false, 8, -1][67, 0, 0, false, false, 9, -1][72, 0, 0, false, false, 6, -1][73, 0, 0, false, false, 7, -1][78, 0, 0, false, false, 4, -1][83, 0, 0, false, false, 5, -1][87, 0, 0, false, false, 7, -1][90, 0, 0, false, false, 6, -1][67, 2, 0, false, false, 12, -1][75, 1, 0, false, false, 4, -1][67, 2, 1, false, false, 12, -1][75, 1, 1, false, false, 4, -1][62, 5, 0, false, false, 5, 19][66, 5, 0, false, false, 1, 12][62, 5, 1, false, false, 5, 19][66, 5, 1, false, false, 1, 12]"
	if concat != expected {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}

}

func TestMatchEvents2ndHalfHardcoded(t *testing.T) {
	// in this test, events64 is hardcoded (coming from a set of events that once gave apparently wrong final results)
	// so we test that team0 indeed scores 3 goals, given the hardcoded events64
	verseSeed := [32]byte{0x2, 0x1}
	teamId0 := "1"
	teamId1 := "2"
	matchLog := [15]uint32{
		0,        //teamSumSkills,
		0,        //winner,
		0,        //nGoals,
		0,        //trainingPoints1stHalf = 0,
		12, 3, 5, //outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],
		4, 14, //yellowCards1[0], yellowCards1[1],
		1, 1, 0, //ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
		0, 0, 0} // halftimesubs: 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
	var events []*big.Int
	events64 := []int64{
		4324,  // seed
		34234, // starttime
		1, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		0, 1, 7, 1, 7,
		1, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		0, 1, 10, 1, 10,
		0, 1, 7, 1, 7,
		0, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}

	for i := 0; i < len(events64); i++ {
		events = append(events, big.NewInt(events64[i]))
	}

	NO_SUBS := uint8(11)
	lineup := [14]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 19, 12, 21}
	substitutions := [3]uint8{5, 1, NO_SUBS}
	subsRounds := [3]uint8{4, 6, 7}

	is2ndHalf := true
	computedEvents, err := matchevents.Generate(
		verseSeed,
		teamId0,
		teamId1,
		matchLog,
		matchLog,
		events,
		lineup,
		lineup,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	concat := ""
	nGoals := [2]uint8{0, 0}
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
		if computedEvents[i].IsGoal {
			nGoals[computedEvents[i].Team]++
		}
		concat += strconv.Itoa(int(computedEvents[i].PrimaryPlayer))
		concat += ", "
		concat += strconv.Itoa(int(computedEvents[i].SecondaryPlayer))
		concat += "]"
	}
	expected := "[46, 0, 1, false, false, 3, -1][51, 0, 1, false, false, 4, -1][56, 0, 0, true, true, 7, 7][57, 0, 1, false, false, 2, -1][62, 0, 0, false, false, 8, -1][67, 0, 1, false, false, 9, -1][72, 0, 1, false, false, 6, -1][73, 0, 0, true, true, 10, 10][78, 0, 0, true, true, 7, 7][83, 0, 0, false, false, 5, -1][87, 0, 1, false, false, 7, -1][90, 0, 0, false, false, 6, -1][67, 2, 0, false, false, 12, -1][75, 1, 0, false, false, 4, -1][67, 2, 1, false, false, 12, -1][75, 1, 1, false, false, 4, -1][62, 5, 0, false, false, 5, 19][66, 5, 0, false, false, 1, 12][62, 5, 1, false, false, 5, 19][66, 5, 1, false, false, 1, 12]"
	allOK := (concat == expected) && (nGoals[0] == 3) && (nGoals[1] == 0)
	if !allOK {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}

}
