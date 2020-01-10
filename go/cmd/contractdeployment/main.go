package main

import (
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func main() {
	_, err := testutils.NewBlockchainNodeDeployAndInitAt("http://ethereum:8545")
	if err != nil {
		panic(err)
	}
}
