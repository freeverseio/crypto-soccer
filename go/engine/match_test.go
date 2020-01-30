package engine_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/engine"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestDefaultMatch(t *testing.T) {
	t.Parallel()
	match := engine.NewMatch()
	golden.Assert(t, match.DumpState(), t.Name()+".golden")
}
func TestDefaultValues(t *testing.T) {
	t.Parallel()
	engine := engine.NewMatch()
	assert.Assert(t, engine != nil, "engine is nil")
}

func TestPlay1stHalfWithEmptyTeam(t *testing.T) {
	t.Parallel()
	match := engine.NewMatch()
	err := match.Play1stHalf(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, match.HomeGoals, uint8(0))
	assert.Equal(t, match.VisitorGoals, uint8(0))
	assert.Equal(t, match.HomeMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
	assert.Equal(t, match.VisitorMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
	golden.Assert(t, match.DumpState(), t.Name()+".golden")
}

func TestPlay2ndHalfWithEmptyTeam(t *testing.T) {
	t.Parallel()
	engine := engine.NewMatch()
	err := engine.Play2ndHalf(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, engine.HomeGoals, uint8(0))
	assert.Equal(t, engine.VisitorGoals, uint8(0))
	assert.Equal(t, engine.HomeMatchLog.String(), "1656419124875239866305548088508421623415643372415742952246406776530927616")
	assert.Equal(t, engine.VisitorMatchLog.String(), "1656419124875239866305548088508421623415643372415742952246406776530927616")
}

func TestPlayGame(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	m.Seed = [32]byte{0x2, 0x1}
	m.StartTime = big.NewInt(1570147200)
	m.HomeTeam.TeamID = big.NewInt(int64(1))
	m.VisitorTeam.TeamID = big.NewInt(int64(2))
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i] = engine.NewPlayerFromSkills("16573429227295117480385309339445376240739796176995438")
		m.VisitorTeam.Players[i] = engine.NewPlayerFromSkills("16573429227295117480385309340654302060354425351701614")
	}
	golden.Assert(t, m.DumpState(), t.Name()+".starting.golden")
	err := m.Play1stHalf(*bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, m.DumpState(), t.Name()+".half.golden")
	err = m.Play2ndHalf(*bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, m.DumpState(), t.Name()+".ended.golden")
}

func TestPlay2ndHalf(t *testing.T) {
	t.Parallel()
	m := engine.NewMatch()
	homePlayer := engine.NewPlayerFromSkills("146156532686539503615416807207209880594713965887498")
	visitorPlayer := engine.NewPlayerFromSkills("730757187618900670896890173308251570218123297685554")
	m.HomeTeam.Players[0] = homePlayer
	m.VisitorTeam.Players[0] = visitorPlayer
	err := m.Play2ndHalf(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, m.HomeGoals, uint8(0))
	assert.Equal(t, m.VisitorGoals, uint8(0))
	assert.Equal(t, m.HomeMatchLog.String(), "1656419124875239866305548088508421623415643372415742952246406776530927616")
	assert.Equal(t, m.VisitorMatchLog.String(), "1656419124875239866305548088508421623415643372415742952246406776530927616")
	assert.Equal(t, m.HomeTeam.Players[0].Skills().String(), "146150823695768679775892574063332082614168434901002")
	assert.Equal(t, m.HomeTeam.Players[1].Skills().String(), "0")
	assert.Equal(t, m.VisitorTeam.Players[0].Skills().String(), "730751478628129847057365940164373772237577766699058")
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
			m.HomeTeam.TeamID = big.NewInt(int64(1))
			m.VisitorTeam.TeamID = big.NewInt(int64(2))
			for i := 0; i < 25; i++ {
				m.HomeTeam.Players[i] = engine.NewPlayerFromSkills("16573429227295117480385309339445376240739796176995438")
				m.VisitorTeam.Players[i] = engine.NewPlayerFromSkills("16573429227295117480385309340654302060354425351701614")
			}
			golden.Assert(t, m.DumpState(), t.Name()+".starting.golden")
			err := m.Play1stHalf(*bc.Contracts)
			assert.NilError(t, err)
			golden.Assert(t, m.DumpState(), t.Name()+".half.golden")
			assert.Equal(t, m.HomeGoals, m.Events.HomeGoals())
			assert.Equal(t, m.VisitorGoals, m.Events.VisitorGoals())
			err = m.Play2ndHalf(*bc.Contracts)
			assert.NilError(t, err)
			golden.Assert(t, m.DumpState(), t.Name()+".ended.golden")
			assert.Equal(t, m.HomeGoals, m.Events.HomeGoals())
			assert.Equal(t, m.VisitorGoals, m.Events.VisitorGoals())
		})
	}
}
