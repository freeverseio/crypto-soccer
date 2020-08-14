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

func TestHash(t *testing.T) {
	big := big.NewInt(123456789)
	unsig := uint(123456789)
	inputs := []string{"hola", big.String(), strconv.FormatUint(uint64(unsig), 10)}
	expectedOutputs := []uint64{4623503348185510199, 492395637191921148, 492395637191921148}
	for i := 0; i < len(inputs); i++ {
		hash := int_hash(inputs[i])
		if hash != expectedOutputs[i] {
			fmt.Println("Processing Hash with Input and Result:")
			fmt.Println(inputs[i])
			fmt.Println(hash)
			t.Fatal("Wrong hash")
		}
	}
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
	lineup0 := [14]uint8{17, 16, 15, 14, 13, 11, 9, 8, 7, 0, 10, 19, 12, 21}
	lineup1 := [14]uint8{3, 4, 5, 6, 0, 1, 2, 14, 8, 0, 10, 17, 18, 19}
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
		lineup0,
		lineup1,
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
	expected := "[1, 0, 0, true, true, 10, 15][7, 0, 1, false, false, 7, -1][10, 0, 1, true, false, 8, 3][13, 0, 0, false, false, 4, -1][16, 0, 0, false, false, 1, -1][23, 0, 0, false, false, 0, -1][26, 0, 0, false, false, 14, -1][29, 0, 0, false, false, 2, -1][32, 0, 0, false, false, 8, -1][39, 0, 0, false, false, 14, -1][41, 0, 0, false, false, 8, -1][46, 0, 0, false, false, 0, -1][23, 2, 0, false, false, 14, -1][19, 1, 0, false, false, 14, -1][23, 1, 0, false, false, 14, -1][23, 2, 1, false, false, 6, -1][6, 1, 1, false, false, 6, -1][23, 1, 1, false, false, 6, -1][16, 5, 0, false, false, 11, 19][26, 5, 0, false, false, 16, 12][16, 5, 1, false, false, 1, 17][26, 5, 1, false, false, 4, 18]"
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
	}
	var NROUNDS = 12
	for i := 0; i < 5*(NROUNDS-3); i++ {
		events64 = append(events64, 0)
	}
	for i := 0; i < len(events64); i++ {
		events = append(events, big.NewInt(events64[i]))
	}

	NO_SUBS := uint8(11)
	lineup0 := [14]uint8{17, 16, 15, 14, 13, 11, 9, 8, 7, 0, 10, 19, 12, 21}
	lineup1 := [14]uint8{3, 4, 5, 6, 0, 1, 2, 14, 8, 0, 10, 17, 18, 19}
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
		lineup0,
		lineup1,
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
	lineup0 := [14]uint8{17, 16, 15, 14, 13, 11, 9, 8, 7, 0, 10, 19, 12, 21}
	lineup1 := [14]uint8{3, 4, 5, 6, 0, 1, 2, 14, 8, 0, 10, 17, 18, 19}
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
		lineup0,
		lineup1,
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
	expected := "[1, 0, 0, true, true, 10, 15][7, 0, 1, false, false, 7, -1][10, 0, 1, true, false, 8, 3][13, 0, 0, false, false, 4, -1][16, 0, 0, false, false, 1, -1][23, 0, 0, false, false, 0, -1][26, 0, 0, false, false, 14, -1][29, 0, 0, false, false, 2, -1][32, 0, 0, false, false, 8, -1][39, 0, 0, false, false, 14, -1][41, 0, 0, false, false, 8, -1][46, 0, 0, false, false, 0, -1][23, 2, 0, false, false, 12, -1][19, 1, 0, false, false, 13, -1][23, 2, 1, false, false, 18, -1][37, 1, 1, false, false, 0, -1][16, 5, 0, false, false, 11, 19][22, 5, 0, false, false, 16, 12][16, 5, 1, false, false, 1, 17][22, 5, 1, false, false, 4, 18]"
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
	lineup0 := [14]uint8{17, 16, 15, 14, 13, 11, 9, 8, 7, 0, 10, 19, 12, 21}
	lineup1 := [14]uint8{3, 4, 5, 6, 0, 1, 2, 14, 8, 0, 10, 17, 18, 19}
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
		lineup0,
		lineup1,
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
	expected := "[46, 0, 0, true, true, 10, 15][52, 0, 1, false, false, 7, -1][55, 0, 1, true, false, 8, 3][58, 0, 0, false, false, 4, -1][61, 0, 0, false, false, 1, -1][68, 0, 0, false, false, 0, -1][71, 0, 0, false, false, 14, -1][74, 0, 0, false, false, 2, -1][77, 0, 0, false, false, 8, -1][84, 0, 0, false, false, 14, -1][86, 0, 0, false, false, 8, -1][91, 0, 0, false, false, 0, -1][68, 2, 0, false, false, 12, -1][64, 1, 0, false, false, 13, -1][68, 2, 1, false, false, 18, -1][82, 1, 1, false, false, 0, -1][61, 5, 0, false, false, 11, 19][67, 5, 0, false, false, 16, 12][61, 5, 1, false, false, 1, 17][67, 5, 1, false, false, 4, 18]"
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
	// this decodedMatchLog comes from: '452312848584470512245079946786433186608365459112320500501947696564481818624'
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
	lineup0 := [14]uint8{17, 16, 15, 14, 13, 11, 9, 8, 7, 0, 10, 19, 12, 21}
	lineup1 := [14]uint8{3, 4, 5, 6, 0, 1, 2, 14, 8, 0, 10, 17, 18, 19}
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
		lineup0,
		lineup1,
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
	expected := "[46, 0, 1, false, false, 0, -1][52, 0, 1, false, false, 7, -1][55, 0, 0, true, true, 8, 8][58, 0, 1, false, false, 16, -1][61, 0, 0, false, false, 1, -1][68, 0, 1, false, false, 13, -1][71, 0, 1, false, false, 8, -1][74, 0, 0, true, true, 10, 10][77, 0, 0, true, true, 8, 8][84, 0, 0, false, false, 14, -1][86, 0, 1, false, false, 7, -1][91, 0, 0, false, false, 0, -1][68, 2, 0, false, false, 12, -1][64, 1, 0, false, false, 13, -1][68, 2, 1, false, false, 18, -1][82, 1, 1, false, false, 0, -1][61, 5, 0, false, false, 11, 19][67, 5, 0, false, false, 16, 12][61, 5, 1, false, false, 1, 17][67, 5, 1, false, false, 4, 18]"
	allOK := (concat == expected) && (nGoals[0] == 3) && (nGoals[1] == 0)
	if !allOK {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}

}

