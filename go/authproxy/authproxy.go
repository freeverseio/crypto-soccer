package authproxy

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

const (
	GodToken = "joshua"
)

type AuthProxy struct {
	timeout       int
	backdoor      bool
	cache         *gocache.Cache
	gracetime     int
	debug         bool
	serverService ServerService
}

func New(
	timeout int,
	gracetime int,
	serverService ServerService,

) *AuthProxy {
	return &AuthProxy{
		timeout:   timeout,
		gracetime: gracetime,
		// create authentication cache
		// default expiration time of 5 minutes, and purges expired items every 2 minute
		cache:         gocache.New(5*time.Minute, 2*time.Minute),
		serverService: serverService,
	}
}

func (b *AuthProxy) SetDebug(active bool) {
	b.debug = active
}

func (b *AuthProxy) SetBackdoor(active bool) {
	b.backdoor = active
}

var opid uint64

func (b *AuthProxy) checkAuthorization(ctx context.Context, r *http.Request) (string, error) {
	// check if token is well formed
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, "Bearer") {
		return "", errors.New("No authorization bearer")
	}
	token := strings.TrimSpace(auth[len("Bearer"):])

	// if backdoor is activated, check if is the godmode token
	if b.backdoor && token == GodToken {
		return GodToken, nil
	}

	// check if token is cached
	if addrHex, ok := b.cache.Get(token); ok {
		// metricsCacheHits.Inc()
		return addrHex.(string), nil
	}

	// verify token
	addr, _, err := VerifyToken(token, time.Duration(b.gracetime)*time.Second)
	if err != nil {
		return "", err
	}

	countTeams, err := b.serverService.CountTeams(ctx, addr)
	if err != nil {
		return "", err
	}
	isTransferFirstBot, err := MatchTransferFirstBotMutation(r)
	if err != nil {
		log.Warningf("failed checking for the transfer first bot match %v", err)
	}
	if isTransferFirstBot {
		if countTeams != 0 {
			return "", errors.New("Already owner of a team")
		}
		return GodToken, nil
	}

	if countTeams < 1 {
		return "", errors.New("User does not has pesmissions")
	}

	// add to cache & return
	// metricsCacheMisses.Inc()
	b.cache.Set(token, addr.Hex(), gocache.DefaultExpiration)
	return addr.Hex(), nil
}

func (b *AuthProxy) Gqlproxy(w http.ResponseWriter, r *http.Request) {

	// auto increment op number
	op := atomic.AddUint64(&opid, 1)

	// define the function for http failing
	fail := func(err error) {
		// metricsOpsFailed.Inc()
		log.Error("[", op, "] ", err)
		http.Error(w, fmt.Sprintf("Internal error traceid:%v", op), http.StatusInternalServerError)
	}

	// set the maximum time for the whole operation
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(b.timeout)*time.Second)
	defer cancel()

	// check if this request can be done
	addr, err := b.checkAuthorization(ctx, r)
	if err != nil {
		// metricsOpsFailed.Inc()
		log.Error("[", op, "]", err)
		http.Error(w, fmt.Sprintf("Invalid authorization token [%v]", op), http.StatusUnauthorized)
		return
	}

	// send the request to the gql
	var request io.Reader
	if b.debug {
		op = atomic.AddUint64(&opid, 1)
		requestBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fail(err)
			return
		}
		log.Debug("> [", op, "] '", string(requestBytes), "'")
		request = bytes.NewReader(requestBytes)
	} else {
		request = r.Body
	}

	req, err := b.serverService.NewRequest("POST", request)
	if err != nil {
		fail(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "AuthProxy")
	req.Header.Set("X-AuthProxy-Address", addr)

	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		fail(err)
		return
	}

	defer resp.Body.Close()

	// read the response from the gql
	var response io.Reader
	if b.debug {
		responseBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fail(err)
			return
		}
		log.Debug("< [", op, "] '", string(responseBytes), "'")
		response = bytes.NewReader(responseBytes)
	} else {
		response = resp.Body
	}

	// send back the reponse
	if _, err := io.Copy(w, response); err != nil {
		fail(err)
		return
	}

	// update successfull ops metric
	// metricsOpsSuccess.Inc()
}

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
		return common.Address{}, time.Time{}, errors.New("malformed token")
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
