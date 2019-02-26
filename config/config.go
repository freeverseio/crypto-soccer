package config

import (
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// C is the package config
var C Config

// Config is the server configurtion
type Config struct {
	Contracts struct {
		LionelAddress   string
		StackersAddress string
	}

	Keystore struct {
		Account string
		Path    string
		Passwd  string
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
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
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
