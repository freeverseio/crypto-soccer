package fillgamedb

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/shurcooL/graphql"
	log "github.com/sirupsen/logrus"
)

func main() {
	universeURL := flag.String("universeApi", "http://universe.api:4000/graphql", "gql universe url")
	gameURL := flag.String("gameApi", "http://game.api:4000/graphql", "gql game url")
	debug := flag.Bool("debug", false, "print debug logs")
	flag.Parse()

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
		universe := graphql.NewClient("http://localhost:4000/graphql", nil)
		game := graphql.NewClient("http://localhost:4040/graphql", nil)

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
