package authproxy_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/authproxy"
	"github.com/stretchr/testify/assert"
)

const (
	pvkhex = "c6d398e89bf7cbda7663ca881bd992eb80ad170e4ca0bd65a8b1c719ee02bc67"
)

func TestBasic(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)
	addr := crypto.PubkeyToAddress(pvk.PublicKey)

	token, err := authproxy.SignToken(pvk, time.Now())
	assert.Nil(t, err)
	addr1, _, err := authproxy.VerifyToken(token, time.Hour)
	assert.Nil(t, err)

	assert.EqualValues(t, addr, addr1)
}

func TestHalfGraceUp(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := authproxy.SignToken(pvk, time.Now().Add(time.Hour/2))
	assert.Nil(t, err)
	_, _, err = authproxy.VerifyToken(token, time.Hour)
	assert.Nil(t, err)
}

func TestHalfGraceDown(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := authproxy.SignToken(pvk, time.Now().Add(-time.Hour/2))
	assert.Nil(t, err)
	_, _, err = authproxy.VerifyToken(token, time.Hour)
	assert.Nil(t, err)
}

func TestDobuleGraceUp(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := authproxy.SignToken(pvk, time.Now().Add(time.Hour*2))
	assert.Nil(t, err)
	_, _, err = authproxy.VerifyToken(token, time.Hour)
	assert.NotNil(t, err)
}

func TestDobuleGraceDown(t *testing.T) {
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := authproxy.SignToken(pvk, time.Now().Add(-time.Hour*2))
	assert.Nil(t, err)
	_, _, err = authproxy.VerifyToken(token, time.Hour)
	assert.NotNil(t, err)
}

func BenchmarkSignatureVerification(b *testing.B) {
	pvkhex := "c6d398e89bf7cbda7663ca881bd992eb80ad170e4ca0bd65a8b1c719ee02bc67"
	pvk, _ := crypto.HexToECDSA(pvkhex)

	token, err := authproxy.SignToken(pvk, time.Now())
	assert.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		authproxy.VerifyToken(token, time.Hour)
	}
}

func TestPhoenixAuthIsValid(t *testing.T) {
	address := common.HexToAddress("12890D2cce102216644c59daE5baed380d84830c")
	token := "1579870008:Wic7VnPY+w4HwIGDvxpoDhXUyFBEJodyGsIOoN7iKmYZXa+lk3Zji5EfJtabXGnNblQOUl9bLPsFITSBqfNWawA="

	addr, _, err := authproxy.VerifyToken(token, time.Hour*10000)

	assert.Nil(t, err)
	assert.EqualValues(t, address, addr)
}

func TestMatchTransferFirstBotMutation(t *testing.T) {
	m := "mutation {transferFirstBotToAddr(timezone: 10, countryIdxInTimezone: 1000, address: \"0x02\")}"
	match, _ := authproxy.MatchTransferFirstBotMutation(m)
	assert.True(t, match)

	// spaces should not affect matching
	m = "mutation  {    transferFirstBotToAddr   (timezone    :   10 ,  countryIdxInTimezone :   10 ,    address: \"0x02\"   ) }"
	match, _ = authproxy.MatchTransferFirstBotMutation(m)
	assert.True(t, match)

	// wrong param order should fail
	m = "mutation {transferFirstBotToAddr(countryIdxInTimezone: 1000, timezone: 10, address: \"0x02\")}"
	match, _ = authproxy.MatchTransferFirstBotMutation(m)
	assert.False(t, match)

	// wrong param name should fail (Timezone should be timezone)
	m = "mutation {transferFirstBotToAddr(Timezone: 10, countryIdxInTimezone: 1000, address: \"0x02\")}"
	match, _ = authproxy.MatchTransferFirstBotMutation(m)
	assert.False(t, match)

	// wrong method name should fail
	m = "mutation {foo(Timezone: 10, countryIdxInTimezone: 1000, address: \"0x02\")}"
	match, _ = authproxy.MatchTransferFirstBotMutation(m)
	assert.False(t, match)

	// not using mutation should fail
	m = "bar {transferFirstBotToAddr(timezone: 10, countryIdxInTimezone: 1000, address: \"0x02\")}"
	match, _ = authproxy.MatchTransferFirstBotMutation(m)
	assert.False(t, match)
}

func TestCheckAuthorizationGodToken(t *testing.T) {
	backdoor := true
	gracetime := 10
	gqlurl := ""
	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+authproxy.GodToken)
	s, err := authproxy.CheckAuthorization(
		context.Background(),
		req,
		backdoor,
		nil,
		gracetime,
		gqlurl,
	)
	assert.Nil(t, err)
	assert.Equal(t, s, authproxy.GodToken)

	backdoor = false
	_, err = authproxy.CheckAuthorization(
		context.Background(),
		req,
		backdoor,
		nil,
		gracetime,
		gqlurl,
	)
	assert.EqualError(t, err, "nil body")
}
