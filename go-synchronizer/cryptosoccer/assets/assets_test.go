package assets

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"testing"
	
	log "github.com/sirupsen/logrus"
)

func TestNewAssets(t *testing.T) {
	conn, err := ethclient.Dial("https://devnet.busyverse.com/web3")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	assets, err := NewAssets(common.HexToAddress("0x05Fdd4d2340bcA823802849c75F385561278c3aB"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}

	count, err := assets.CountTeams(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v", err)
	}

	t.Error("number of teams : ", count)
}