package config

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// C is the package config
var C Config

// Config is the server configurtion
type Config struct {
	Contracts struct {
		StakersAddress string
		AssetsAddress  string
		LeaguesAddress string
		StateAddress   string
	}

	Stakers struct {
		Accounts []struct {
			Address   string
			OnionHash string
		}
		Keystore  string
		Keypasswd string
	}

	DB struct {
		Path string
	}

	Web3 struct {
		MaxGasPrice uint64
		RPCURL      string
	}

	API struct {
		Port int
	}
}

func MustRead(c *cli.Context) error {
	configfile := strings.TrimSuffix(c.GlobalString("config"), ".yaml")
	return ReadConfig(configfile)
}

func ReadConfig(configfile string) error {
	if len(configfile) == 0 {
		configfile = "config"
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(configfile)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		return err
	}
	return nil
}
