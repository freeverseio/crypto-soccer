package cmd

import (
	cfg "github.com/freeverseio/go-soccer/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// cfgFile is the configuration file path.
	cfgFile string
	// verbose is the verbosity level used in logrus.
	verbose string
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if logLevel, err := log.ParseLevel(verbose); err == nil {
		log.SetLevel(logLevel)
	} else {
		panic(err)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("gosoccer") // name ofconfig file (without extension)
	viper.AddConfigPath(".")        // adding current directory as first search path
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.SetEnvPrefix("GOSO")      // so viper.AutomaticEnv will get matching envvars starting with O2M_
	viper.AutomaticEnv()            // read in environment variables that match

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	log.WithField("file", viper.ConfigFileUsed()).Debug("Using config file")

	if err := viper.Unmarshal(&cfg.C); err != nil {
		panic(err)
	}

}
