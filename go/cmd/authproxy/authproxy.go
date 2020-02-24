package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/didip/tollbooth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/authproxy"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsOpsSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authproxy_ops_success",
		Help: "The total number of processed events",
	})
	metricsOpsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authproxy_ops_failed",
		Help: "The total number of failed events",
	})
	metricsOpsDropped = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authproxy_ops_dropped",
		Help: "The total number of droped events",
	})
	metricsCacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authproxy_cache_hits",
		Help: "The total number of cache hits",
	})
	metricsCacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authproxy_cache_misses",
		Help: "The total number of cache misses",
	})
)

const (
	godtoken = "joshua"
)

var timeout *int
var gqlurl *string
var debug *bool
var backdoor *bool
var cache *gocache.Cache
var gracetime *int

var opid uint64

func checkPermissions(ctx context.Context, addr common.Address) (bool, error) {

	// create request
	gqlQuery := `{allTeams (condition: {owner: "` + addr.Hex() + `"}){totalCount}}`
	query, err := json.Marshal(map[string]string{"query": gqlQuery})
	if err != nil {
		return false, errors.Wrap(err, "failed bulding auth query")
	}
	req, err := http.NewRequest(http.MethodPost, *gqlurl, bytes.NewReader(query))
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

func checkAuthorization(ctx context.Context, r *http.Request) (string, error) {

	// check if token is well formed
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, "Bearer") {
		return "", errors.New("No authorization bearer")
	}
	token := strings.TrimSpace(auth[len("Bearer"):])

	// if backdoor is activated, check if is the godmode token
	if *backdoor && token == godtoken {
		return godtoken, nil
	}

	// check if token is cached
	if addrHex, ok := cache.Get(token); ok {
		metricsCacheHits.Inc()
		return addrHex.(string), nil
	}

	// verify token
	addr, _, err := authproxy.VerifyToken(token, time.Duration(*gracetime)*time.Second)
	if err != nil {
		return "", err
	}

	ok, err := checkPermissions(ctx, addr)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("User does not has pesmissions")
	}

	// add to cache & return
	metricsCacheMisses.Inc()
	cache.Set(token, addr.Hex(), gocache.DefaultExpiration)
	return addr.Hex(), nil
}

func gqlproxy(w http.ResponseWriter, r *http.Request) {

	// auto increment op number
	op := atomic.AddUint64(&opid, 1)

	// define the function for http failing
	fail := func(err error) {
		metricsOpsFailed.Inc()
		log.Error("[", op, "] ", err)
		http.Error(w, fmt.Sprintf("Internal error traceid:%v", op), http.StatusInternalServerError)
	}

	// set the maximum time for the whole operation
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	// check if this request can be done
	addr, err := checkAuthorization(ctx, r)
	if err != nil {
		metricsOpsFailed.Inc()
		log.Error("[", op, "]", err)
		http.Error(w, fmt.Sprintf("Invalid authorization token [%v]", op), http.StatusUnauthorized)
		return
	}

	// send the request to the gql
	var request io.Reader
	if *debug {
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

	req, err := http.NewRequest("POST", *gqlurl, request)
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
	if *debug {
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
	metricsOpsSuccess.Inc()
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func startProxyServer(serviceurl *string, ratelimit *int) {

	// check gql server
	_, err := checkPermissions(context.Background(), common.HexToAddress("0x83A909262608c650BD9b0ae06E29D90D0F67aC5e"))
	if err != nil {
		log.Warnf("Caution, cannot access to gql service at %s: %w", gqlurl, err)
	}

	// parse service url
	serviceURL, err := url.Parse(*serviceurl)
	if err != nil {
		log.Fatal("malformed URL ", *serviceurl)
	}

	// create authentication cache
	// default expiration time of 5 minutes, and purges expired items every 2 minute
	cache = gocache.New(5*time.Minute, 2*time.Minute)

	// create server handing requests with reqest limits
	proxyserver := http.NewServeMux()

	lmt := tollbooth.NewLimiter(float64(*ratelimit), nil)
	lmt.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		metricsOpsDropped.Inc()
	})
	lmtHandler := tollbooth.LimitFuncHandler(lmt, gqlproxy)
	proxyserver.Handle(serviceURL.Path, lmtHandler)

	// start the server
	if serviceURL.Scheme == "https" {
		log.Info("Starting at :443")
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(serviceURL.Host),
			Cache:      autocert.DirCache("cache-path"),
		}
		server := &http.Server{
			Handler:   proxyserver,
			Addr:      ":https",
			TLSConfig: m.TLSConfig(),
		}
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatal("cannot start https server: ", err)
		}
	} else if serviceURL.Scheme == "http" {
		log.Info("Starting proxy server at :8080")
		server := &http.Server{
			Handler: proxyserver,
			Addr:    ":8080",
		}
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("cannot start http server", err)
		}
	} else {
		log.Fatal("Unknown scheme ", serviceURL.Scheme)
	}
}

func startMetricsServer() {
	log.Info("Starting metrics server at :4000/metrics")
	metricsserver := http.NewServeMux()
	metricsserver.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Handler: metricsserver,
		Addr:    ":4000",
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("cannot start metrics server", err)
	}
}

func main() {

	gqlurl = flag.String("gqlurl", "http://dev1.gorengine.com:4000/graphql", "graphql url")
	serviceurl := flag.String("serviceurl", "http://localhost:8080/graphql", "service url, http or https for autocert")
	debug = flag.Bool("debug", false, "debug")
	timeout = flag.Int("timeout", 5, "max timeout")
	ratelimit := flag.Int("ratelimit", 1000000, "max requests per second")
	backdoor = flag.Bool("backdoor", false, " allow god mode for token 'hi!'")
	gracetime = flag.Int("gracetime", 3600, " grace time for tickets in seconds")
	flag.Parse()

	log.Info("-timeout=", *timeout)
	log.Info("-gqlurl=", *gqlurl)
	log.Info("-serviceurl=", *serviceurl)
	log.Info("-debug=", *debug)
	log.Info("-ratelimit=", *ratelimit)
	log.Info("-backdoor=", *backdoor, " (Bearer ", godtoken, ")")
	log.Info("-gracetime=", *gracetime)

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	go startMetricsServer()
	startProxyServer(serviceurl, ratelimit)
}

/*
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer joshua" --data '{"query":"{allTeams (condition: {owner: \"0x83A909262608c650BD9b0ae06E29D90D0F67aC5e\"}){totalCount}}"}'  `, *serviceurl
ab -n 1000 -c 100 -p data.json -H "Authorization: Bearer joshua"
*/
