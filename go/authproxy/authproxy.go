package authproxy

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/didip/tollbooth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
)

const (
	GodToken = "joshua"
)

type AuthProxy struct {
	timeout   int
	gqlurl    string
	backdoor  bool
	cache     *gocache.Cache
	gracetime int
	debug     bool
}

func New(
	gqlurl string,
	timeout int,
	gracetime int,
) *AuthProxy {
	return &AuthProxy{
		timeout:   timeout,
		gqlurl:    gqlurl,
		gracetime: gracetime,
		// create authentication cache
		// default expiration time of 5 minutes, and purges expired items every 2 minute
		cache: gocache.New(5*time.Minute, 2*time.Minute),
	}
}

func (b *AuthProxy) SetDebug(active bool) {
	b.debug = active
}

func (b *AuthProxy) SetBackdoor(active bool) {
	b.backdoor = active
}

var opid uint64

func (b *AuthProxy) checkPermissions(ctx context.Context, addr common.Address) (bool, error) {

	// create request
	gqlQuery := `{allTeams (condition: {owner: "` + addr.Hex() + `"}){totalCount}}`
	query, err := json.Marshal(map[string]string{"query": gqlQuery})
	if err != nil {
		return false, errors.Wrap(err, "failed bulding auth query")
	}
	req, err := http.NewRequest(http.MethodPost, b.gqlurl, bytes.NewReader(query))
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

	isTransferFirstBot, err := MatchTransferFirstBotMutation(r)
	if err != nil {
		return "", errors.Wrap(err, "failed checking for the transfer first bot match")
	}
	if isTransferFirstBot {
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

	ok, err := b.checkPermissions(ctx, addr)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("User does not has pesmissions")
	}

	// add to cache & return
	// metricsCacheMisses.Inc()
	b.cache.Set(token, addr.Hex(), gocache.DefaultExpiration)
	return addr.Hex(), nil
}

func (b *AuthProxy) gqlproxy(w http.ResponseWriter, r *http.Request) {

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

	req, err := http.NewRequest("POST", b.gqlurl, request)
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

func (b *AuthProxy) StartProxyServer(serviceport int, ratelimit int) {

	// check gql server
	_, err := b.checkPermissions(context.Background(), common.HexToAddress("0x83A909262608c650BD9b0ae06E29D90D0F67aC5e"))
	if err != nil {
		log.Warnf("Caution, cannot access to gql service at %s: %w", b.gqlurl, err)
	}

	// create server handing requests with reqest limits
	proxyserver := http.NewServeMux()

	lmt := tollbooth.NewLimiter(float64(ratelimit), nil)
	lmt.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		// metricsOpsDropped.Inc()
	})
	lmtHandler := tollbooth.LimitFuncHandler(lmt, b.gqlproxy)
	proxyserver.Handle("/", lmtHandler)

	bind := fmt.Sprintf(":%v", serviceport)
	log.Infof("Starting proxy server at %v/", bind)
	server := &http.Server{
		Handler: proxyserver,
		Addr:    bind,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("cannot start http server", err)
	}
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
