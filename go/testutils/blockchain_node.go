package testutils

import (
	"crypto/ecdsa"
	"os"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"

	log "github.com/sirupsen/logrus"
)

type ContractAddresses struct {
	Leagues        string
	Assets         string
	Evolution      string
	Engine         string
	Engineprecomp  string
	Updates        string
	Market         string
	Utils          string
	Playandevolve  string
	Shop           string
	Trainingpoints string
}

type BlockchainNode struct {
	Client    *ethclient.Client
	Owner     *ecdsa.PrivateKey
	Contracts *contracts.Contracts
	Addresses ContractAddresses
}

// AssertNoErr - log fatal and panic on error and print params
func AssertNoErr(err error, params ...interface{}) {
	if err != nil {
		log.Fatal(err, params)
	}
}

func NewBlockchainNodeDeployAndInitAt(address string) (*BlockchainNode, error) {
	node, err := NewBlockchainNodeAt(address)
	if err != nil {
		return nil, err
	}
	err = node.DeployContracts(node.Owner)
	if err != nil {
		return nil, err
	}
	err = node.InitOneTimezone(1)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func NewBlockchainNodeAt(address string) (*BlockchainNode, error) {
	client, err := ethclient.Dial(address)
	if err != nil {
		return nil, err
	}
	creatorPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	if err != nil {
		return nil, err
	}

	return &BlockchainNode{
		client,
		creatorPrivateKey,
		nil,
		ContractAddresses{},
	}, nil
}

func NewBlockchainNodeDeployAndInit() (*BlockchainNode, error) {
	return NewBlockchainNodeDeployAndInitAt("http://localhost:8545")
}

func NewBlockchainNode() (*BlockchainNode, error) {
	return NewBlockchainNodeAt("http://localhost:8545")
}

func deplyByTruffle() (map[string]string, error) {
	cryptoRoot, err := exec.Command("/usr/bin/git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Repo root at: %s", cryptoRoot)
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if err = os.Chdir(string(cryptoRoot[:len(cryptoRoot)-1]) + "/truffle-core"); err != nil {
		return nil, err
	}
	cmd := exec.Command("./node_modules/.bin/truffle", "migrate", "--network", "local", "--reset")
	log.Infof("Deploy by truffle: %v", cmd.String())
	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	// log.Infof("%s", o)
	output := string(o)
	startIdx := strings.Index(output, "-----------AddressesStart-----------") + len("-----------AddressesStart-----------")
	endIdx := strings.Index(output, "-----------AddressesEnd-----------")
	var contracts map[string]string
	contracts = make(map[string]string)
	addresses := strings.Split(output[startIdx+1:endIdx-1], "\n")
	for _, address := range addresses {
		log.Info(address)
		pair := strings.SplitN(address, "=", 2)
		contracts[pair[0]] = pair[1]
	}
	if err = os.Chdir(workingDir); err != nil {
		return nil, err
	}
	return contracts, nil
}

func (b *BlockchainNode) DeployContracts(owner *ecdsa.PrivateKey) error {
	contractMap, err := deplyByTruffle()
	if err != nil {
		return err
	}

	b.Contracts, _ = contracts.New(
		b.Client,
		contractMap["LEAGUES_CONTRACT_ADDRESS"],
		contractMap["ASSETS_CONTRACT_ADDRESS"],
		contractMap["EVOLUTION_CONTRACT_ADDRESS"],
		contractMap["ENGINE_CONTRACT_ADDRESS"],
		contractMap["ENGINEPRECOMP_CONTRACT_ADDRESS"],
		contractMap["UPDATES_CONTRACT_ADDRESS"],
		contractMap["MARKET_CONTRACT_ADDRESS"],
		contractMap["UTILS_CONTRACT_ADDRESS"],
		contractMap["PLAYANDEVOLVE_CONTRACT_ADDRESS"],
		contractMap["SHOP_CONTRACT_ADDRESS"],
		contractMap["TRAININGPOINTS_CONTRACT_ADDRESS"],
		contractMap["CONSTANTSGETTERS_CONTRACT_ADDRESS"],
	)

	b.Addresses = ContractAddresses{
		contractMap["LEAGUES_CONTRACT_ADDRESS"],
		contractMap["ASSETS_CONTRACT_ADDRESS"],
		contractMap["EVOLUTION_CONTRACT_ADDRESS"],
		contractMap["ENGINE_CONTRACT_ADDRESS"],
		contractMap["ENGINEPRECOMP_CONTRACT_ADDRESS"],
		contractMap["UPDATES_CONTRACT_ADDRESS"],
		contractMap["MARKET_CONTRACT_ADDRESS"],
		contractMap["UTILS_CONTRACT_ADDRESS"],
		contractMap["PLAYANDEVOLVE_CONTRACT_ADDRESS"],
		contractMap["SHOP_CONTRACT_ADDRESS"],
		contractMap["TRAININGPOINTS_CONTRACT_ADDRESS"],
	}
	return nil
}

func (b *BlockchainNode) Init() error {
	// Initing
	tx, err := b.Contracts.Assets.Init(bind.NewKeyedTransactor(b.Owner))
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx, 10)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlockchainNode) InitOneTimezone(timezoneIdx uint8) error {
	// Initing
	tx, err := b.Contracts.Assets.InitSingleTZ(bind.NewKeyedTransactor(b.Owner), timezoneIdx)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx, 10)
	if err != nil {
		return err
	}
	return nil
}
