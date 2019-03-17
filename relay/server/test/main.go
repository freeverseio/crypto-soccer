package main

import (
	"github.com/freeverseio/go-soccer/relay/server"
)

func main() {
	server := server.Server{}
	server.Start()
}
