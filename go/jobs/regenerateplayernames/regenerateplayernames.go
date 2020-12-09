package regenerateplayernames

import (
	"fmt"
	"os"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
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

func NewJob(universeDB *string, universeURL *string, gameURL *string, debug *bool) {

	log.Infof("[PARAM] universedb                 : %v", *universeDB)
	log.Infof("[PARAM] universe                   : %v", *universeURL)
	log.Infof("[PARAM] game                       : %v", *gameURL)
	log.Infof("[PARAM] debug                      : %v", *debug)
	log.Infof("-------------------------------------------------------------------")

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := func() error {
		start := time.Now()
		log.Info("Create the clients to GQL APIs")
		game := graphql.NewClient(*gameURL, nil)
		universe := graphql.NewClient(*universeURL, nil)

		log.Info("Create the db client")
		universedb, err := postgres.New(*universeDB)
		if err != nil {
			return err
		}
		defer universedb.Close()
		playerService := NewPlayerService(universedb, universe, game)

		log.Info("Retrieving all players")
		players, err := playerService.getAllPlayers()
		if err != nil {
			return err
		}

		log.Info("Updating all players names and races")
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
