package main

import (
	banner "github.com/CrowdSurge/banner"
	cmd "github.com/freeverseio/go-soccer/cmd"
)

func main() {

	banner.Print("go-soccer")
	cmd.ExecuteCmd()
}
