package engine_test

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/engine"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestDefaultMatch(t *testing.T) {
	t.Parallel()
	match, err := engine.NewMatch(bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, match.DumpState(), t.Name()+".golden")
}
func TestDefaultValues(t *testing.T) {
	t.Parallel()
	engine, err := engine.NewMatch(bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, engine != nil, "engine is nil")
}

func TestMatchState(t *testing.T) {
	t.Parallel()
	m, _ := engine.NewMatch(bc.Contracts)
	m.Seed = [32]byte{0x2, 0x1}
	m.StartTime = big.NewInt(1570147200)
	m.HomeTeam.TeamID = big.NewInt(int64(1))
	m.VisitorTeam.TeamID = big.NewInt(int64(2))
	homeSkill := uint16(1134)
	visitorSkill := uint16(2344)
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i] = engine.CreateDummyPlayer(t, bc.Contracts, 33, homeSkill, homeSkill, homeSkill, homeSkill, homeSkill)
		m.VisitorTeam.Players[i] = engine.CreateDummyPlayer(t, bc.Contracts, 18, visitorSkill, visitorSkill, visitorSkill, visitorSkill, visitorSkill)
	}
	assert.Equal(t, m.State, engine.Starting)
	err := m.Play1stHalf()
	assert.NilError(t, err)
	assert.Equal(t, m.State, engine.Half)
	err = m.Play2ndHalf()
	assert.NilError(t, err)
	assert.Equal(t, m.State, engine.Ended)
}

func TestPlay1stHalfTwice(t *testing.T) {
	t.Parallel()
	match, _ := engine.NewMatch(bc.Contracts)
	err := match.Play1stHalf()
	assert.NilError(t, err)
	err = match.Play1stHalf()
	assert.Error(t, err, "Wrong state Half")
}

func TestPlay2ndHalfWithout1st(t *testing.T) {
	t.Parallel()
	m, _ := engine.NewMatch(bc.Contracts)
	err := m.Play2ndHalf()
	assert.Error(t, err, "Wrong state Starting")
}

func TestPlay1stHalfWithEmptyTeam(t *testing.T) {
	t.Parallel()
	match, _ := engine.NewMatch(bc.Contracts)
	err := match.Play1stHalf()
	assert.NilError(t, err)
	assert.Equal(t, match.HomeGoals, uint8(0))
	assert.Equal(t, match.VisitorGoals, uint8(0))
	assert.Equal(t, match.HomeMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
	assert.Equal(t, match.VisitorMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
	golden.Assert(t, match.DumpState(), t.Name()+".golden")
}

func TestPlay2ndHalfWithEmptyTeam(t *testing.T) {
	t.Skip("TODO: *************************  REACTIVE ***************************")
	t.Parallel()
	engine, _ := engine.NewMatch(bc.Contracts)
	err := engine.Play2ndHalf()
	assert.NilError(t, err)
	assert.Equal(t, engine.HomeGoals, uint8(0))
	assert.Equal(t, engine.VisitorGoals, uint8(0))
	assert.Equal(t, engine.HomeMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
	assert.Equal(t, engine.VisitorMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
}

func TestPlayGameStress(t *testing.T) {
	t.Skip("TODO: *************************  REACTIVE ***************************")
	t.Parallel()
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("age%v", i), func(t *testing.T) {
			m, _ := engine.NewMatch(bc.Contracts)
			m.Seed = [32]byte{0x2, 0x1}
			m.StartTime = big.NewInt(1570147200)
			m.HomeTeam.TeamID = big.NewInt(int64(1))
			m.VisitorTeam.TeamID = big.NewInt(int64(2))
			visitorSkill := uint16(2344)
			for i := 0; i < 25; i++ {
				m.HomeTeam.Players[i] = engine.NewPlayerFromSkills("16573429227295117480385309339445376240739796176995438")
				// m.VisitorTeam.Players[i] = engine.NewPlayerFromSkills("34257599038999042790649896760714712618135701934311720")
				m.VisitorTeam.Players[i] = engine.CreateDummyPlayer(t, bc.Contracts, uint16(i), visitorSkill, visitorSkill, visitorSkill, visitorSkill, visitorSkill)
			}
			err := m.Play1stHalf()
			assert.NilError(t, err)
			err = m.Play2ndHalf()
			assert.NilError(t, err)
			assert.Equal(t, m.HomeGoals, m.Events.HomeGoals())
			assert.Equal(t, m.VisitorGoals, m.Events.VisitorGoals())
		})
	}
}

func TestPlayGame(t *testing.T) {
	t.Parallel()
	m, _ := engine.NewMatch(bc.Contracts)
	m.Seed = [32]byte{0x2, 0x1}
	m.StartTime = big.NewInt(1570147200)
	m.HomeTeam.TeamID = big.NewInt(int64(1))
	m.VisitorTeam.TeamID = big.NewInt(int64(2))
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i] = engine.NewPlayerFromSkills("16573429227295117480385309339445376240739796176995438")
		m.VisitorTeam.Players[i] = engine.NewPlayerFromSkills("16573429227295117480385309340654302060354425351701614")
	}
	golden.Assert(t, m.DumpState(), t.Name()+".starting.golden")
	err := m.Play1stHalf()
	assert.NilError(t, err)
	golden.Assert(t, m.DumpState(), t.Name()+".half.golden")
	err = m.Play2ndHalf()
	assert.NilError(t, err)
	golden.Assert(t, m.DumpState(), t.Name()+".ended.golden")
}

