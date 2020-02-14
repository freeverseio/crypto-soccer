package engine_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
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
	err := match.Play1stHalf(*bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, dump.Sdump(match), t.Name()+".golden")
	assert.Equal(t, match.HomeGoals, uint8(0))
	assert.Equal(t, match.VisitorGoals, uint8(0))
}

func TestPlay2ndHalfWithEmptyTeam(t *testing.T) {
	t.Parallel()
	engine := engine.NewMatch()
	err := engine.Play2ndHalf(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, engine.HomeGoals, uint8(0))
	assert.Equal(t, engine.VisitorGoals, uint8(0))
	assert.Equal(t, engine.HomeMatchLog.String(), "1822502747332067472423741020992967168868778470381174956540879256639596658688")
	assert.Equal(t, engine.VisitorMatchLog.String(), "1822502747332067472423741020992967168868778470381174956540879256639596658688")
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
	assert.Equal(t, m.HomeMatchLog.String(), "1822502747332067472423741020992967168868778470381174956540879256639596658688")
	assert.Equal(t, m.VisitorMatchLog.String(), "1822502747332067472423741020992967168868778470381174956540879256639596658688")
	assert.Equal(t, m.HomeTeam.Players[0].Skills().String(), "146173659658851975133989506638843274536350558846986")
	assert.Equal(t, m.HomeTeam.Players[1].Skills().String(), "0")
	assert.Equal(t, m.VisitorTeam.Players[0].Skills().String(), "730774314591213142415462872739884964159759890645042")
	assert.Equal(t, m.VisitorTeam.Players[1].Skills().String(), "0")
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

// func TestMatchFromStorage(t *testing.T) {
// 	t.Parallel()
// 	tx, err := db.Begin()
// 	assert.NilError(t, err)
// 	defer tx.Rollback()
// 	stoMatch := storage.Match{}
// 	stoHomeTeam := storage.Team{}
// 	stoVisitorTeam := storage.Team{}
// 	stoHomePlayers := []*storage.Player{&storage.Player{}}
// 	stoHomePlayers[0].ShirtNumber = 4
// 	stoHomePlayers[0].EncodedSkills = SkillsFromString(t, "40439920000726868070503716865792521545121682176182486071370780491777")
// 	assert.NilError(t, stoHomePlayers[0].Insert(tx))
// 	stoVisitorPlayers := []*storage.Player{}
// 	match := engine.NewMatchFromStorage(
// 		stoMatch,
// 		stoHomeTeam,
// 		stoVisitorTeam,
// 		stoHomePlayers,
// 		stoVisitorPlayers,
// 	)
// 	golden.Assert(t, dump.Sdump(match), t.Name()+".golden")
// 	assert.NilError(t, match.ToStorage(*bc.Contracts, tx))
// 	golden.Assert(t, dump.Sdump(match), t.Name()+".after.toStorage.golden")
// }

// func TestMatchToStorage(t *testing.T) {
// 	t.Parallel()
// 	tx, err := db.Begin()
// 	assert.NilError(t, err)
// 	defer tx.Rollback()

// 	match := engine.NewMatch()
// 	err = match.ToStorage(*bc.Contracts, tx)
// 	assert.NilError(t, err)
// }
