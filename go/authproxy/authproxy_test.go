package authproxy_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestA(t *testing.T) {
	timeout := 10
	gracetime := 10

	t.Run("no authorization", func(t *testing.T) {
		serverService := authproxy.MockServerService{}
		ap := authproxy.New(timeout, gracetime, serverService)

		req, err := http.NewRequest("GET", "/health-check", nil)
		assert.Nil(t, err)
		rr := httptest.NewRecorder()
		ap.Gqlproxy(rr, req)
		assert.Equal(t, 401, rr.Code)
		msg, err := ioutil.ReadAll(rr.Body)
		assert.Nil(t, err)
		assert.Equal(t, "Invalid authorization token [1]\n", string(msg))
	})

	t.Run("malformed token", func(t *testing.T) {
		serverService := authproxy.MockServerService{}
		ap := authproxy.New(timeout, gracetime, serverService)

		req, err := http.NewRequest("GET", "/health-check", nil)
		assert.Nil(t, err)
		req.Header.Set("Authorization", "Bearer abc123")
		rr := httptest.NewRecorder()
		ap.Gqlproxy(rr, req)
		assert.Equal(t, 401, rr.Code)
		msg, err := ioutil.ReadAll(rr.Body)
		assert.Nil(t, err)
		assert.Equal(t, "Invalid authorization token [2]\n", string(msg))
	})

	t.Run("no backdoor and godtoken token", func(t *testing.T) {
		serverService := authproxy.MockServerService{}
		ap := authproxy.New(timeout, gracetime, serverService)

		req, err := http.NewRequest("GET", "/health-check", nil)
		assert.Nil(t, err)
		req.Header.Set("Authorization", "Bearer "+authproxy.GodToken)
		rr := httptest.NewRecorder()
		ap.Gqlproxy(rr, req)
		assert.Equal(t, 401, rr.Code)
		msg, err := ioutil.ReadAll(rr.Body)
		assert.Nil(t, err)
		assert.Equal(t, "Invalid authorization token [3]\n", string(msg))
	})

	t.Run("backdoor and godtoken token", func(t *testing.T) {
		serverService := authproxy.MockServerService{}
		serverService.CountTeamFn = func() (int, error) { return 0, nil }
		serverService.NewRequestFn = func() (*http.Request, error) {
			return httptest.NewRequest(http.MethodGet, "/health-check", http.NoBody), nil
		}
		ap := authproxy.New(timeout, gracetime, serverService)
		ap.SetBackdoor(true)

		req, err := http.NewRequest("GET", "/health-check", nil)
		assert.Nil(t, err)
		req.Header.Set("Authorization", "Bearer "+authproxy.GodToken)
		rr := httptest.NewRecorder()
		ap.Gqlproxy(rr, req)
		assert.Equal(t, 500, rr.Code)
		msg, err := ioutil.ReadAll(rr.Body)
		assert.Nil(t, err)
		assert.Equal(t, "Internal error traceid:4\n", string(msg))
	})
}
