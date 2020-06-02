package authproxy

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

const (
	GodToken = "joshua"
)

func hashPrefixedMessage(data []byte) []byte {
	prefixed := []byte{0x19}
	prefixed = append(prefixed, []byte(fmt.Sprintf("Ethereum Signed Message:\n%v", len(data)))...)
	prefixed = append(prefixed, data...)
	return prefixed
}
func SignToken(privateKey *ecdsa.PrivateKey, t time.Time) (string, error) {
	ts := fmt.Sprintf("%v", t.Unix())
	hash := crypto.Keccak256(hashPrefixedMessage([]byte(ts)))
	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return "", fmt.Errorf("unable to sign token %w", err)
	}
	token := ts + ":" + base64.StdEncoding.EncodeToString(sig)
	return token, nil
}
func VerifyToken(token string, grace time.Duration) (common.Address, time.Time, error) {
	// token have: unix_time:signature
	tokenfields := strings.Split(token, ":")
	if len(tokenfields) != 2 {
		return common.Address{}, time.Time{}, errors.New("malformed token:" + token)
	}
	// check date
	tsunix, err := strconv.Atoi(tokenfields[0])
	if err != nil {
		return common.Address{}, time.Time{}, fmt.Errorf("malformed token timestamp %w", err)
	}
	ts, now := time.Unix(int64(tsunix), 0), time.Now()
	if math.Abs(now.Sub(ts).Seconds()) > grace.Seconds() {
		return common.Address{}, time.Time{}, errors.New("token out of time")
	}
	// retrieve public key
	hash := crypto.Keccak256(hashPrefixedMessage([]byte(tokenfields[0])))
	sig, err := base64.StdEncoding.DecodeString(tokenfields[1])
	if err != nil {
		return common.Address{}, time.Time{}, fmt.Errorf("malformed token signature encoding : %w", err)
	}
	pbk, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return common.Address{}, time.Time{}, fmt.Errorf("malformed token signature: %w", err)
	}
	addr := crypto.PubkeyToAddress(*pbk)
	return addr, ts, nil
}

func MatchTransferFirstBotMutation(data string) (bool, error) {
	ex := `mutation(\s*).*(\s*){(\s*)transferFirstBotToAddr(\s*)\((\s*)timezone(\s*):(\s*)\d{1,2}(\s*),(\s*)countryIdxInTimezone(\s*):(\s*)[0-9]+(\s*),(\s*)address(\s*):(\s*)"[a-zA-Z0-9]+"(\s*)\)(\s*)}`
	return regexp.MatchString(ex, data)
}

func matchTransferFirstBotMutation(r *http.Request) (bool, error) {
	if r.Body == nil {
		return false, errors.New("nil body")
	}
	var query struct {
		Data string `json:"query"`
	}
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		return false, err
	}
	return MatchTransferFirstBotMutation(query.Data)
}

func CheckAuthorization(
	ctx context.Context,
	r *http.Request,
	backdoor bool,
	cache *gocache.Cache,
	gracetime int,
	gqlurl string,
) (string, error) {
	if r == nil {
		return "", errors.New("r is nil")
	}
	// check if token is well formed
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, "Bearer") {
		return "", errors.New("No authorization bearer")
	}
	token := strings.TrimSpace(auth[len("Bearer"):])

	// if backdoor is activated, check if is the godmode token
	if backdoor && token == GodToken {
		return GodToken, nil
	}

	match, err := matchTransferFirstBotMutation(r)
	if match {
		// Always ALLOW this query
		return GodToken, nil
	}
	if err != nil {
		return "", err
	}

	// check if token is cached
	if addrHex, ok := cache.Get(token); ok {
		// metricsCacheHits.Inc()
		return addrHex.(string), nil
	}

	// verify token
	addr, _, err := VerifyToken(token, time.Duration(gracetime)*time.Second)
	if err != nil {
		return "", err
	}

	ok, err := CheckPermissions(ctx, addr, gqlurl)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("User does not has pesmissions")
	}

	// add to cache & return
	// metricsCacheMisses.Inc()
	cache.Set(token, addr.Hex(), gocache.DefaultExpiration)
	return addr.Hex(), nil
}

func CheckPermissions(ctx context.Context, addr common.Address, gqlurl string) (bool, error) {

	// create request
	gqlQuery := `{allTeams (condition: {owner: "` + addr.Hex() + `"}){totalCount}}`
	query, err := json.Marshal(map[string]string{"query": gqlQuery})
	if err != nil {
		return false, errors.Wrap(err, "failed bulding auth query")
	}
	req, err := http.NewRequest(http.MethodPost, gqlurl, bytes.NewReader(query))
	if err != nil {
		return false, errors.Wrap(err, "failed bulding auth request")
	}
	req.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return false, errors.Wrap(err, "failed sending auth request")
	}

	// check http response is ok
	if resp.StatusCode != 200 {
		errstr := "unknown"
		if resp.Body != nil {
			errbody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errstr = string(errbody)
			}
		}
		return false, fmt.Errorf("failed sending auth request, errcode=%v, err=%s", resp.StatusCode, errstr)
	}

	// parse qgl response, and return
	var response struct {
		Data struct {
			AllTeams struct {
				TotalCount int
			}
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, errors.Wrap(err, "failed decoding auth response")
	}

	ok := response.Data.AllTeams.TotalCount > 0
	return ok, nil
}