func TestMatchEventsWithInjuredGKEndOfMatch(t *testing.T) {
	// rounds go from 0...11, so round = 12 is reserved for endOfMatch (e.g. injuries for GKs)
	NROUNDS := uint32(12)
	verseSeed := [32]byte{0x2, 0x1}
	teamId0 := "1"
	teamId1 := "2"
	matchLog := [15]uint32{
		0,             //teamSumSkills,
		0,             //winner,
		0,             //nGoals,
		0,             //trainingPoints1stHalf = 0,
		0, 2, NROUNDS, //outOfGames[0], typesOutOfGames[0], outOfGameRounds[0],  // = hardInjury
		4, 14, //yellowCards1[0], yellowCards1[1],
		1, 1, 5, //ingameSubs1[0], ingameSubs1[1], ingameSubs1[2],
		0, 0, 0} // halftimesubs: 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
	var events []*big.Int
	events64 := []int64{
		4324,  // seed
		34234, // starttime
	}
	for i := uint32(0); i < 5*(NROUNDS); i++ {
		events64 = append(events64, 0)
	}
	for i := 0; i < len(events64); i++ {
		events = append(events, big.NewInt(events64[i]))
	}

	NO_SUBS := uint8(11)
	lineup0 := [14]uint8{17, 16, 15, 14, 13, 11, 9, 8, 7, 0, 10, 19, 12, 21}
	lineup1 := [14]uint8{3, 4, 5, 6, 0, 1, 2, 14, 8, 0, 10, 17, 18, 19}
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
		lineup0,
		lineup1,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
	assert.NilError(t, err)
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
	expected := "[1, 0, 0, false, false, 0, -1][7, 0, 0, false, false, 8, -1][10, 0, 0, false, false, 5, -1][13, 0, 0, false, false, 4, -1][16, 0, 0, false, false, 1, -1][23, 0, 0, false, false, 0, -1][26, 0, 0, false, false, 14, -1][29, 0, 0, false, false, 2, -1][32, 0, 0, false, false, 8, -1][39, 0, 0, false, false, 14, -1][41, 0, 0, false, false, 8, -1][46, 0, 0, false, false, 0, -1][47, 4, 0, false, false, 17, -1][19, 1, 0, false, false, 13, -1][47, 4, 1, false, false, 3, -1][37, 1, 1, false, false, 0, -1][16, 5, 0, false, false, 11, 19][26, 5, 0, false, false, 16, 12][16, 5, 1, false, false, 1, 17][26, 5, 1, false, false, 4, 18]"
	allOK := (concat == expected) && (nGoals[0] == 0) && (nGoals[1] == 0)
	if !allOK {
		fmt.Println("the obtained result is: ")
		fmt.Println(concat)
		t.Fatal("result of generating matchevents not as expected")
	}
}

func TestMatchEventsBadOutOfGame(t *testing.T) {
	// this tests the range of the entires 4-5-6 of matchLog:
	// matchLog[4,5,6] = outOfGamePlayer, outOfGameType, outOfGameRound
	// First, make a choice that works OK, as in the other tests.
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
	lineup0 := [14]uint8{17, 16, 15, 14, 13, 11, 9, 8, 7, 0, 10, 19, 12, 21}
	lineup1 := [14]uint8{3, 4, 5, 6, 0, 1, 2, 14, 8, 0, 10, 17, 18, 19}
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
		lineup0,
		lineup1,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	// Show that typeOfOutOfGame must be 0,1,2, or 3. Fails otherwise
	badLog := matchLog
	badLog[5] = 4
	_, err = matchevents.Generate(
		verseSeed,
		teamId0,
		teamId1,
		badLog,
		matchLog,
		events,
		lineup0,
		lineup1,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
	if err == nil {
		t.Fatalf("error: this command should have failed, but it didnt")
	}

	// Show that typeOfOutOfRound must be, at most, 12.
	badLog2 := matchLog
	badLog2[6] = 13
	_, err = matchevents.Generate(
		verseSeed,
		teamId0,
		teamId1,
		badLog2,
		matchLog,
		events,
		lineup0,
		lineup1,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
	if err == nil {
		t.Fatalf("error: this command should have failed, but it didnt")
	}

	// Show that if there is an outOfGame (for which outOfGamePlayer < 14)
	// Then the type cannot be 0.
	badLog4 := matchLog
	badLog4[4] = 5
	badLog4[5] = 0
	_, err = matchevents.Generate(
		verseSeed,
		teamId0,
		teamId1,
		badLog4,
		matchLog,
		events,
		lineup0,
		lineup1,
		substitutions,
		substitutions,
		subsRounds,
		subsRounds,
		is2ndHalf,
	)
	if err == nil {
		t.Fatalf("error: this command should have failed, but it didnt")
	}
}
