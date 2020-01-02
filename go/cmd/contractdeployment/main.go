package main

import (
	"fmt"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func main() {
	b, err := testutils.NewBlockchainNodeDeployAndInitAt("http://ethereum:8545")
	if err != nil {
		panic(err)
	}
	fmt.Println("Copy and paste the following data into configmap.yaml")
	fmt.Println("data:")
	fmt.Println("\tfreeverse_username: freeverse")
	fmt.Println("\tfreeverse_password: freeverse")
	fmt.Printf("\tassets_contract_address: %v\n", b.Addresses.Assets)
	fmt.Printf("\tleagues_contract_address: %v\n", b.Addresses.Leagues)
	fmt.Printf("\tupdates_contract_address: %v\n", b.Addresses.Updates)
	fmt.Printf("\tengine_contract_address: %v\n", b.Addresses.Engine)
	fmt.Printf("\tmarket_contract_address: %v\n", b.Addresses.Market)
	fmt.Printf("\tevolution_contract_address: %v\n", b.Addresses.Evolution)
	fmt.Printf("\tengineprecomp_contract_address: %v\n", b.Addresses.Engineprecomp)
	fmt.Printf("\tmatchevents_contract_address: %v\n", b.Addresses.Matchevents)
	fmt.Printf("\tutils_match_log_contract_address: %v\n", b.Addresses.Utilsmatchlog)
	fmt.Println("\tenode: enode://133f77f423d96282613afe4a3bd2c09a0645be853bd8d27d75da3064b1692cfc869ddeca586dc7969cfa4a30b9dbc9856f5cb02bd20fcb5fc0697c2b1fe2ce46@165.22.66.118:30303")
}
