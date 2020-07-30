package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/shurcooL/graphql"
	log "github.com/sirupsen/logrus"
)

func main() {
	universeURL := flag.String("universeApi", "http://universe.api:4000/graphql", "gql universe url")
	gameURL := flag.String("gameApi", "http://game.api:4000/graphql", "gql game url")
	debug := flag.Bool("debug", false, "print debug logs")
	bufferSize := flag.Int("buffer_size", 10000, "size of event buffer")
	processWait := flag.Int("process_wait", 5, "secs to wait for next process")
	flag.Parse()

	log.Infof("[PARAM] universe                   : %v", *universeURL)
	log.Infof("[PARAM] game                       : %v", *gameURL)
	log.Infof("[PARAM] Buffer size                : %v", *bufferSize)
	log.Infof("[PARAM] Process wait               : %v", *processWait)
	log.Infof("[PARAM] debug                      : %v", *debug)
	log.Infof("-------------------------------------------------------------------")

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := func() error {
		log.Info("Create the connection to GQL APIs")
		// var allPlayers struct {
		// 	TotalCount int
		// 	Nodes      struct {
		// 		PlayerId graphql.String
		// 		Name     graphql.String
		// 	}
		// }
		var query struct {
			AllPlayers struct {
				TotalCount graphql.Int
				Nodes      []struct {
					PlayerId graphql.String
					Name     graphql.String
				}
			}
		}
		universe := graphql.NewClient("http://localhost:4000/graphql", nil)
		//game := graphql.NewClient("http://game.api:4000/graphql", nil)
		universe.Query(context.Background(), &query, nil)
		fmt.Println(query.AllPlayers.TotalCount)
		fmt.Println(query.AllPlayers.Nodes[0])

		return nil
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
