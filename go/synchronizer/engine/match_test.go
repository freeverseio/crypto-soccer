package engine_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/matchevents"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestDefaultMatch(t *testing.T) {
	t.Parallel()
	match := engine.NewMatch()
	golden.Assert(t, dump.Sdump(match), t.Name()+".golden")
}
func TestDefaultValues(t *testing.T) {
	t.Parallel()
	engine := engine.NewMatch()
	assert.Assert(t, engine != nil, "engine is nil")
}

func TestPlay1stHalfWithEmptyTeam(t *testing.T) {
	t.Parallel()
	match := engine.NewMatch()
	match.Seed = [32]byte{0x2, 0x1}
	match.StartTime = big.NewInt(1570147200)
	match.HomeTeam.TeamID = "1"
	match.VisitorTeam.TeamID = "2"
	for i := 0; i < 25; i++ {
		match.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		match.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	err := match.Play1stHalf(*bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, dump.Sdump(match), t.Name()+".golden")
	assert.Equal(t, match.HomeGoals, uint8(0))
	assert.Equal(t, match.VisitorGoals, uint8(0))
	assert.Equal(t, match.HomeTeamSumSkills, uint32(0))
	assert.Equal(t, match.VisitorTeamSumSkills, uint32(0))
}

func TestPlay1stHalfConsumeTheTrainingPoints(t *testing.T) {
	t.Parallel()
	match := engine.NewMatch()
	match.Seed = [32]byte{0x2, 0x1}
	match.StartTime = big.NewInt(1570147200)
	match.HomeTeam.TeamID = "1"
	match.VisitorTeam.TeamID = "2"
	assert.NilError(t, match.Play1stHalf(*bc.Contracts))
	assert.Equal(t, match.HomeTeam.TrainingPoints, uint16(0))
	assert.Equal(t, match.VisitorTeam.TrainingPoints, uint16(0))
	assert.NilError(t, match.Play2ndHalf(*bc.Contracts))
	assert.Equal(t, match.HomeTeam.TrainingPoints, uint16(34))
	assert.Equal(t, match.VisitorTeam.TrainingPoints, uint16(34))
	assert.NilError(t, match.Play1stHalf(*bc.Contracts))
	assert.Equal(t, match.HomeTeam.TrainingPoints, uint16(0))
	assert.Equal(t, match.VisitorTeam.TrainingPoints, uint16(0))
}

func TestPlay2ndHalfWithEmptyTeam(t *testing.T) {
	t.Parallel()
	engine := engine.NewMatch()
	err := engine.Play2ndHalf(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, engine.HomeGoals, uint8(0))
	assert.Equal(t, engine.VisitorGoals, uint8(0))
	assert.Equal(t, engine.HomeTeamSumSkills, uint32(0))
	assert.Equal(t, engine.VisitorTeamSumSkills, uint32(0))
	assert.Equal(t, engine.HomeTeam.MatchLog, "1823386170864456664588543800539808540283317251593298733231417759322490273792")
	assert.Equal(t, engine.VisitorTeam.MatchLog, "1823386170864456664588543800539808540283317251593298733231417759322490273792")
}

func TestPlayGame(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.Seed = [32]byte{0x2, 0x1}
	m.StartTime = big.NewInt(1570147200)
	m.HomeTeam.TeamID = "1"
	m.VisitorTeam.TeamID = "2"
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309339445376240739796176995438"))
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
		m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	golden.Assert(t, dump.Sdump(m), t.Name()+".starting.golden")
	err := m.Play1stHalf(*bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, dump.Sdump(m), t.Name()+".half.golden")
	err = m.Play2ndHalf(*bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, dump.Sdump(m), t.Name()+".ended.golden")
}

func TestPlay2ndHalf(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()

	m.Seed = sha256.Sum256([]byte("sdadfefe"))
	m.StartTime = big.NewInt(1570147200)
	m.HomeTeam.TeamID = "1"
	m.VisitorTeam.TeamID = "2"

	playerAligned1stHalf := SkillsFromString(t, "155218553451227483908160832387024632959402541639272557524812776")
	playerNotAligned1stHalf := SkillsFromString(t, "155218553445241173201653454034062339884596646390761857828783080")

	for i := 1; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, playerAligned1stHalf)
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
		m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	m.HomeTeam.Players[4].SetSkills(*bc.Contracts, playerNotAligned1stHalf)

	err := m.Play2ndHalf(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, m.HomeGoals, uint8(12))
	assert.Equal(t, m.VisitorGoals, uint8(0))
	assert.Equal(t, m.HomeTeamSumSkills, uint32(310768))
	assert.Equal(t, m.VisitorTeamSumSkills, uint32(0))
	assert.Equal(t, m.HomeTeam.MatchLog, "1854314176407607222731365934996113197647416408691401809684068772276829647276")
	assert.Equal(t, m.VisitorTeam.MatchLog, "1813668511995011514317266015948482734985807286120620942036680286308591992832")
	assert.Equal(t, m.HomeTeam.Players[0].Skills().String(), "0")
	assert.Equal(t, m.HomeTeam.Players[1].Skills().String(), "155218553469186416027682967445911512183820227384804656612901864")
	assert.Equal(t, m.VisitorTeam.Players[0].Skills().String(), "0")
	assert.Equal(t, m.VisitorTeam.Players[1].Skills().String(), "4600807814280360774460723191042511563333025959642222")
	assert.Equal(t, m.HomeTeam.TrainingPoints, uint16(102))
	assert.Equal(t, m.VisitorTeam.TrainingPoints, uint16(10))
}