func TestPlay2ndHalf(t *testing.T) {
	t.Skip("TODO: *************************  REACTIVE ***************************")
	t.Parallel()
	m, _ := engine.NewMatch(bc.Contracts)
	homePlayer := engine.NewPlayerFromSkills("146156532686539503615416807207209880594713965887498")
	visitorPlayer := engine.NewPlayerFromSkills("730757187618900670896890173308251570218123297685554")
	m.HomeTeam.Players[0] = homePlayer
	m.VisitorTeam.Players[0] = visitorPlayer
	err := m.Play2ndHalf()
	assert.NilError(t, err)
	assert.Equal(t, m.HomeGoals, uint8(0))
	assert.Equal(t, m.VisitorGoals, uint8(0))
	assert.Equal(t, m.HomeMatchLog.String(), "166195960289441810257652497224293923324982848796288083926844440576")
	assert.Equal(t, m.VisitorMatchLog.String(), "824397783217924227119640170247234125318077195049720029266288050176")
	assert.Equal(t, m.HomeTeam.Players[0].Skills().String(), "146156532686539503615416807207209880594713965887498")
	assert.Equal(t, m.HomeTeam.Players[1].Skills().String(), "0")
	assert.Equal(t, m.VisitorTeam.Players[0].Skills().String(), "730757187618900670896890173308251570218123297685554")
	assert.Equal(t, m.VisitorTeam.Players[1].Skills().String(), "0")
}

func TestMatchPlayCheckGoalsWithEventGoals(t *testing.T) {
	t.Parallel()
	cases := []struct{ Seed string }{
		{"sdadfefe"},
		{"pippo"},
		{"4gfsg3564e5t"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			m, _ := engine.NewMatch(bc.Contracts)
			m.Seed = sha256.Sum256([]byte(tc.Seed))
			m.StartTime = big.NewInt(1570147200)
			m.HomeTeam.TeamID = big.NewInt(int64(1))
			m.VisitorTeam.TeamID = big.NewInt(int64(2))
			for i := 0; i < 25; i++ {
				m.HomeTeam.Players[i] = engine.NewPlayerFromSkills("16573429227295117480385309339445376240739796176995438")
				m.VisitorTeam.Players[i] = engine.NewPlayerFromSkills("16573429227295117480385309340654302060354425351701614")
			}
			golden.Assert(t, m.DumpState(), t.Name()+".starting.golden")
			err := m.Play1stHalf()
			assert.NilError(t, err)
			golden.Assert(t, m.DumpState(), t.Name()+".half.golden")
			assert.Equal(t, m.HomeGoals, m.Events.HomeGoals())
			assert.Equal(t, m.VisitorGoals, m.Events.VisitorGoals())
			err = m.Play2ndHalf()
			assert.NilError(t, err)
			golden.Assert(t, m.DumpState(), t.Name()+".ended.golden")
			assert.Equal(t, m.HomeGoals, m.Events.HomeGoals())
			assert.Equal(t, m.VisitorGoals, m.Events.VisitorGoals())
		})
	}
}

func TestPlay2ndHalf_goals(t *testing.T) {
	t.Skip("TODO: *************************  REACTIVE ***************************")
	t.Parallel()
	cases := []struct {
		HomeAge             uint16
		VisitorAge          uint16
		HomeSkill           uint16
		VisitorSkill        uint16
		ExpectedHomeGoal    uint8
		ExpectedVisitorGoal uint8
	}{
		{21, 30, 10, 50, 0, 3},
		{30, 18, 1233, 2344, 0, 0},
		{55, 18, 12, 2344, 0, 10},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("HAge:%v VAge:%v HSkill:%v VSkills:%v HGoals:%v VGoals:%v", tc.HomeAge, tc.VisitorAge, tc.HomeSkill, tc.VisitorSkill, tc.ExpectedHomeGoal, tc.ExpectedVisitorGoal), func(t *testing.T) {
			m, _ := engine.NewMatch(bc.Contracts)
			m.Seed = [32]byte{0x1, 0x1f}
			m.StartTime = big.NewInt(1570147200)
			m.HomeTeam.TeamID = big.NewInt(int64(1))
			m.VisitorTeam.TeamID = big.NewInt(int64(2))
			for i := 0; i < 25; i++ {
				m.HomeTeam.Players[i] = engine.CreateDummyPlayer(t, bc.Contracts, tc.HomeAge, tc.HomeSkill, tc.HomeSkill, tc.HomeSkill, tc.HomeSkill, tc.HomeSkill)
				m.VisitorTeam.Players[i] = engine.CreateDummyPlayer(t, bc.Contracts, tc.VisitorAge, tc.VisitorSkill, tc.VisitorSkill, tc.VisitorSkill, tc.VisitorSkill, tc.VisitorSkill)
			}
			err := m.Play2ndHalf()
			assert.NilError(t, err)
			assert.Equal(t, m.HomeGoals, tc.ExpectedHomeGoal)
			assert.Equal(t, m.VisitorGoals, tc.ExpectedVisitorGoal)
		})
	}
}
