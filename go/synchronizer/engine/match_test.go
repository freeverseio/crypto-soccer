package engine_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

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
	match.HomeTeam.Owner = "0x433"
	match.VisitorTeam.Owner = "0x433"
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
	match.HomeTeam.Owner = "0x433"
	match.VisitorTeam.Owner = "0x433"
	assert.NilError(t, match.Play1stHalf(*bc.Contracts))
	assert.Equal(t, match.HomeTeam.TrainingPoints, uint16(0))
	assert.Equal(t, match.VisitorTeam.TrainingPoints, uint16(0))
	assert.NilError(t, match.Play2ndHalf(*bc.Contracts))
	assert.Equal(t, match.HomeTeam.TrainingPoints, uint16(10))
	assert.Equal(t, match.VisitorTeam.TrainingPoints, uint16(10))
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
	assert.Equal(t, engine.HomeTeam.MatchLog, "453417128002043887693956307131195271170887450874663967603880824822023847936")
	assert.Equal(t, engine.VisitorTeam.MatchLog, "453417128002043887693956307131195271170887450874663967603880824822023847936")
}

func TestPlayGame(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.Seed = [32]byte{0x2, 0x1}
	m.StartTime = big.NewInt(1570147200)
	m.HomeTeam.TeamID = "1"
	m.VisitorTeam.TeamID = "2"
	m.HomeTeam.Owner = "0x433"
	m.VisitorTeam.Owner = "0x433"
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
	m.HomeTeam.Owner = "0x433"
	m.VisitorTeam.Owner = "0x433"

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
	assert.Equal(t, m.HomeTeam.MatchLog, "1367431245678592349487822957430078931600781052235089130055218559588686578092")
	assert.Equal(t, m.VisitorTeam.MatchLog, "453417127998752878579313895046885332805772749864698495872613665095326629888")
	assert.Equal(t, m.HomeTeam.Players[0].Skills().String(), "0")
	assert.Equal(t, m.HomeTeam.Players[1].Skills().String(), "155218553469186416027682967445911512183820227384804656612901864")
	assert.Equal(t, m.VisitorTeam.Players[0].Skills().String(), "0")
	assert.Equal(t, m.VisitorTeam.Players[1].Skills().String(), "4600807814280360774460723191042511563333025959642222")
	assert.Equal(t, m.HomeTeam.TrainingPoints, uint16(95))
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
			m.HomeTeam.Owner = "0x433"
			m.VisitorTeam.Owner = "0x433"
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
	m.HomeTeam.Owner = "0x433"
	m.VisitorTeam.Owner = "0x433"
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
	m.HomeTeam.Owner = "0x433"
	m.VisitorTeam.Owner = "0x433"
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
	m.HomeTeam.Owner = "0x433"
	m.VisitorTeam.Owner = "0x433"
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
	assert.Equal(t, player.Skills().String(), "41189051263162475633439470877683173079621711484548237131060872184")
	assert.Assert(t, player.RedCard)
	assert.Assert(t, player.YellowCard1stHalf)
}

func TestMatchHardInjury(t *testing.T) {
	// note: the strings used in this test are a bit crazy, and players are superold
	// note: that is why after evolving, they become normal players (children)
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte("10"))
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	m.HomeTeam.Owner = "0x433"
	m.VisitorTeam.Owner = "0x433"
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "14606248079918261338806855269144928920528183545627247"))
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "15324955408660339766675662613581148386300673993530595607528"))
		m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	golden.Assert(t, dump.Sdump(m), t.Name()+".starting.golden")
	assert.NilError(t, m.Play1stHalf(*bc.Contracts))
	golden.Assert(t, dump.Sdump(m), t.Name()+".half.golden")
	event := m.Events[12]
	assert.Equal(t, event.Type, matchevents.EVNT_HARD)
	assert.Equal(t, event.PrimaryPlayer, int16(6))
	assert.Equal(t, event.Team, int16(0))
	player := m.HomeTeam.Players[6]
	assert.Equal(t, player.Skills().String(), "14860978346969016050588129736533607305717269083675649869455819768")
	assert.Equal(t, player.InjuryMatchesLeft, uint8(5))
}

