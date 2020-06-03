package authproxy_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/authproxy"
	"github.com/stretchr/testify/assert"
)

func TestMatchTransferFirstBotMutation(t *testing.T) {
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
