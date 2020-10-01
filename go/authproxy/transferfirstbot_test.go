package authproxy_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/authproxy"
	"github.com/stretchr/testify/assert"
)

func TestIsTransferFirstBotMutation(t *testing.T) {
	t.Skip("has to be reactivated")
	m := "mutation {transferFirstBotToAddr(timezone: 10, countryIdxInTimezone: 1000, address: \"0x02\")}"
	match, _ := authproxy.IsTransferFirstBotMutation(m)
	assert.True(t, match)

	// spaces should not affect matching
	m = "mutation  {    transferFirstBotToAddr   (timezone    :   10 ,  countryIdxInTimezone :   10 ,    address: \"0x02\"   ) }"
	match, _ = authproxy.IsTransferFirstBotMutation(m)
	assert.True(t, match)

	// wrong param order should fail
	m = "mutation {transferFirstBotToAddr(countryIdxInTimezone: 1000, timezone: 10, address: \"0x02\")}"
	match, _ = authproxy.IsTransferFirstBotMutation(m)
	assert.False(t, match)

	// wrong param name should fail (Timezone should be timezone)
	m = "mutation {transferFirstBotToAddr(Timezone: 10, countryIdxInTimezone: 1000, address: \"0x02\")}"
	match, _ = authproxy.IsTransferFirstBotMutation(m)
	assert.False(t, match)

	// wrong method name should fail
	m = "mutation {foo(Timezone: 10, countryIdxInTimezone: 1000, address: \"0x02\")}"
	match, _ = authproxy.IsTransferFirstBotMutation(m)
	assert.False(t, match)

	// not using mutation should fail
	m = "bar {transferFirstBotToAddr(timezone: 10, countryIdxInTimezone: 1000, address: \"0x02\")}"
	match, _ = authproxy.IsTransferFirstBotMutation(m)
	assert.False(t, match)
}

func TestQueryFromPhoenix(t *testing.T) {
	m := `mutation TransferTeamToPlayer($timezoneIdx: Int!, $countryIdx: ID!, $address: String!) {  transferFirstBotToAddr(  timezone: $timezoneIdx  countryIdxInTimezone: $countryIdx  address: $address  )  }`
	match, _ := authproxy.IsTransferFirstBotMutation(m)
	assert.True(t, match)
}

func TestMatchTransferFirstBotMutation(t *testing.T) {
	t.Run("nil request", func(t *testing.T) {
		_, err := authproxy.MatchTransferFirstBotMutation(nil)
		assert.EqualError(t, err, "nil request")
	})
	t.Run("request without a body", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, "http://localhost", http.NoBody)
		assert.Nil(t, err)
		result, err := authproxy.MatchTransferFirstBotMutation(request)
		assert.Nil(t, err)
		assert.False(t, result)
	})
	t.Run("body of request is not nil after read", func(t *testing.T) {
		body := []byte{0x2}
		request, err := http.NewRequest(http.MethodPost, "http://localhost", bytes.NewBuffer(body))
		assert.Nil(t, err)
		resultBody, err := ioutil.ReadAll(request.Body)
		assert.Nil(t, err)
		assert.Equal(t, body, resultBody)
		resultBody, err = ioutil.ReadAll(request.Body)
		assert.Nil(t, err)
		assert.Equal(t, resultBody, []byte{})
	})
	t.Run("request body persists", func(t *testing.T) {
		t.Skip("has to be reactivated")
		body, err := json.Marshal(map[string]string{
			"query": "mutation {transferFirstBotToAddr(timezone: 1, countryIdxInTimezone: 0, address: \"0x02\")}",
		})
		assert.Nil(t, err)
		request, err := http.NewRequest(http.MethodPost, "http://localhost", bytes.NewBuffer(body))
		assert.Nil(t, err)
		isMatch, err := authproxy.MatchTransferFirstBotMutation(request)
		assert.Nil(t, err)
		assert.True(t, isMatch)
		resultBody, err := ioutil.ReadAll(request.Body)
		assert.Nil(t, err)
		assert.Equal(t, body, resultBody)
	})
}

func TestQueryFromPhoenix1(t *testing.T) {
	m := `
query playerInfo($playerId: String!) {
  playerByPlayerId(playerId: $playerId) {
    teamId
    playerId
    name
    preferredPosition
    dayOfBirth
    countryOfBirth
    race
    endurance
    defence
    shoot
    pass
    speed
    potential
    tiredness
    redCard
    injuryMatchesLeft
    teamByTeamId {
      teamId
      owner
      name
      managerName
      gkCount: playersByTeamId(condition: { preferredPosition: "GK" }) {
        totalCount
      }
    }
    offersByPlayerId(condition: { state: STARTED }, orderBy: PRICE_DESC) {
      nodes {
        id
        teamByBuyerTeamId {
          teamId
          name
          owner
          managerName
        }
        rnd
        price
        validUntil
        auctionId
      }
    }
  }
}
	`
	match, _ := authproxy.IsTransferFirstBotMutation(m)
	assert.False(t, match)
}