func TestMatchHardInjuryAmongBots(t *testing.T) {
	// note: the strings used in this test correcpond to all players = 1000, but superold
	// note: since the teams are bots, they dont evolve, and remain at 1000
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte("10"))
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	// m.HomeTeam.Owner = "0x433"
	// m.VisitorTeam.Owner = "0x433"
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "15324955408660339766675662613581148386300673993530595607528"))
		m.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "15324955408660339766675662613581148386300673993530595607528"))
		m.HomeTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
		m.VisitorTeam.Players[i].SetPlayerId(new(big.Int).SetUint64(21342314523))
	}
	// check that skills are 1000 before and after
	assert.Equal(t, m.HomeTeam.Players[0].Defence, uint64(1000))
	golden.Assert(t, dump.Sdump(m), t.Name()+".starting.golden")
	assert.NilError(t, m.Play1stHalf(*bc.Contracts))
	golden.Assert(t, dump.Sdump(m), t.Name()+".half.golden")
	for i := 0; i < 25; i++ {
		assert.Equal(t, m.HomeTeam.Players[i].RedCard, false)
		assert.Equal(t, m.HomeTeam.Players[i].YellowCard1stHalf, false)
		assert.Equal(t, m.HomeTeam.Players[i].Tiredness, 0)
		assert.Equal(t, m.HomeTeam.Players[i].InjuryMatchesLeft, uint8(0))
		assert.Equal(t, m.HomeTeam.Players[i].Defence, uint64(1000))
	}
	assert.NilError(t, m.Play2ndHalf(*bc.Contracts))
	golden.Assert(t, dump.Sdump(m), t.Name()+".end.golden")
	for i := 0; i < 25; i++ {
		assert.Equal(t, m.HomeTeam.Players[i].RedCard, false)
		assert.Equal(t, m.HomeTeam.Players[i].YellowCard1stHalf, false)
		assert.Equal(t, m.HomeTeam.Players[i].Tiredness, 0)
		assert.Equal(t, m.HomeTeam.Players[i].InjuryMatchesLeft, uint8(0))
		assert.Equal(t, m.HomeTeam.Players[i].Defence, uint64(1000))
	}
}

func TestMatchSoftInjury(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte("161"))
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	m.HomeTeam.Owner = "0x433"
	m.VisitorTeam.Owner = "0x433"
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
	assert.Equal(t, event.PrimaryPlayer, int16(6))
	assert.Equal(t, event.Team, int16(0))
	player := m.HomeTeam.Players[6]
	assert.Equal(t, player.Skills().String(), "14860978346394330222763421414649227170535903139818622698636968952")
	assert.Equal(t, player.InjuryMatchesLeft, uint8(2))
}

func TestMatchEvents(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
	m.Seed = sha256.Sum256([]byte(string(4)))
	m.HomeTeam.TeamID = "274877906944"
	m.VisitorTeam.TeamID = "274877906945"
	m.HomeTeam.Owner = "0x433"
	m.VisitorTeam.Owner = "0x433"
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
	assert.Equal(t, fmt.Sprintf("%x", m.Hash()), "4a77551588ccd5cf4fbbb3ea36d6d651b56a5c354621353b2305e20e22ecb800")
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
			m.HomeTeam.Owner = "0x433"
			m.VisitorTeam.Owner = "0x433"
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
	// These tests used to fail badly, returning error. Now they fail gracefully, returning 0-0 and valid skills, logs, etc.
	t.Parallel()
	t.Run("0498232f79495530fa199c6d51fa51b2bfb22989b01e5f390eced6e729b04102.1st.error.json", func(t *testing.T) {
		input := golden.Get(t, t.Name())
		match, err := engine.NewMatchFromJson(input)
		assert.NilError(t, err)
		match.HomeTeam.Owner = "0x433"
		match.VisitorTeam.Owner = "0x433"
		assert.Error(t, match.Play1stHalf(*bc.Contracts), "failed calculating visitor assignedTP: one of the assigned TPs is too large")
	})
	t.Run("fe6e996fc594c5043f29040561cc95c02c0f68ccdc80047a30e42e74f3b402f8.2nd.error.json", func(t *testing.T) {
		input := golden.Get(t, t.Name())
		match, err := engine.NewMatchFromJson(input)
		assert.NilError(t, err)
		match.HomeTeam.Owner = "0x433"
		match.VisitorTeam.Owner = "0x433"
		assert.Error(t, match.Play2ndHalf(*bc.Contracts), "BLOCKCHAIN ERROR!!!! play2ndHalfAndEvolve: Blockchain returned error code: 16")
	})
	t.Run("a102d90303aafcdae29c09bc6b338a50048b9cd4d8fa1942cf315bb7e3736aac.2nd.error.json", func(t *testing.T) {
		input := golden.Get(t, t.Name())
		match, err := engine.NewMatchFromJson(input)
		assert.NilError(t, err)
		match.HomeTeam.Owner = "0x433"
		match.VisitorTeam.Owner = "0x433"
		assert.Error(t, match.Play2ndHalf(*bc.Contracts), "BLOCKCHAIN ERROR!!!! play2ndHalfAndEvolve: Blockchain returned error code: 13")
	})
}
