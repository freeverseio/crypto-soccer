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
	homePlayer := engine.NewPlayer()
	homePlayer.SetSkills(*bc.Contracts, SkillsFromString(t, "146156532686539503615416807207209880594713965887498"))
	visitorPlayer := engine.NewPlayer()
	visitorPlayer.SetSkills(*bc.Contracts, SkillsFromString(t, "730757187618900670896890173308251570218123297685554"))
	m.HomeTeam.Players[0] = homePlayer
	m.VisitorTeam.Players[0] = visitorPlayer
	err := m.Play2ndHalf(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, m.HomeGoals, uint8(0))
	assert.Equal(t, m.VisitorGoals, uint8(0))
	assert.Equal(t, m.HomeTeam.MatchLog, "1823386170864456664588543800539808540283317251593298733231417759322490273792")
	assert.Equal(t, m.VisitorTeam.MatchLog, "1823386170864456664588543800539808540283317251593298733231417759322490273792")
	assert.Equal(t, m.HomeTeam.Players[0].Skills().String(), "146173659658851975133989506638843274536350558846986")
	assert.Equal(t, m.HomeTeam.Players[1].Skills().String(), "0")
	assert.Equal(t, m.VisitorTeam.Players[0].Skills().String(), "730774314591213142415462872739884964159759890645042")
	assert.Equal(t, m.VisitorTeam.Players[1].Skills().String(), "0")
	assert.Equal(t, m.HomeTeam.TrainingPoints, uint16(32))
	assert.Equal(t, m.VisitorTeam.TrainingPoints, uint16(32))
}

func TestMatchPlayCheckGoalsWithEventGoals(t *testing.T) {
	t.Parallel()
	cases := []struct{ Seed [32]byte }{
		{sha256.Sum256([]byte("sdadfefe"))},
		{sha256.Sum256([]byte("pippo"))},
		{sha256.Sum256([]byte("4gfsg3564e5t"))},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", hex.EncodeToString(tc.Seed[:])), func(t *testing.T) {
			m := engine.NewMatch()
			m.Seed = tc.Seed
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
	assert.Equal(t, m.HomeTeam.Players[0].Defence, uint64(955))
	assert.NilError(t, m.Play1stHalf(*bc.Contracts))
	assert.Equal(t, m.HomeTeam.Players[0].Defence, uint64(1237))
	assert.NilError(t, m.Play2ndHalf(*bc.Contracts))
	assert.Equal(t, m.HomeTeam.Players[0].Defence, uint64(1237))
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
	assert.Equal(t, player.Skills().String(), "444839120007985571215348664084887401221731547818249502887980205736758")
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
	assert.Equal(t, player.Skills().String(), "444839120007985571216250684626677567866560384550941583814174101603126")
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
	assert.Equal(t, player.Skills().String(), "444839120007985571215702621512678479272234002738672977681803126899510")
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
	fmt.Sprintf("%x", m.Hash())
	assert.Equal(t, fmt.Sprintf("%x", m.Hash()), "191607c1c3760dbeedca48c9bd33d2a793d34c19b7f01469aafdf08ee309fb5a")
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

func TestFromTheField(t *testing.T) {
	input := golden.Get(t, t.Name()+"/be021ab477ff432a6152482cc409acfb3bf525037fbd69bc2b41f83c5d67433a.toStorage.error.json")
	match, err := engine.NewMatchFromJson(input)
	assert.NilError(t, err)
	matchLog, _ := new(big.Int).SetString(match.HomeTeam.MatchLog, 10)
	decodedHomeMatchLog, err := bc.Contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, matchLog, true)
	assert.NilError(t, err)
	matchLog, _ = new(big.Int).SetString(match.VisitorTeam.MatchLog, 10)
	decodedVisitorMatchLog, err := bc.Contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, matchLog, true)
	assert.NilError(t, err)

	t.Log(decodedHomeMatchLog)
	t.Log(decodedVisitorMatchLog)
}
