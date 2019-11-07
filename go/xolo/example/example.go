package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"

	"github.com/freeverseio/crypto-soccer/go/xolo"
	"github.com/freeverseio/crypto-soccer/go/xolo/example/abigen"
)

const key = `{"address":"be3a732e058fdfdb3457ba1bb1d87f9c200982f2","crypto":{"cipher":"aes-128-ctr","ciphertext":"8fe134c7059aebde9043f6454a8b6451d52b3d4e4c9162728fd35f1fee05c229","cipherparams":{"iv":"1eb6f75e10f5de357f9019176fd9a8d7"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"8f42b2a446065308d01778c16af2cea37f0829096de4796c0d81b1ed140817b5"},"mac":"227761a06aa2881ee51a7fcd7cb2f35317dbbc1ba9b4749e562c7eb045650ba2"},"id":"68c3a22d-aa36-49cb-a2e9-b97c717f3958","version":3}`

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	rpcURL := "https://goerli.infura.io/v3/9c24cba6b7b647d28deb48817cf605ce"
	rpclient, err := ethclient.Dial(rpcURL)
	assert(err)

	xserver, err := xolo.NewServer()
	assert(err)

	assert(xserver.AddPool("main", time.Second*4))
	assert(xserver.AddRpcClient("main", rpclient))
	_, err = xserver.AddSigner("main", key, "11111111")
	assert(err)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.POST("/tx", xserver.ServePostTx)
	engine.GET("/tx/:txhash", xserver.ServeGetTx)

	go xserver.Start(engine)
	go engine.Run("0.0.0.0:8004")

	go func() {
		for {
			fmt.Println(xserver.Info())
			time.Sleep(10 * time.Second)
		}
	}()

	xbackend, err := xolo.NewAbiBackend(rpcURL, "http://localhost:8004", "main")
	assert(err)

	counter, err := abigen.NewCounter(common.HexToAddress("0x7cf3ab3954ac41a53294d55262b5bc5c62c2b000"), xbackend)
	assert(err)

	session := abigen.CounterSession{
		Contract:     counter,
		CallOpts:     bind.CallOpts{},
		TransactOpts: *xbackend.TransactOps(),
	}

	previousI, err := session.I()
	assert(err)

	tx, err := session.Inc()
	assert(err)
	receipt, err := xbackend.WaitReceipt(context.Background(), tx)
	assert(err)

	if receipt.Status != types.ReceiptStatusSuccessful {
		assert(errors.New("!ReceiptStatusSuccessful"))
	}

	nextI, err := session.I()
	assert(err)

	if nextI.Uint64()-previousI.Uint64() != 1 {
		panic("failed to increment")
	}

	fmt.Printf("SUCESS!")

}
