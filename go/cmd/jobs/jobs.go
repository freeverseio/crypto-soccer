package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/jobs/fillgamedb"
	log "github.com/sirupsen/logrus"
)

func main() {
	jobName := os.Getenv("JOB_NAME")
	universeURL := os.Getenv("UNIVERSE_URL")
	gameURL := os.Getenv("GAME_URL")
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	log.Infof("[PARAM] job name                   : %v", jobName)
	log.Infof("[PARAM] debug                      : %v", debug)
	log.Infof("[PARAM] universe                   : %v", universeURL)
	log.Infof("[PARAM] game                       : %v", gameURL)
	log.Infof("[PARAM] debug                      : %v", debug)
	log.Infof("-------------------------------------------------------------------")

	if debug {
		log.SetLevel(log.DebugLevel)
	}
	if err := func() error {
		switch jobName {
		case "fillgamedb":
			fillgamedb.NewJob(&universeURL, &gameURL, &debug)
			return nil
		default:
			fmt.Println("Job does not exist ", jobName)
		}

		return nil
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