func TestMatchPlayCheckGoalsWithEventGoals(t *testing.T) {
	t.Parallel()
	cases := []struct{ Seed [32]byte }{
		{sha256.Sum256([]byte("sdadfefe"))},
		{sha256.Sum256([]byte("pippo"))},
		{sha256.Sum256([]byte("4gfsg3564e5t"))},
	}
	playerNotAligned1stHalf := SkillsFromString(t, "155218553445241173201653454034062339884596646390761857828783080")

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", hex.EncodeToString(tc.Seed[:])), func(t *testing.T) {
			m := engine.NewMatch()
			m.Seed = tc.Seed
			m.StartTime = big.NewInt(1570147200)
			m.HomeTeam.TeamID = "1"
			m.VisitorTeam.TeamID = "2"
			for i := 0; i < 25; i++ {
				m.HomeTeam.Players[i].SetSkills(*bc.Contracts, playerNotAligned1stHalf)
				m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, playerNotAligned1stHalf)
				m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
				m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
			}
			golden.Assert(t, dump.Sdump(m), t.Name()+".starting.golden")
			err := m.Play1stHalf(*bc.Contracts)
			assert.NilError(t, err)
			golden.Assert(t, dump.Sdump(m), t.Name()+".half.golden")
			assert.Equal(t, m.HomeGoals, m.Events.HomeGoals())
			assert.Equal(t, m.VisitorGoals, m.Events.VisitorGoals())
			err = m.Play2ndHalf(*bc.Contracts)
			assert.NilError(t, err)
			golden.Assert(t, dump.Sdump(m), t.Name()+".ended.golden")
			assert.Equal(t, m.HomeGoals, m.Events.HomeGoals())
			assert.Equal(t, m.VisitorGoals, m.Events.VisitorGoals())
		})
	}
}

func TestMatchPlayerEvolution(t *testing.T) {
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "14606248079918261338806855269144928920528183545627247"))
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
	}
	// the player has great skills, but is initially very old
	assert.Equal(t, m.HomeTeam.Players[0].Defence, uint64(23264))
	assert.NilError(t, m.Play1stHalf(*bc.Contracts))
	// so after the evolution stage before 1st half begings, it generates a child/academy player:
	assert.Equal(t, m.HomeTeam.Players[0].Defence, uint64(1226))
	assert.NilError(t, m.Play2ndHalf(*bc.Contracts))
	// which is maintaind after playing 2nd half:
	assert.Equal(t, m.HomeTeam.Players[0].Defence, uint64(1226))
}

func TestMatchTeamSkillsEvolution(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte("18"))
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "14606248079918261338806855269144928920528183545627247"))
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
		m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	golden.Assert(t, dump.Sdump(m.HomeTeam.Skills()), t.Name()+".before.golden")
	assert.NilError(t, m.Play1stHalf(*bc.Contracts))
	golden.Assert(t, dump.Sdump(m.HomeTeam.Skills()), t.Name()+".half.golden")
	assert.NilError(t, m.Play2ndHalf(*bc.Contracts))
	golden.Assert(t, dump.Sdump(m.HomeTeam.Skills()), t.Name()+".end.golden")
}

