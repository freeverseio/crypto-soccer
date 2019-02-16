package commands

import (
	log "github.com/sirupsen/logrus"
)

func must(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
