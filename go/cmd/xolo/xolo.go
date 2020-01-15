package main

import (
	"net/http"
	"os"
	"strings"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/freeverseio/crypto-soccer/go/xolo"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type Config struct {
	Rpc        struct {
		URL string
	}
	Keystore struct {
		Path   string
		Passwd string
	}
}

var C Config

func init() {
	cobra.OnInitialize(MustLoadConfig)
}

func MustLoadConfig() {

	viper.SetConfigType("yaml")
	viper.SetConfigName("xolo")
	viper.AddConfigPath(".") // adding home directory as first search path

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}
}

func must(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func serverStart(c *cli.Context) {

	MustLoadConfig()

	rpclient, err := ethclient.Dial(C.Rpc.URL)
	if err != nil {
		log.Fatalf("Failed to connect to RPC: %v %v", C.Rpc.URL, err)
	}
	keystore, err := ioutil.ReadFile(C.Keystore.Path)
	if err != nil {
		log.Fatalf("Unable to read key file '%v' %v", C.Keystore.Path, err)
	}
	signer, err := bind.NewTransactor(strings.NewReader(string(keystore)), C.Keystore.Passwd)
	if err != nil {
		log.Fatalf("Unable use key %v", err)
	}

	xserver, err := xolo.NewServer(signer, rpclient)
	if err != nil {
		log.Fatalf("Cannot create server", err)
	}

	engine := gin.Default()
	srv := &http.Server{
		Addr:    ":8004",
		Handler: engine,
	}
	engine.POST("/tx", xserver.HttpPostTx)
	srv.ListenAndServe()
}

var ServerCommands = []cli.Command{
	{
		Name:  "server",
		Usage: "manage server",
		Subcommands: []cli.Command{
			{
				Name:   "start",
				Usage:  "start the service",
				Action: serverStart,
			},
		},
	},
}

func main() {
	app := cli.NewApp()
	app.Description = "signer server"
	app.Name = "Xoloitzcuintle"
	app.Version = "0.0.1-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, ServerCommands...)
	err := app.Run(os.Args)
	if err != nil {
		color.Red(err.Error())
	}

}
