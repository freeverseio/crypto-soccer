package fillgamedb

import (
	"fmt"
	"os"
	"time"

	"github.com/shurcooL/graphql"
	log "github.com/sirupsen/logrus"
)

type FillGameDb struct {
	universeURL string
	gameURL     string
	debug       *bool
}

func New(
	universeURL string,
	gameURL string,
	debug *bool,

) *FillGameDb {
	return &FillGameDb{
		universeURL: universeURL,
		gameURL:     gameURL,
		debug:       debug,
	}
}

func NewJob(universeURL *string, gameURL *string, debug *bool) {

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
		universe := graphql.NewClient(*universeURL, nil)
		game := graphql.NewClient(*gameURL, nil)

		playerService := NewPlayerService(universe, game)
		teamService := NewTeamService(universe, game)

		teams, err := teamService.getAllTeams()
		if err != nil {
			return err
		}

		players, err := playerService.getAllPlayers()
		if err != nil {
			return err
		}

		err = teamService.upsertAllTeamProps(teams)
		if err != nil {
			return err
		}

		err = playerService.upsertAllPlayerProps(players)
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
