package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"time"

	"github.com/freeverseio/crypto-soccer/go/purchasevoider"
	"github.com/freeverseio/crypto-soccer/go/purchasevoider/google"
	"github.com/freeverseio/crypto-soccer/go/purchasevoider/postgres"

	log "github.com/sirupsen/logrus"
)

func main() {
	debug := flag.Bool("debug", false, "print debug logs")
	universeURL := flag.String("universe_url", "postgres://freeverse:freeverse@crypto-soccer_devcontainer_dockerhost_1:5432/cryptosoccer?sslmode=disable", "postgres url")
	marketURL := flag.String("market_url", "postgres://freeverse:freeverse@crypto-soccer_devcontainer_dockerhost_1:5432/market?sslmode=disable", "postgres url")
	googleKey := flag.String("google_key", "", "google credentials")
	packageName := flag.String("package_name", "", "packege name to scan")
	periodSec := flag.Int64("period", 10, "period")

	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		log.Infof("[param] %v : %v", f.Name, f.Value)
	})

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := func() error {
		log.Infof("connecting to universe DBMS %v", *universeURL)
		universedb, err := sql.Open("postgres", *universeURL)
		if err != nil {
			return err
		}
		defer universedb.Close()

		log.Infof("connecting to market DBMS %v", *marketURL)
		marketdb, err := sql.Open("postgres", *marketURL)
		if err != nil {
			return err
		}
		defer marketdb.Close()

		googleCredentials, err := ioutil.ReadFile(*googleKey)
		if err != nil {
			return err
		}
		pvService, err := google.NewVoidPurchaseService(googleCredentials, *packageName)
		if err != nil {
			return err
		}

		processor, err := purchasevoider.New(
			pvService,
			&postgres.UniverseService{universedb},
			&postgres.MarketService{marketdb},
		)
		if err != nil {
			return err
		}

		sleepDuration := time.Duration(*periodSec) * time.Second

		log.Info("start ...")
		for {
			log.Infof("sleep for %v seconds", sleepDuration.Seconds())
			time.Sleep(sleepDuration)
			if err != processor.Run() {
				return err
			}
		}
	}(); err != nil {
		log.Fatal(err)
	}
}