func TestMatchRedCards(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte(string(4)))
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "14606248079918261338806855269144928920528183545627247"))
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
		m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	assert.NilError(t, m.Play1stHalf(*bc.Contracts))
	golden.Assert(t, dump.Sdump(m), t.Name()+".golden")
	event := m.Events[12]
	assert.Equal(t, event.Type, matchevents.EVNT_RED)
	assert.Equal(t, event.PrimaryPlayer, int16(10))
	assert.Equal(t, event.Team, int16(0))
	player := m.HomeTeam.Players[10]
	assert.Equal(t, player.Skills().String(), "1696941887453530621720210496306760960036050709342320410694255608")
	assert.Assert(t, player.RedCard)
}

func TestMatchHardInjury(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte("10"))
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "14606248079918261338806855269144928920528183545627247"))
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
		m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	assert.NilError(t, m.Play1stHalf(*bc.Contracts))
	golden.Assert(t, dump.Sdump(m), t.Name()+".golden")
	event := m.Events[12]
	assert.Equal(t, event.Type, matchevents.EVNT_HARD)
	assert.Equal(t, event.PrimaryPlayer, int16(10))
	assert.Equal(t, event.Team, int16(0))
	player := m.HomeTeam.Players[10]
	assert.Equal(t, player.Skills().String(), "1696941888399367713348376276074803265855382158607010962666947576")
	assert.Equal(t, player.InjuryMatchesLeft, uint8(5))
}

func TestMatchSoftInjury(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte("161"))
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "14606248079918261338806855269144928920528183545627247"))
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
		m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	assert.NilError(t, m.Play1stHalf(*bc.Contracts))
	golden.Assert(t, dump.Sdump(m), t.Name()+".golden")
	event := m.Events[12]
	assert.Equal(t, event.Type, matchevents.EVNT_SOFT)
	assert.Equal(t, event.PrimaryPlayer, int16(12))
	assert.Equal(t, event.Team, int16(0))
	player := m.HomeTeam.Players[12]
	assert.Equal(t, player.Skills().String(), "1696941887824681885523667954190423130674016214749983791848096760")
	assert.Equal(t, player.InjuryMatchesLeft, uint8(2))
}

func TestMatchEvents(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte(string(4)))
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "14606248079918261338806855269144928920528183545627247"))
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
		m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	golden.Assert(t, m.Events.DumpState(), t.Name()+".atStart.golden")
	assert.NilError(t, m.Play1stHalf(*bc.Contracts))
	golden.Assert(t, m.Events.DumpState(), t.Name()+".first.golden")
	assert.NilError(t, m.Play2ndHalf(*bc.Contracts))
	golden.Assert(t, m.Events.DumpState(), t.Name()+".second.golden")
}

func TestMatchJson(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte("161"))
	m.HomeTeam.Players[0].SetSkills(*bc.Contracts, SkillsFromString(t, "14606248079918261338806855269144928920528183545627247"))

	golden.Assert(t, string(m.ToJson()), t.Name()+".golden")
	input := golden.Get(t, t.Name()+".golden")
	match, err := engine.NewMatchFromJson(input)
	assert.NilError(t, err)
	assert.Equal(t, m.HomeTeam.Players[0].Skills().String(), match.HomeTeam.Players[0].Skills().String())
	assert.Equal(t, m.StartTime.String(), match.StartTime.String())
	assert.Equal(t, m.Seed, match.Seed)
}

func TestMatchHash(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	assert.Equal(t, fmt.Sprintf("%x", m.Hash()), "313cc8b62f5a3dd84903050f69c6dc30ed7b5d39d50a9286414f29fb76c92faf")
}

