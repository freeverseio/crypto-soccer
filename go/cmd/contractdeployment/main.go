package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/freeverseio/crypto-soccer/go/testutils"
)

type ConfigMap struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Data struct {
		Freeverse_username               string `yaml:"freeverse_username"`
		Freeverse_password               string `yaml:"freeverse_password"`
		Assets_contract_address          string `yaml:"assets_contract_address"`
		Leagues_contract_address         string `yaml:"leagues_contract_address"`
		Updates_contract_address         string `yaml:"updates_contract_address"`
		Engine_contract_address          string `yaml:"engine_contract_address"`
		Market_contract_address          string `yaml:"market_contract_address"`
		Evolution_contract_address       string `yaml:"evolution_contract_address"`
		Engineprecomp_contract_address   string `yaml:"engineprecomp_contract_address"`
		Matchevents_contract_address     string `yaml:"matchevents_contract_address"`
		Utils_match_log_contract_address string `yaml:"utils_match_log_contract_address"`
		Enode                            string `yaml:"enode"`
	} `yaml:"data"`
}

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

func main() {
	b, err := testutils.NewBlockchainNodeDeployAndInitAt("http://ethereum:8545")
	check(err)

	c := ConfigMap{}
	c.APIVersion = "v1"
	c.Kind = "ConfigMap"
	c.Metadata.Name = "freeverse-configmap"
	c.Data.Freeverse_username = "freeverse"
	c.Data.Freeverse_password = "freeverse"
	c.Data.Enode = "enode://133f77f423d96282613afe4a3bd2c09a0645be853bd8d27d75da3064b1692cfc869ddeca586dc7969cfa4a30b9dbc9856f5cb02bd20fcb5fc0697c2b1fe2ce46@165.22.66.118:30303"
	c.Data.Assets_contract_address = b.Addresses.Assets
	c.Data.Leagues_contract_address = b.Addresses.Leagues
	c.Data.Updates_contract_address = b.Addresses.Updates
	c.Data.Engine_contract_address = b.Addresses.Engine
	c.Data.Market_contract_address = b.Addresses.Market
	c.Data.Evolution_contract_address = b.Addresses.Evolution
	c.Data.Engineprecomp_contract_address = b.Addresses.Engineprecomp
	c.Data.Matchevents_contract_address = b.Addresses.Matchevents
	c.Data.Utils_match_log_contract_address = b.Addresses.Utilsmatchlog

	data, err := yaml.Marshal(&c)
	check(err)
	fmt.Printf("\n\n%v\n", string(data))

	if len(os.Args) > 1 {
		outputfilename := os.Args[1]
		fmt.Printf("writing configmap to yaml file: %v\n", outputfilename)
		err := ioutil.WriteFile(outputfilename, data, 0644)
		check(err)
	}
}
