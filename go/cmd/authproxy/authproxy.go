package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/freeverseio/crypto-soccer/go/authproxy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// var (
// 	metricsOpsSuccess = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "authproxy_ops_success",
// 		Help: "The total number of processed events",
// 	})
// 	metricsOpsFailed = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "authproxy_ops_failed",
// 		Help: "The total number of failed events",
// 	})
// 	metricsOpsDropped = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "authproxy_ops_dropped",
// 		Help: "The total number of droped events",
// 	})
// 	metricsCacheHits = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "authproxy_cache_hits",
// 		Help: "The total number of cache hits",
// 	})
// 	metricsCacheMisses = promauto.NewCounter(prometheus.CounterOpts{
// 		Name: "authproxy_cache_misses",
// 		Help: "The total number of cache misses",
// 	})
// )

func main() {

	gqlurl := flag.String("gqlurl", "http://dev1.gorengine.com:4000/graphql", "graphql url")
	serviceport := flag.Int("serviceport", 8080, "service port")
	metricsport := flag.Int("metricsport", 4000, "metrics port")
	debug := flag.Bool("debug", false, "debug")
	timeout := flag.Int("timeout", 5, "max timeout")
	ratelimit := flag.Int("ratelimit", 1000000, "max requests per second")
	backdoor := flag.Bool("backdoor", false, " allow god mode for token 'hi!'")
	gracetime := flag.Int("gracetime", 3600, " grace time for tickets in seconds")
	flag.Parse()

	log.Info("-timeout=", *timeout)
	log.Info("-gqlurl=", *gqlurl)
	log.Info("-serviceport=", *serviceport)
	log.Info("-metricsport=", *metricsport)
	log.Info("-debug=", *debug)
	log.Info("-ratelimit=", *ratelimit)
	log.Info("-backdoor=", *backdoor, " (Bearer ", authproxy.GodToken, ")")
	log.Info("-gracetime=", *gracetime)

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	ap := authproxy.New(
		*gqlurl,
		*timeout,
		*gracetime,
	)
	ap.SetDebug(*debug)
	ap.SetBackdoor(*backdoor)

	// go startMetricsServer(*metricsport)
	ap.StartProxyServer(*serviceport, *ratelimit)
}

func startMetricsServer(metricsport int) {

	bind := fmt.Sprintf(":%v", metricsport)
	log.Infof("Starting metrics server at %v/metrics", bind)

	metricsserver := http.NewServeMux()
	metricsserver.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Handler: metricsserver,
		Addr:    bind,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("cannot start metrics server", err)
	}
}

/*
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer joshua" --data '{"query":"{allTeams (condition: {owner: \"0x83A909262608c650BD9b0ae06E29D90D0F67aC5e\"}){totalCount}}"}'  `, *serviceurl
ab -n 1000 -c 100 -p data.json -H "Authorization: Bearer joshua"
*/
