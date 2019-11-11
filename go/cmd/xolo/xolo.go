package main

import (
	"net/http"
	"os"
	"time"

	"io/ioutil"

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
	Pools map[string]struct {
		WaitReceipt string
		Rpcs        []struct {
			URL string
		}
		Signers []struct {
			Path   string
			Passwd string
		}
	}
	waitOutQueue string
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

	xserver, err := xolo.NewServer()
	if err != nil {
		log.Fatalf("Failed to start xserver: %v", err)
	}

	for poolname, pool := range C.Pools {
		waitReceipt, err := time.ParseDuration(pool.WaitReceipt)
		if err != nil {
			log.Fatal("Invalid waitReceipt ", waitReceipt, " ", err)
		}
		if err := xserver.AddPool(poolname, waitReceipt); err != nil {
			log.Fatal("Cannot create pool ", poolname, " ", err)
		}
		for _, rpc := range pool.Rpcs {
			rpclient, err := ethclient.Dial(rpc.URL)
			if err != nil {
				log.Fatalf("Failed to connect to RPC: %v %v", rpc.URL, err)
			}
			if err := xserver.AddRpcClient(poolname, rpclient); err != nil {
				log.Fatalf("Failed to add to RPC: %v %v", rpc.URL, err)
			}
			log.Info("Added rpc ", rpc.URL, " for pool ", poolname)
		}
		for _, signerdef := range pool.Signers {
			content, err := ioutil.ReadFile(signerdef.Path)
			if err != nil {
				log.Fatalf("Unable to read key file '%v' %v", signerdef.Path, err)
			}
			signer, err := xserver.AddSigner(poolname, string(content), signerdef.Passwd)
			if err != nil {
				log.Fatalf("Unable use key %v %v", signerdef.Path, err)
			}
			log.Info("Added signer ", signer.Address().Hex(), " for pool ", poolname)
		}

	}

	engine := gin.Default()
	engine.GET("/info", func(c *gin.Context) {
		c.String(http.StatusOK, xserver.Info())
	})
	engine.POST("/tx", xserver.ServePostTx)
	engine.GET("/tx/:txhash", xserver.ServeGetTx)

	go xserver.Start(engine)
	go func() {
		for {
			log.Info(xserver.Info())
			time.Sleep(time.Second * 5)
		}
	}()
	/*
		// catch ^C to send the stop signal
		ossig := make(chan os.Signal, 1)
		signal.Notify(ossig, os.Interrupt)

		go func() {
			for sig := range ossig {
				if sig == os.Interrupt {
					serv.Stop()
					os.Exit(1) // FIX
				}
			}
		}()
	*/
	engine.Run("0.0.0.0:8004")

	/*
		serv, err := service.NewService(stkrs, storage)
		must(err)


		serv.Start()
		serv.Join()
	*/
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
