package regenerateplayernames

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/shurcooL/graphql"
	log "github.com/sirupsen/logrus"
)

type FillGameDb struct {
	universeDB  string
	universeURL string
	gameURL     string
	debug       *bool
}

func New(
	universeDB string,
	universeURL string,
	gameURL string,
	debug *bool,

) *FillGameDb {
	return &FillGameDb{
		universeDB:  universeDB,
		universeURL: universeURL,
		gameURL:     gameURL,
		debug:       debug,
	}
}

func NewJob(universeDB *string, universeURL *string, debug *bool) {

	log.Infof("[PARAM] universedb                 : %v", *universeDB)
	log.Infof("[PARAM] universe                   : %v", *universeURL)
	log.Infof("[PARAM] debug                      : %v", *debug)
	log.Infof("-------------------------------------------------------------------")

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := func() error {
		start := time.Now()
		log.Info("Create the clients to GQL APIs")
		universe := graphql.NewClient(*universeURL, nil)

		log.Info("Create the db client")

		universedb, err := NewDB(*universeDB)
		if err != nil {
			log.Info("The error %v", err)
			return err
		}
		defer universedb.Close()
		playerService := NewPlayerService(universedb, universe)

		log.Info("Retrieving all players")
		players, err := playerService.getAllPlayers()
		if err != nil {
			return err
		}

		log.Info("Updating all players names and races for tz 7,8,9,11")
		err = playerService.updatetAllPlayerNamesAndRaces(players)
		if err != nil {
			return err
		}

		elapsed := time.Since(start)
		fmt.Printf("Exectuion took %s\n", elapsed)
		return nil
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func NewDB(url string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("Postgres connect err: (%v)", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Postgres Ping err: (%v)", err)
	}
	return db, nil
}