func TestMatchError1stHalf(t *testing.T) {
	t.Parallel()
	cases := []struct {
		File   string
		Output string
	}{
		{"3859fc1422bc9d7e58621e77466eb42c7db8cc2305687bfe41b23bc137e14d70.1st.error.json", "failed calculating visitor assignedTP: VM execution error."},
		// {"530796ade7bacc9b7d2e83246cc6fd46da9fb205d0fab24d80d6c8946a58b294.1st.error.json", "failed calculating home assignedTP: VM execution error."},
		{"9a78b84120c90d40da0fce05cbab1bf539bb3a68cb835886e01af6ddaaf4aca9.1st.error.json", "failed calculating visitor assignedTP: VM execution error."},
		{"9cf953e0438bdd61de9b78b713c04384d67d15feb6e809de10f616ee1f812c65.1st.error.json", "failed calculating home assignedTP: VM execution error."},
	}
	for _, tc := range cases {
		t.Run(tc.File, func(t *testing.T) {
			input := golden.Get(t, t.Name())
			match, err := engine.NewMatchFromJson(input)
			assert.NilError(t, err)
			matchLog, _ := new(big.Int).SetString(match.HomeTeam.MatchLog, 10)
			decodedHomeMatchLog, err := bc.Contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, matchLog, true)
			assert.NilError(t, err)
			assert.Equal(t, uint32(match.HomeTeam.TrainingPoints), decodedHomeMatchLog[3])
			matchLog, _ = new(big.Int).SetString(match.VisitorTeam.MatchLog, 10)
			decodedVisitorMatchLog, err := bc.Contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, matchLog, true)
			assert.NilError(t, err)
			assert.Equal(t, uint32(match.VisitorTeam.TrainingPoints), decodedVisitorMatchLog[3])
			err = match.Play1stHalf(*bc.Contracts)
			assert.Assert(t, err != nil)
			assert.Equal(t, err.Error(), tc.Output)
		})
	}
}

func TestMatchError2ndHalf(t *testing.T) {
	t.Parallel()
	cases := []struct {
		File   string
		Output string
	}{
		// {"86dc12c640c057604e8e384f88bb15a3c6b15a43b8b76f34fbe6e320516095e1.2nd.error.json", "VM execution error."},
		// {"a1943a63802b87bf2247f96fc1ee7d80482354623d997ea9f115e59a3d94d2db.2nd.error.json", "VM execution error."},
		// {"bfdd1ce80cebf6417dd98a419da55fd8428b7e7122123464ac037d5ef4a3aaec.2nd.error.json", "VM execution error."},
		// {"d2783956a1153d9da33a222293ce5c0751bdc98f253736203e8a92b5b6f081b8.2nd.error.json", "VM execution error."},
	}
	for _, tc := range cases {
		t.Run(tc.File, func(t *testing.T) {
			input := golden.Get(t, t.Name())
			match, err := engine.NewMatchFromJson(input)
			assert.NilError(t, err)
			assert.Error(t, match.Play2ndHalf(*bc.Contracts), tc.Output)
		})
	}
}

func TestMatchEventsGeneration(t *testing.T) {
	t.Parallel()
	for j := 0; j < 100; j++ {
		t.Run(fmt.Sprintf("%d", j), func(t *testing.T) {
			t.Parallel()
			m := engine.NewMatch()
			m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
			m.Seed = sha256.Sum256([]byte(fmt.Sprintf("%d", j)))
			m.HomeTeam.TeamID = "274877906944"
			m.VisitorTeam.TeamID = "274877906945"
			for i := 0; i < 24; i++ {
				m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "14606248079918261338806855269144928920528183545627247"))
				m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
			}
			for i := 1; i < 4; i++ {
				m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
				m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
			}
			assert.NilError(t, m.Play1stHalf(*bc.Contracts))
		})
	}
}

func TestFromTheField(t *testing.T) {
	t.Parallel()
	t.Run("InconsistentPositionPlayerId", func(t *testing.T) {
		input := golden.Get(t, t.Name()+"/b65d48b5a6a4075098e6a996bece8f5aeec8b2ac73c6d62a8de8a18bc28a5230.1st.error.json")
		match, err := engine.NewMatchFromJson(input)
		assert.NilError(t, err)
		assert.NilError(t, match.Play1stHalf(*bc.Contracts))
	})
	// the following should fail because user saw 45 TPs when he actually had only 44 available:
	t.Run("Failing0", func(t *testing.T) {
		input := golden.Get(t, t.Name()+"/0498232f79495530fa199c6d51fa51b2bfb22989b01e5f390eced6e729b04102.1st.error.json")
		match, err := engine.NewMatchFromJson(input)
		assert.NilError(t, err)
		assert.Error(t, match.Play1stHalf(*bc.Contracts), "failed calculating visitor assignedTP: VM execution error.")
	})
	t.Run("Failing1", func(t *testing.T) {
		input := golden.Get(t, t.Name()+"/fe6e996fc594c5043f29040561cc95c02c0f68ccdc80047a30e42e74f3b402f8.2nd.error.json")
		match, err := engine.NewMatchFromJson(input)
		assert.NilError(t, err)
		assert.Error(t, match.Play2ndHalf(*bc.Contracts), "failed play2ndHalfAndEvolve: VM execution error.")
		// assert.NilError(t, match.Play2ndHalf(*bc.Contracts))
	})
}
