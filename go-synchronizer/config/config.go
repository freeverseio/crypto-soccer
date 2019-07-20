package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	EthereumClient         string `json:"ethereumClient`
	AssetsContractAddress  string `json:"assetsContractAddress`
	StatesContractAddress  string `json:"statesContractAddress`
	LeaguesContractAddress string `json:"leaguesContractAddress`
}

func New(configFile string) (*Config, error) {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (b *Config) Print() {
	log.Info("(CONFIG) ethereumClient          : ", b.EthereumClient)
	log.Info("(CONFIG) assetsContractAddress   : ", b.AssetsContractAddress)
	log.Info("(CONFIG) statesContractAddress   : ", b.StatesContractAddress)
	log.Info("(CONFIG) leaguesContractAddress  : ", b.LeaguesContractAddress)
}
