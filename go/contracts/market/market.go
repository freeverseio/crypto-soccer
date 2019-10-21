// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package market

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MarketABI is the input ABI used to generate the binding from.
const MarketABI = "[{\"inputs\":[],\"constant\":true,\"name\":\"POST_AUCTION_TIME\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"PLAYERS_PER_TEAM_MAX\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"FREE_PLAYER_ID\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"MAX_VALID_UNTIL\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"AUCTION_TIME\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"playerId\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"auctionData\"},{\"indexed\":false,\"type\":\"bool\",\"name\":\"frozen\"}],\"type\":\"event\",\"name\":\"PlayerFreeze\",\"anonymous\":false},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"teamId\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"auctionData\"},{\"indexed\":false,\"type\":\"bool\",\"name\":\"frozen\"}],\"type\":\"event\",\"name\":\"TeamFreeze\",\"anonymous\":false},{\"inputs\":[{\"type\":\"address\",\"name\":\"addr\"}],\"constant\":false,\"name\":\"setAssetsAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"bytes32[3]\",\"name\":\"sig\"},{\"type\":\"uint8\",\"name\":\"sigV\"}],\"constant\":false,\"name\":\"freezePlayer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"bytes32\",\"name\":\"buyerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"buyerTeamId\"},{\"type\":\"bytes32[3]\",\"name\":\"sig\"},{\"type\":\"uint8\",\"name\":\"sigV\"},{\"type\":\"bool\",\"name\":\"isOffer2StartAuction\"}],\"constant\":false,\"name\":\"completePlayerAuction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"bytes32[3]\",\"name\":\"sig\"},{\"type\":\"uint8\",\"name\":\"sigV\"}],\"constant\":false,\"name\":\"freezeTeam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"bytes32\",\"name\":\"buyerHiddenPrice\"},{\"type\":\"bytes32[3]\",\"name\":\"sig\"},{\"type\":\"uint8\",\"name\":\"sigV\"},{\"type\":\"bool\",\"name\":\"isOffer2StartAuction\"}],\"constant\":false,\"name\":\"completeTeamAuction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"bytes32[3]\",\"name\":\"sig\"},{\"type\":\"uint8\",\"name\":\"sigV\"}],\"constant\":true,\"name\":\"areFreezeTeamRequirementsOK\",\"outputs\":[{\"type\":\"bool\",\"name\":\"ok\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"bytes32\",\"name\":\"buyerHiddenPrice\"},{\"type\":\"bytes32[3]\",\"name\":\"sig\"},{\"type\":\"uint8\",\"name\":\"sigV\"},{\"type\":\"bool\",\"name\":\"isOffer2StartAuction\"}],\"constant\":true,\"name\":\"areCompleteTeamAuctionRequirementsOK\",\"outputs\":[{\"type\":\"bool\",\"name\":\"ok\"},{\"type\":\"address\",\"name\":\"buyerAddress\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"bytes32[3]\",\"name\":\"sig\"},{\"type\":\"uint8\",\"name\":\"sigV\"}],\"constant\":true,\"name\":\"areFreezePlayerRequirementsOK\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"bytes32\",\"name\":\"buyerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"buyerTeamId\"},{\"type\":\"bytes32[3]\",\"name\":\"sig\"},{\"type\":\"uint8\",\"name\":\"sigV\"},{\"type\":\"bool\",\"name\":\"isOffer2StartAuction\"}],\"constant\":true,\"name\":\"areCompletePlayerAuctionRequirementsOK\",\"outputs\":[{\"type\":\"bool\",\"name\":\"ok\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"currencyId\"},{\"type\":\"uint256\",\"name\":\"price\"},{\"type\":\"uint256\",\"name\":\"rnd\"}],\"constant\":true,\"name\":\"hashPrivateMsg\",\"outputs\":[{\"type\":\"bytes32\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"hiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"assetId\"}],\"constant\":true,\"name\":\"buildPutAssetForSaleTxMsg\",\"outputs\":[{\"type\":\"bytes32\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"hiddenPrice\"},{\"type\":\"uint256\",\"name\":\"validUntil\"},{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint256\",\"name\":\"buyerTeamId\"}],\"constant\":true,\"name\":\"buildOfferToBuyTxMsg\",\"outputs\":[{\"type\":\"bytes32\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerTxHash\"},{\"type\":\"bytes32\",\"name\":\"buyerHiddenPrice\"},{\"type\":\"uint256\",\"name\":\"buyerTeamId\"},{\"type\":\"bool\",\"name\":\"isOffer2StartAuction\"}],\"constant\":true,\"name\":\"buildAgreeToBuyPlayerTxMsg\",\"outputs\":[{\"type\":\"bytes32\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"sellerTxHash\"},{\"type\":\"bytes32\",\"name\":\"buyerHiddenPrice\"},{\"type\":\"bool\",\"name\":\"isOffer2StartAuction\"}],\"constant\":true,\"name\":\"buildAgreeToBuyTeamTxMsg\",\"outputs\":[{\"type\":\"bytes32\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"address\",\"name\":\"_addr\"},{\"type\":\"bytes32\",\"name\":\"msgHash\"},{\"type\":\"uint8\",\"name\":\"v\"},{\"type\":\"bytes32\",\"name\":\"r\"},{\"type\":\"bytes32\",\"name\":\"s\"}],\"constant\":true,\"name\":\"isSigned\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"hash\"}],\"constant\":true,\"name\":\"prefixed\",\"outputs\":[{\"type\":\"bytes32\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"getBlockchainNowTime\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"isPlayerFrozen\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"isTeamFrozen\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"}]"

// MarketBin is the compiled bytecode used for deploying new contracts.
const MarketBin = `0x608060405234801561001057600080fd5b50612142806100206000396000f3fe608060405234801561001057600080fd5b50600436106101585760003560e01c80638677ebe8116100c3578063b896857d1161007c578063b896857d146108a1578063c258012b1461096a578063dc29805714610988578063e592301a146109a6578063e7ff8b0e146109c4578063eb0f1ea014610a0a57610158565b80638677ebe8146105e55780638adddc9d1461066c5780638ce954a8146106905780638cfd3125146107385780639cdbbe7b14610798578063a51dcaa3146107fa57610158565b806360f4818e1161011557806360f4818e14610367578063628be8241461041857806363ef7460146104c05780636612cdff14610519578063663c2e4b1461053757806375a554711461055557610158565b80630db2b6ae1461015d5780631490a174146101a35780631f275713146101e75780633cfaf14c146102295780634a74f553146102815780635e8f77e7146102d7575b600080fd5b6101896004803603602081101561017357600080fd5b8101908080359060200190929190505050610afc565b604051808215151515815260200191505060405180910390f35b6101e5600480360360208110156101b957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610c45565b005b610213600480360360208110156101fd57600080fd5b8101908080359060200190929190505050610c88565b6040518082815260200191505060405180910390f35b61026b6004803603606081101561023f57600080fd5b810190808035906020019092919080359060200190929190803515159060200190929190505050610ce0565b6040518082815260200191505060405180910390f35b6102c16004803603606081101561029757600080fd5b81019080803590602001909291908035906020019092919080359060200190929190505050610d26565b6040518082815260200191505060405180910390f35b610365600480360360e08110156102ed57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff169060200190929190505050610d68565b005b610416600480360361014081101561037e57600080fd5b81019080803590602001909291908035906020019092919080359060200190929190803590602001909291908035906020019092919080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff169060200190929190803515159060200190929190505050610e6e565b005b6104a6600480360360e081101561042e57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff169060200190929190505050610fd6565b604051808215151515815260200191505060405180910390f35b610503600480360360608110156104d657600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190505050611306565b6040518082815260200191505060405180910390f35b61052161134e565b6040518082815260200191505060405180910390f35b61053f611354565b6040518082815260200191505060405180910390f35b6105e3600480360360e081101561056b57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff16906020019092919050505061135c565b005b610652600480360360a08110156105fb57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803560ff1690602001909291908035906020019092919080359060200190929190505050611462565b604051808215151515815260200191505060405180910390f35b610674611504565b604051808260ff1660ff16815260200191505060405180910390f35b61071e600480360360e08110156106a657600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff169060200190929190505050611509565b604051808215151515815260200191505060405180910390f35b6107826004803603608081101561074e57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080359060200190929190505050611814565b6040518082815260200191505060405180910390f35b6107e4600480360360808110156107ae57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919080351515906020019092919050505061185f565b6040518082815260200191505060405180910390f35b61089f600480360361012081101561081157600080fd5b810190808035906020019092919080359060200190929190803590602001909291908035906020019092919080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff1690602001909291908035151590602001909291905050506118ae565b005b61095060048036036101408110156108b857600080fd5b81019080803590602001909291908035906020019092919080359060200190929190803590602001909291908035906020019092919080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff169060200190929190803515159060200190929190505050611a4b565b604051808215151515815260200191505060405180910390f35b610972611d67565b6040518082815260200191505060405180910390f35b610990611d6c565b6040518082815260200191505060405180910390f35b6109ae611d73565b6040518082815260200191505060405180910390f35b6109f0600480360360208110156109da57600080fd5b8101908080359060200190929190505050611d7a565b604051808215151515815260200191505060405180910390f35b610aaf6004803603610120811015610a2157600080fd5b810190808035906020019092919080359060200190929190803590602001909291908035906020019092919080606001906003806020026040519081016040528092919082600360200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff169060200190929190803515159060200190929190505050611ec3565b60405180831515151581526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663bc1a97c1836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015610b7057600080fd5b505afa158015610b84573d6000803e3d6000fd5b505050506040513d6020811015610b9a57600080fd5b8101908080519060200190929190505050610c1d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f756e6578697374656e7420706c6179657200000000000000000000000000000081525060200191505060405180910390fd5b426154606403ffffffff60016000868152602001908152602001600020541601119050919050565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008160405160200180807f19457468657265756d205369676e6564204d6573736167653a0a333200000000815250601c01828152602001915050604051602081830303815290604052805190602001209050919050565b6000838383604051602001808481526020018381526020018215151515815260200193505050506040516020818303038152906040528051906020012090509392505050565b60008383836040516020018084815260200183815260200182815260200193505050506040516020818303038152906040528051906020012090509392505050565b610d758585858585610fd6565b610de7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f46726565506c6179657220726571756972656d656e7473206e6f74206d65740081525060200191505060405180910390fd5b602285901b60001c840160016000858152602001908152602001600020819055507f44b30f34d8f6f8c3cc8737fe3476b3bfd6fff21e03ef669b48b787a213b1f11083600160008681526020019081526020016000205460016040518084815260200183815260200182151515158152602001935050505060405180910390a15050505050565b610e7e8888888888888888611a4b565b610ed3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c8152602001806120eb602c913960400191505060405180910390fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663800257d587866040518363ffffffff1660e01b81526004018083815260200182815260200192505050600060405180830381600087803b158015610f4f57600080fd5b505af1158015610f63573d6000803e3d6000fd5b505050506001806000888152602001908152602001600020819055507f44b30f34d8f6f8c3cc8737fe3476b3bfd6fff21e03ef669b48b787a213b1f11086600160006040518084815260200183815260200182151515158152602001935050505060405180910390a15050505050505050565b60008442108015610fed5750610feb84610afc565b155b80156110ab57506110a96000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166338c96b5c866040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561106957600080fd5b505afa15801561107d573d6000803e3d6000fd5b505050506040513d602081101561109357600080fd5b8101908080519060200190929190505050611d7a565b155b80156111905750600073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638f9da214866040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561113c57600080fd5b505afa158015611150573d6000803e3d6000fd5b505050506040513d602081101561116657600080fd5b810190808051906020019092919050505073ffffffffffffffffffffffffffffffffffffffff1614155b80156112ba57506111df83600060ff16600381106111aa57fe5b60200201518385600160ff16600381106111c057fe5b602002015186600260ff16600381106111d557fe5b6020020151612054565b73ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638f9da214866040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561126757600080fd5b505afa15801561127b573d6000803e3d6000fd5b505050506040513d602081101561129157600080fd5b810190808051906020019092919050505073ffffffffffffffffffffffffffffffffffffffff16145b80156112fb57506112d46112cf878787610d26565b610c88565b83600060ff16600381106112e457fe5b60200201511480156112fa57506201a5e0420185105b5b905095945050505050565b6000838383604051602001808460ff1660ff16815260200183815260200182815260200193505050506040516020818303038152906040528051906020012090509392505050565b61546081565b600042905090565b6113698585858585611509565b6113db576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f46726565506c6179657220726571756972656d656e7473206e6f74206d65740081525060200191505060405180910390fd5b602285901b60001c840160026000858152602001908152602001600020819055507fbe5cb47c7008dd92757f872cf47b2f27eae3e4d3efbba0bd54969ed71c927d0e83600260008681526020019081526020016000205460016040518084815260200183815260200182151515158152602001935050505060405180910390a15050505050565b60008573ffffffffffffffffffffffffffffffffffffffff1660018686868660405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156114d8573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff1614905095945050505050565b601981565b6000806000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663492afc69866040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561157e57600080fd5b505afa158015611592573d6000803e3d6000fd5b505050506040513d60208110156115a857600080fd5b8101908080519060200190929190505050905085421080156115d057506115ce85611d7a565b155b80156116095750600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614155b8015611687575061165884600060ff166003811061162357fe5b60200201518486600160ff166003811061163957fe5b602002015187600260ff166003811061164e57fe5b6020020151612054565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16145b80156116b857506116a161169c888888610d26565b610c88565b84600060ff16600381106116b157fe5b6020020151145b80156116c857506201a5e0420186105b9150816116d957600091505061180b565b6116e16120c7565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eabf6a4b876040518263ffffffff1660e01b8152600401808281526020019150506103206040518083038186803b15801561175457600080fd5b505afa158015611768573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525061032081101561178e57600080fd5b8101908091905050905060008090505b601960ff168160ff161015611807576001828260ff16601981106117be57fe5b6020020151141580156117e957506117e8828260ff16601981106117de57fe5b6020020151610afc565b5b156117fa576000935050505061180b565b808060010191505061179e565b5050505b95945050505050565b60008484848460405160200180858152602001848152602001838152602001828152602001945050505050604051602081830303815290604052805190602001209050949350505050565b6000848484846040516020018085815260200184815260200183815260200182151515158152602001945050505050604051602081830303815290604052805190602001209050949350505050565b6000806118c089898989898989611ec3565b915091508161191a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c8152602001806120eb602c913960400191505060405180910390fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e945e96a88836040518363ffffffff1660e01b8152600401808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200192505050600060405180830381600087803b1580156119c257600080fd5b505af11580156119d6573d6000803e3d6000fd5b50505050600160026000898152602001908152602001600020819055507fbe5cb47c7008dd92757f872cf47b2f27eae3e4d3efbba0bd54969ed71c927d0e87600160006040518084815260200183815260200182151515158152602001935050505060405180910390a1505050505050505050565b600080611a61611a5c8b8b8b610d26565b610c88565b9050600073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663492afc69886040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015611aed57600080fd5b505afa158015611b01573d6000803e3d6000fd5b505050506040513d6020811015611b1757600080fd5b810190808051906020019092919050505073ffffffffffffffffffffffffffffffffffffffff1614158015611b8a57506022600160008a815260200190815260200160002054901c7b400000000000000000000000000000000000000000000000000000008b60001c81611b8757fe5b06145b8015611cb45750611bd985600060ff1660038110611ba457fe5b60200201518587600160ff1660038110611bba57fe5b602002015188600260ff1660038110611bcf57fe5b6020020151612054565b73ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663492afc69886040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015611c6157600080fd5b505afa158015611c75573d6000803e3d6000fd5b505050506040513d6020811015611c8b57600080fd5b810190808051906020019092919050505073ffffffffffffffffffffffffffffffffffffffff16145b8015611cc55750611cc488610afc565b5b8015611cf75750611ce0611cdb8289898761185f565b610c88565b85600060ff1660038110611cf057fe5b6020020151145b91508215611d3157818015611d2a5750620151806403ffffffff600160008b815260200190815260200160002054160389115b9150611d5a565b818015611d5757506403ffffffff600160008a8152602001908152602001600020541689145b91505b5098975050505050505050565b600181565b6201a5e081565b6201518081565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166398981756836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015611dee57600080fd5b505afa158015611e02573d6000803e3d6000fd5b505050506040513d6020811015611e1857600080fd5b8101908080519060200190929190505050611e9b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f756e6578697374656e74207465616d000000000000000000000000000000000081525060200191505060405180910390fd5b426154606403ffffffff60026000868152602001908152602001600020541601119050919050565b6000806000611edb611ed68b8b8b610d26565b610c88565b9050611f2586600060ff1660038110611ef057fe5b60200201518688600160ff1660038110611f0657fe5b602002015189600260ff1660038110611f1b57fe5b6020020151612054565b9150600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614158015611fa257506022600260008a815260200190815260200160002054901c7b400000000000000000000000000000000000000000000000000000008b60001c81611f9f57fe5b06145b8015611fb35750611fb288611d7a565b5b8015611fe45750611fcd611fc8828987610ce0565b610c88565b86600060ff1660038110611fdd57fe5b6020020151145b9250831561201e578280156120175750620151806403ffffffff600260008b815260200190815260200160002054160389115b9250612047565b82801561204457506403ffffffff600260008a8152602001908152602001600020541689145b92505b5097509795505050505050565b600060018585858560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156120b3573d6000803e3d6000fd5b505050602060405103519050949350505050565b60405180610320016040528060199060208202803883398082019150509050509056fe726571756972656d656e747320746f20636f6d706c6574652061756374696f6e20617265206e6f74206d6574a165627a7a723058207e4d8a27c57df9b59adaed6160a9675dbd72d21835ac385fbb5df10acd3bb2290029`

// DeployMarket deploys a new Ethereum contract, binding an instance of Market to it.
func DeployMarket(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Market, error) {
	parsed, err := abi.JSON(strings.NewReader(MarketABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MarketBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Market{MarketCaller: MarketCaller{contract: contract}, MarketTransactor: MarketTransactor{contract: contract}, MarketFilterer: MarketFilterer{contract: contract}}, nil
}

// Market is an auto generated Go binding around an Ethereum contract.
type Market struct {
	MarketCaller     // Read-only binding to the contract
	MarketTransactor // Write-only binding to the contract
	MarketFilterer   // Log filterer for contract events
}

// MarketCaller is an auto generated read-only Go binding around an Ethereum contract.
type MarketCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MarketTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MarketTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MarketFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MarketFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MarketSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MarketSession struct {
	Contract     *Market           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MarketCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MarketCallerSession struct {
	Contract *MarketCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MarketTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MarketTransactorSession struct {
	Contract     *MarketTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MarketRaw is an auto generated low-level Go binding around an Ethereum contract.
type MarketRaw struct {
	Contract *Market // Generic contract binding to access the raw methods on
}

// MarketCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MarketCallerRaw struct {
	Contract *MarketCaller // Generic read-only contract binding to access the raw methods on
}

// MarketTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MarketTransactorRaw struct {
	Contract *MarketTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMarket creates a new instance of Market, bound to a specific deployed contract.
func NewMarket(address common.Address, backend bind.ContractBackend) (*Market, error) {
	contract, err := bindMarket(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Market{MarketCaller: MarketCaller{contract: contract}, MarketTransactor: MarketTransactor{contract: contract}, MarketFilterer: MarketFilterer{contract: contract}}, nil
}

// NewMarketCaller creates a new read-only instance of Market, bound to a specific deployed contract.
func NewMarketCaller(address common.Address, caller bind.ContractCaller) (*MarketCaller, error) {
	contract, err := bindMarket(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MarketCaller{contract: contract}, nil
}

// NewMarketTransactor creates a new write-only instance of Market, bound to a specific deployed contract.
func NewMarketTransactor(address common.Address, transactor bind.ContractTransactor) (*MarketTransactor, error) {
	contract, err := bindMarket(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MarketTransactor{contract: contract}, nil
}

// NewMarketFilterer creates a new log filterer instance of Market, bound to a specific deployed contract.
func NewMarketFilterer(address common.Address, filterer bind.ContractFilterer) (*MarketFilterer, error) {
	contract, err := bindMarket(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MarketFilterer{contract: contract}, nil
}

// bindMarket binds a generic wrapper to an already deployed contract.
func bindMarket(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MarketABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Market *MarketRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Market.Contract.MarketCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Market *MarketRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Market.Contract.MarketTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Market *MarketRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Market.Contract.MarketTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Market *MarketCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Market.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Market *MarketTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Market.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Market *MarketTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Market.Contract.contract.Transact(opts, method, params...)
}

// AUCTIONTIME is a free data retrieval call binding the contract method 0xe592301a.
//
// Solidity: function AUCTION_TIME() constant returns(uint256)
func (_Market *MarketCaller) AUCTIONTIME(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "AUCTION_TIME")
	return *ret0, err
}

// AUCTIONTIME is a free data retrieval call binding the contract method 0xe592301a.
//
// Solidity: function AUCTION_TIME() constant returns(uint256)
func (_Market *MarketSession) AUCTIONTIME() (*big.Int, error) {
	return _Market.Contract.AUCTIONTIME(&_Market.CallOpts)
}

// AUCTIONTIME is a free data retrieval call binding the contract method 0xe592301a.
//
// Solidity: function AUCTION_TIME() constant returns(uint256)
func (_Market *MarketCallerSession) AUCTIONTIME() (*big.Int, error) {
	return _Market.Contract.AUCTIONTIME(&_Market.CallOpts)
}

// FREEPLAYERID is a free data retrieval call binding the contract method 0xc258012b.
//
// Solidity: function FREE_PLAYER_ID() constant returns(uint256)
func (_Market *MarketCaller) FREEPLAYERID(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "FREE_PLAYER_ID")
	return *ret0, err
}

// FREEPLAYERID is a free data retrieval call binding the contract method 0xc258012b.
//
// Solidity: function FREE_PLAYER_ID() constant returns(uint256)
func (_Market *MarketSession) FREEPLAYERID() (*big.Int, error) {
	return _Market.Contract.FREEPLAYERID(&_Market.CallOpts)
}

// FREEPLAYERID is a free data retrieval call binding the contract method 0xc258012b.
//
// Solidity: function FREE_PLAYER_ID() constant returns(uint256)
func (_Market *MarketCallerSession) FREEPLAYERID() (*big.Int, error) {
	return _Market.Contract.FREEPLAYERID(&_Market.CallOpts)
}

// MAXVALIDUNTIL is a free data retrieval call binding the contract method 0xdc298057.
//
// Solidity: function MAX_VALID_UNTIL() constant returns(uint256)
func (_Market *MarketCaller) MAXVALIDUNTIL(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "MAX_VALID_UNTIL")
	return *ret0, err
}

// MAXVALIDUNTIL is a free data retrieval call binding the contract method 0xdc298057.
//
// Solidity: function MAX_VALID_UNTIL() constant returns(uint256)
func (_Market *MarketSession) MAXVALIDUNTIL() (*big.Int, error) {
	return _Market.Contract.MAXVALIDUNTIL(&_Market.CallOpts)
}

// MAXVALIDUNTIL is a free data retrieval call binding the contract method 0xdc298057.
//
// Solidity: function MAX_VALID_UNTIL() constant returns(uint256)
func (_Market *MarketCallerSession) MAXVALIDUNTIL() (*big.Int, error) {
	return _Market.Contract.MAXVALIDUNTIL(&_Market.CallOpts)
}

// PLAYERSPERTEAMMAX is a free data retrieval call binding the contract method 0x8adddc9d.
//
// Solidity: function PLAYERS_PER_TEAM_MAX() constant returns(uint8)
func (_Market *MarketCaller) PLAYERSPERTEAMMAX(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "PLAYERS_PER_TEAM_MAX")
	return *ret0, err
}

// PLAYERSPERTEAMMAX is a free data retrieval call binding the contract method 0x8adddc9d.
//
// Solidity: function PLAYERS_PER_TEAM_MAX() constant returns(uint8)
func (_Market *MarketSession) PLAYERSPERTEAMMAX() (uint8, error) {
	return _Market.Contract.PLAYERSPERTEAMMAX(&_Market.CallOpts)
}

// PLAYERSPERTEAMMAX is a free data retrieval call binding the contract method 0x8adddc9d.
//
// Solidity: function PLAYERS_PER_TEAM_MAX() constant returns(uint8)
func (_Market *MarketCallerSession) PLAYERSPERTEAMMAX() (uint8, error) {
	return _Market.Contract.PLAYERSPERTEAMMAX(&_Market.CallOpts)
}

// POSTAUCTIONTIME is a free data retrieval call binding the contract method 0x6612cdff.
//
// Solidity: function POST_AUCTION_TIME() constant returns(uint256)
func (_Market *MarketCaller) POSTAUCTIONTIME(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "POST_AUCTION_TIME")
	return *ret0, err
}

// POSTAUCTIONTIME is a free data retrieval call binding the contract method 0x6612cdff.
//
// Solidity: function POST_AUCTION_TIME() constant returns(uint256)
func (_Market *MarketSession) POSTAUCTIONTIME() (*big.Int, error) {
	return _Market.Contract.POSTAUCTIONTIME(&_Market.CallOpts)
}

// POSTAUCTIONTIME is a free data retrieval call binding the contract method 0x6612cdff.
//
// Solidity: function POST_AUCTION_TIME() constant returns(uint256)
func (_Market *MarketCallerSession) POSTAUCTIONTIME() (*big.Int, error) {
	return _Market.Contract.POSTAUCTIONTIME(&_Market.CallOpts)
}

// AreCompletePlayerAuctionRequirementsOK is a free data retrieval call binding the contract method 0xb896857d.
//
// Solidity: function areCompletePlayerAuctionRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) constant returns(bool ok)
func (_Market *MarketCaller) AreCompletePlayerAuctionRequirementsOK(opts *bind.CallOpts, sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, buyerHiddenPrice [32]byte, buyerTeamId *big.Int, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "areCompletePlayerAuctionRequirementsOK", sellerHiddenPrice, validUntil, playerId, buyerHiddenPrice, buyerTeamId, sig, sigV, isOffer2StartAuction)
	return *ret0, err
}

// AreCompletePlayerAuctionRequirementsOK is a free data retrieval call binding the contract method 0xb896857d.
//
// Solidity: function areCompletePlayerAuctionRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) constant returns(bool ok)
func (_Market *MarketSession) AreCompletePlayerAuctionRequirementsOK(sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, buyerHiddenPrice [32]byte, buyerTeamId *big.Int, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (bool, error) {
	return _Market.Contract.AreCompletePlayerAuctionRequirementsOK(&_Market.CallOpts, sellerHiddenPrice, validUntil, playerId, buyerHiddenPrice, buyerTeamId, sig, sigV, isOffer2StartAuction)
}

// AreCompletePlayerAuctionRequirementsOK is a free data retrieval call binding the contract method 0xb896857d.
//
// Solidity: function areCompletePlayerAuctionRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) constant returns(bool ok)
func (_Market *MarketCallerSession) AreCompletePlayerAuctionRequirementsOK(sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, buyerHiddenPrice [32]byte, buyerTeamId *big.Int, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (bool, error) {
	return _Market.Contract.AreCompletePlayerAuctionRequirementsOK(&_Market.CallOpts, sellerHiddenPrice, validUntil, playerId, buyerHiddenPrice, buyerTeamId, sig, sigV, isOffer2StartAuction)
}

// AreCompleteTeamAuctionRequirementsOK is a free data retrieval call binding the contract method 0xeb0f1ea0.
//
// Solidity: function areCompleteTeamAuctionRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32 buyerHiddenPrice, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) constant returns(bool ok, address buyerAddress)
func (_Market *MarketCaller) AreCompleteTeamAuctionRequirementsOK(opts *bind.CallOpts, sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, buyerHiddenPrice [32]byte, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (struct {
	Ok           bool
	BuyerAddress common.Address
}, error) {
	ret := new(struct {
		Ok           bool
		BuyerAddress common.Address
	})
	out := ret
	err := _Market.contract.Call(opts, out, "areCompleteTeamAuctionRequirementsOK", sellerHiddenPrice, validUntil, teamId, buyerHiddenPrice, sig, sigV, isOffer2StartAuction)
	return *ret, err
}

// AreCompleteTeamAuctionRequirementsOK is a free data retrieval call binding the contract method 0xeb0f1ea0.
//
// Solidity: function areCompleteTeamAuctionRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32 buyerHiddenPrice, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) constant returns(bool ok, address buyerAddress)
func (_Market *MarketSession) AreCompleteTeamAuctionRequirementsOK(sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, buyerHiddenPrice [32]byte, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (struct {
	Ok           bool
	BuyerAddress common.Address
}, error) {
	return _Market.Contract.AreCompleteTeamAuctionRequirementsOK(&_Market.CallOpts, sellerHiddenPrice, validUntil, teamId, buyerHiddenPrice, sig, sigV, isOffer2StartAuction)
}

// AreCompleteTeamAuctionRequirementsOK is a free data retrieval call binding the contract method 0xeb0f1ea0.
//
// Solidity: function areCompleteTeamAuctionRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32 buyerHiddenPrice, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) constant returns(bool ok, address buyerAddress)
func (_Market *MarketCallerSession) AreCompleteTeamAuctionRequirementsOK(sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, buyerHiddenPrice [32]byte, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (struct {
	Ok           bool
	BuyerAddress common.Address
}, error) {
	return _Market.Contract.AreCompleteTeamAuctionRequirementsOK(&_Market.CallOpts, sellerHiddenPrice, validUntil, teamId, buyerHiddenPrice, sig, sigV, isOffer2StartAuction)
}

// AreFreezePlayerRequirementsOK is a free data retrieval call binding the contract method 0x628be824.
//
// Solidity: function areFreezePlayerRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32[3] sig, uint8 sigV) constant returns(bool)
func (_Market *MarketCaller) AreFreezePlayerRequirementsOK(opts *bind.CallOpts, sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, sig [3][32]byte, sigV uint8) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "areFreezePlayerRequirementsOK", sellerHiddenPrice, validUntil, playerId, sig, sigV)
	return *ret0, err
}

// AreFreezePlayerRequirementsOK is a free data retrieval call binding the contract method 0x628be824.
//
// Solidity: function areFreezePlayerRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32[3] sig, uint8 sigV) constant returns(bool)
func (_Market *MarketSession) AreFreezePlayerRequirementsOK(sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, sig [3][32]byte, sigV uint8) (bool, error) {
	return _Market.Contract.AreFreezePlayerRequirementsOK(&_Market.CallOpts, sellerHiddenPrice, validUntil, playerId, sig, sigV)
}

// AreFreezePlayerRequirementsOK is a free data retrieval call binding the contract method 0x628be824.
//
// Solidity: function areFreezePlayerRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32[3] sig, uint8 sigV) constant returns(bool)
func (_Market *MarketCallerSession) AreFreezePlayerRequirementsOK(sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, sig [3][32]byte, sigV uint8) (bool, error) {
	return _Market.Contract.AreFreezePlayerRequirementsOK(&_Market.CallOpts, sellerHiddenPrice, validUntil, playerId, sig, sigV)
}

// AreFreezeTeamRequirementsOK is a free data retrieval call binding the contract method 0x8ce954a8.
//
// Solidity: function areFreezeTeamRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32[3] sig, uint8 sigV) constant returns(bool ok)
func (_Market *MarketCaller) AreFreezeTeamRequirementsOK(opts *bind.CallOpts, sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, sig [3][32]byte, sigV uint8) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "areFreezeTeamRequirementsOK", sellerHiddenPrice, validUntil, teamId, sig, sigV)
	return *ret0, err
}

// AreFreezeTeamRequirementsOK is a free data retrieval call binding the contract method 0x8ce954a8.
//
// Solidity: function areFreezeTeamRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32[3] sig, uint8 sigV) constant returns(bool ok)
func (_Market *MarketSession) AreFreezeTeamRequirementsOK(sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, sig [3][32]byte, sigV uint8) (bool, error) {
	return _Market.Contract.AreFreezeTeamRequirementsOK(&_Market.CallOpts, sellerHiddenPrice, validUntil, teamId, sig, sigV)
}

// AreFreezeTeamRequirementsOK is a free data retrieval call binding the contract method 0x8ce954a8.
//
// Solidity: function areFreezeTeamRequirementsOK(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32[3] sig, uint8 sigV) constant returns(bool ok)
func (_Market *MarketCallerSession) AreFreezeTeamRequirementsOK(sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, sig [3][32]byte, sigV uint8) (bool, error) {
	return _Market.Contract.AreFreezeTeamRequirementsOK(&_Market.CallOpts, sellerHiddenPrice, validUntil, teamId, sig, sigV)
}

// BuildAgreeToBuyPlayerTxMsg is a free data retrieval call binding the contract method 0x9cdbbe7b.
//
// Solidity: function buildAgreeToBuyPlayerTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bool isOffer2StartAuction) constant returns(bytes32)
func (_Market *MarketCaller) BuildAgreeToBuyPlayerTxMsg(opts *bind.CallOpts, sellerTxHash [32]byte, buyerHiddenPrice [32]byte, buyerTeamId *big.Int, isOffer2StartAuction bool) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "buildAgreeToBuyPlayerTxMsg", sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction)
	return *ret0, err
}

// BuildAgreeToBuyPlayerTxMsg is a free data retrieval call binding the contract method 0x9cdbbe7b.
//
// Solidity: function buildAgreeToBuyPlayerTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bool isOffer2StartAuction) constant returns(bytes32)
func (_Market *MarketSession) BuildAgreeToBuyPlayerTxMsg(sellerTxHash [32]byte, buyerHiddenPrice [32]byte, buyerTeamId *big.Int, isOffer2StartAuction bool) ([32]byte, error) {
	return _Market.Contract.BuildAgreeToBuyPlayerTxMsg(&_Market.CallOpts, sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction)
}

// BuildAgreeToBuyPlayerTxMsg is a free data retrieval call binding the contract method 0x9cdbbe7b.
//
// Solidity: function buildAgreeToBuyPlayerTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bool isOffer2StartAuction) constant returns(bytes32)
func (_Market *MarketCallerSession) BuildAgreeToBuyPlayerTxMsg(sellerTxHash [32]byte, buyerHiddenPrice [32]byte, buyerTeamId *big.Int, isOffer2StartAuction bool) ([32]byte, error) {
	return _Market.Contract.BuildAgreeToBuyPlayerTxMsg(&_Market.CallOpts, sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction)
}

// BuildAgreeToBuyTeamTxMsg is a free data retrieval call binding the contract method 0x3cfaf14c.
//
// Solidity: function buildAgreeToBuyTeamTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, bool isOffer2StartAuction) constant returns(bytes32)
func (_Market *MarketCaller) BuildAgreeToBuyTeamTxMsg(opts *bind.CallOpts, sellerTxHash [32]byte, buyerHiddenPrice [32]byte, isOffer2StartAuction bool) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "buildAgreeToBuyTeamTxMsg", sellerTxHash, buyerHiddenPrice, isOffer2StartAuction)
	return *ret0, err
}

// BuildAgreeToBuyTeamTxMsg is a free data retrieval call binding the contract method 0x3cfaf14c.
//
// Solidity: function buildAgreeToBuyTeamTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, bool isOffer2StartAuction) constant returns(bytes32)
func (_Market *MarketSession) BuildAgreeToBuyTeamTxMsg(sellerTxHash [32]byte, buyerHiddenPrice [32]byte, isOffer2StartAuction bool) ([32]byte, error) {
	return _Market.Contract.BuildAgreeToBuyTeamTxMsg(&_Market.CallOpts, sellerTxHash, buyerHiddenPrice, isOffer2StartAuction)
}

// BuildAgreeToBuyTeamTxMsg is a free data retrieval call binding the contract method 0x3cfaf14c.
//
// Solidity: function buildAgreeToBuyTeamTxMsg(bytes32 sellerTxHash, bytes32 buyerHiddenPrice, bool isOffer2StartAuction) constant returns(bytes32)
func (_Market *MarketCallerSession) BuildAgreeToBuyTeamTxMsg(sellerTxHash [32]byte, buyerHiddenPrice [32]byte, isOffer2StartAuction bool) ([32]byte, error) {
	return _Market.Contract.BuildAgreeToBuyTeamTxMsg(&_Market.CallOpts, sellerTxHash, buyerHiddenPrice, isOffer2StartAuction)
}

// BuildOfferToBuyTxMsg is a free data retrieval call binding the contract method 0x8cfd3125.
//
// Solidity: function buildOfferToBuyTxMsg(bytes32 hiddenPrice, uint256 validUntil, uint256 playerId, uint256 buyerTeamId) constant returns(bytes32)
func (_Market *MarketCaller) BuildOfferToBuyTxMsg(opts *bind.CallOpts, hiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, buyerTeamId *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "buildOfferToBuyTxMsg", hiddenPrice, validUntil, playerId, buyerTeamId)
	return *ret0, err
}

// BuildOfferToBuyTxMsg is a free data retrieval call binding the contract method 0x8cfd3125.
//
// Solidity: function buildOfferToBuyTxMsg(bytes32 hiddenPrice, uint256 validUntil, uint256 playerId, uint256 buyerTeamId) constant returns(bytes32)
func (_Market *MarketSession) BuildOfferToBuyTxMsg(hiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, buyerTeamId *big.Int) ([32]byte, error) {
	return _Market.Contract.BuildOfferToBuyTxMsg(&_Market.CallOpts, hiddenPrice, validUntil, playerId, buyerTeamId)
}

// BuildOfferToBuyTxMsg is a free data retrieval call binding the contract method 0x8cfd3125.
//
// Solidity: function buildOfferToBuyTxMsg(bytes32 hiddenPrice, uint256 validUntil, uint256 playerId, uint256 buyerTeamId) constant returns(bytes32)
func (_Market *MarketCallerSession) BuildOfferToBuyTxMsg(hiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, buyerTeamId *big.Int) ([32]byte, error) {
	return _Market.Contract.BuildOfferToBuyTxMsg(&_Market.CallOpts, hiddenPrice, validUntil, playerId, buyerTeamId)
}

// BuildPutAssetForSaleTxMsg is a free data retrieval call binding the contract method 0x4a74f553.
//
// Solidity: function buildPutAssetForSaleTxMsg(bytes32 hiddenPrice, uint256 validUntil, uint256 assetId) constant returns(bytes32)
func (_Market *MarketCaller) BuildPutAssetForSaleTxMsg(opts *bind.CallOpts, hiddenPrice [32]byte, validUntil *big.Int, assetId *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "buildPutAssetForSaleTxMsg", hiddenPrice, validUntil, assetId)
	return *ret0, err
}

// BuildPutAssetForSaleTxMsg is a free data retrieval call binding the contract method 0x4a74f553.
//
// Solidity: function buildPutAssetForSaleTxMsg(bytes32 hiddenPrice, uint256 validUntil, uint256 assetId) constant returns(bytes32)
func (_Market *MarketSession) BuildPutAssetForSaleTxMsg(hiddenPrice [32]byte, validUntil *big.Int, assetId *big.Int) ([32]byte, error) {
	return _Market.Contract.BuildPutAssetForSaleTxMsg(&_Market.CallOpts, hiddenPrice, validUntil, assetId)
}

// BuildPutAssetForSaleTxMsg is a free data retrieval call binding the contract method 0x4a74f553.
//
// Solidity: function buildPutAssetForSaleTxMsg(bytes32 hiddenPrice, uint256 validUntil, uint256 assetId) constant returns(bytes32)
func (_Market *MarketCallerSession) BuildPutAssetForSaleTxMsg(hiddenPrice [32]byte, validUntil *big.Int, assetId *big.Int) ([32]byte, error) {
	return _Market.Contract.BuildPutAssetForSaleTxMsg(&_Market.CallOpts, hiddenPrice, validUntil, assetId)
}

// GetBlockchainNowTime is a free data retrieval call binding the contract method 0x663c2e4b.
//
// Solidity: function getBlockchainNowTime() constant returns(uint256)
func (_Market *MarketCaller) GetBlockchainNowTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "getBlockchainNowTime")
	return *ret0, err
}

// GetBlockchainNowTime is a free data retrieval call binding the contract method 0x663c2e4b.
//
// Solidity: function getBlockchainNowTime() constant returns(uint256)
func (_Market *MarketSession) GetBlockchainNowTime() (*big.Int, error) {
	return _Market.Contract.GetBlockchainNowTime(&_Market.CallOpts)
}

// GetBlockchainNowTime is a free data retrieval call binding the contract method 0x663c2e4b.
//
// Solidity: function getBlockchainNowTime() constant returns(uint256)
func (_Market *MarketCallerSession) GetBlockchainNowTime() (*big.Int, error) {
	return _Market.Contract.GetBlockchainNowTime(&_Market.CallOpts)
}

// HashPrivateMsg is a free data retrieval call binding the contract method 0x63ef7460.
//
// Solidity: function hashPrivateMsg(uint8 currencyId, uint256 price, uint256 rnd) constant returns(bytes32)
func (_Market *MarketCaller) HashPrivateMsg(opts *bind.CallOpts, currencyId uint8, price *big.Int, rnd *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "hashPrivateMsg", currencyId, price, rnd)
	return *ret0, err
}

// HashPrivateMsg is a free data retrieval call binding the contract method 0x63ef7460.
//
// Solidity: function hashPrivateMsg(uint8 currencyId, uint256 price, uint256 rnd) constant returns(bytes32)
func (_Market *MarketSession) HashPrivateMsg(currencyId uint8, price *big.Int, rnd *big.Int) ([32]byte, error) {
	return _Market.Contract.HashPrivateMsg(&_Market.CallOpts, currencyId, price, rnd)
}

// HashPrivateMsg is a free data retrieval call binding the contract method 0x63ef7460.
//
// Solidity: function hashPrivateMsg(uint8 currencyId, uint256 price, uint256 rnd) constant returns(bytes32)
func (_Market *MarketCallerSession) HashPrivateMsg(currencyId uint8, price *big.Int, rnd *big.Int) ([32]byte, error) {
	return _Market.Contract.HashPrivateMsg(&_Market.CallOpts, currencyId, price, rnd)
}

// IsPlayerFrozen is a free data retrieval call binding the contract method 0x0db2b6ae.
//
// Solidity: function isPlayerFrozen(uint256 playerId) constant returns(bool)
func (_Market *MarketCaller) IsPlayerFrozen(opts *bind.CallOpts, playerId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "isPlayerFrozen", playerId)
	return *ret0, err
}

// IsPlayerFrozen is a free data retrieval call binding the contract method 0x0db2b6ae.
//
// Solidity: function isPlayerFrozen(uint256 playerId) constant returns(bool)
func (_Market *MarketSession) IsPlayerFrozen(playerId *big.Int) (bool, error) {
	return _Market.Contract.IsPlayerFrozen(&_Market.CallOpts, playerId)
}

// IsPlayerFrozen is a free data retrieval call binding the contract method 0x0db2b6ae.
//
// Solidity: function isPlayerFrozen(uint256 playerId) constant returns(bool)
func (_Market *MarketCallerSession) IsPlayerFrozen(playerId *big.Int) (bool, error) {
	return _Market.Contract.IsPlayerFrozen(&_Market.CallOpts, playerId)
}

// IsSigned is a free data retrieval call binding the contract method 0x8677ebe8.
//
// Solidity: function isSigned(address _addr, bytes32 msgHash, uint8 v, bytes32 r, bytes32 s) constant returns(bool)
func (_Market *MarketCaller) IsSigned(opts *bind.CallOpts, _addr common.Address, msgHash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "isSigned", _addr, msgHash, v, r, s)
	return *ret0, err
}

// IsSigned is a free data retrieval call binding the contract method 0x8677ebe8.
//
// Solidity: function isSigned(address _addr, bytes32 msgHash, uint8 v, bytes32 r, bytes32 s) constant returns(bool)
func (_Market *MarketSession) IsSigned(_addr common.Address, msgHash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Market.Contract.IsSigned(&_Market.CallOpts, _addr, msgHash, v, r, s)
}

// IsSigned is a free data retrieval call binding the contract method 0x8677ebe8.
//
// Solidity: function isSigned(address _addr, bytes32 msgHash, uint8 v, bytes32 r, bytes32 s) constant returns(bool)
func (_Market *MarketCallerSession) IsSigned(_addr common.Address, msgHash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Market.Contract.IsSigned(&_Market.CallOpts, _addr, msgHash, v, r, s)
}

// IsTeamFrozen is a free data retrieval call binding the contract method 0xe7ff8b0e.
//
// Solidity: function isTeamFrozen(uint256 teamId) constant returns(bool)
func (_Market *MarketCaller) IsTeamFrozen(opts *bind.CallOpts, teamId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "isTeamFrozen", teamId)
	return *ret0, err
}

// IsTeamFrozen is a free data retrieval call binding the contract method 0xe7ff8b0e.
//
// Solidity: function isTeamFrozen(uint256 teamId) constant returns(bool)
func (_Market *MarketSession) IsTeamFrozen(teamId *big.Int) (bool, error) {
	return _Market.Contract.IsTeamFrozen(&_Market.CallOpts, teamId)
}

// IsTeamFrozen is a free data retrieval call binding the contract method 0xe7ff8b0e.
//
// Solidity: function isTeamFrozen(uint256 teamId) constant returns(bool)
func (_Market *MarketCallerSession) IsTeamFrozen(teamId *big.Int) (bool, error) {
	return _Market.Contract.IsTeamFrozen(&_Market.CallOpts, teamId)
}

// Prefixed is a free data retrieval call binding the contract method 0x1f275713.
//
// Solidity: function prefixed(bytes32 hash) constant returns(bytes32)
func (_Market *MarketCaller) Prefixed(opts *bind.CallOpts, hash [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Market.contract.Call(opts, out, "prefixed", hash)
	return *ret0, err
}

// Prefixed is a free data retrieval call binding the contract method 0x1f275713.
//
// Solidity: function prefixed(bytes32 hash) constant returns(bytes32)
func (_Market *MarketSession) Prefixed(hash [32]byte) ([32]byte, error) {
	return _Market.Contract.Prefixed(&_Market.CallOpts, hash)
}

// Prefixed is a free data retrieval call binding the contract method 0x1f275713.
//
// Solidity: function prefixed(bytes32 hash) constant returns(bytes32)
func (_Market *MarketCallerSession) Prefixed(hash [32]byte) ([32]byte, error) {
	return _Market.Contract.Prefixed(&_Market.CallOpts, hash)
}

// CompletePlayerAuction is a paid mutator transaction binding the contract method 0x60f4818e.
//
// Solidity: function completePlayerAuction(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) returns()
func (_Market *MarketTransactor) CompletePlayerAuction(opts *bind.TransactOpts, sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, buyerHiddenPrice [32]byte, buyerTeamId *big.Int, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (*types.Transaction, error) {
	return _Market.contract.Transact(opts, "completePlayerAuction", sellerHiddenPrice, validUntil, playerId, buyerHiddenPrice, buyerTeamId, sig, sigV, isOffer2StartAuction)
}

// CompletePlayerAuction is a paid mutator transaction binding the contract method 0x60f4818e.
//
// Solidity: function completePlayerAuction(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) returns()
func (_Market *MarketSession) CompletePlayerAuction(sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, buyerHiddenPrice [32]byte, buyerTeamId *big.Int, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (*types.Transaction, error) {
	return _Market.Contract.CompletePlayerAuction(&_Market.TransactOpts, sellerHiddenPrice, validUntil, playerId, buyerHiddenPrice, buyerTeamId, sig, sigV, isOffer2StartAuction)
}

// CompletePlayerAuction is a paid mutator transaction binding the contract method 0x60f4818e.
//
// Solidity: function completePlayerAuction(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32 buyerHiddenPrice, uint256 buyerTeamId, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) returns()
func (_Market *MarketTransactorSession) CompletePlayerAuction(sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, buyerHiddenPrice [32]byte, buyerTeamId *big.Int, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (*types.Transaction, error) {
	return _Market.Contract.CompletePlayerAuction(&_Market.TransactOpts, sellerHiddenPrice, validUntil, playerId, buyerHiddenPrice, buyerTeamId, sig, sigV, isOffer2StartAuction)
}

// CompleteTeamAuction is a paid mutator transaction binding the contract method 0xa51dcaa3.
//
// Solidity: function completeTeamAuction(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32 buyerHiddenPrice, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) returns()
func (_Market *MarketTransactor) CompleteTeamAuction(opts *bind.TransactOpts, sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, buyerHiddenPrice [32]byte, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (*types.Transaction, error) {
	return _Market.contract.Transact(opts, "completeTeamAuction", sellerHiddenPrice, validUntil, teamId, buyerHiddenPrice, sig, sigV, isOffer2StartAuction)
}

// CompleteTeamAuction is a paid mutator transaction binding the contract method 0xa51dcaa3.
//
// Solidity: function completeTeamAuction(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32 buyerHiddenPrice, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) returns()
func (_Market *MarketSession) CompleteTeamAuction(sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, buyerHiddenPrice [32]byte, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (*types.Transaction, error) {
	return _Market.Contract.CompleteTeamAuction(&_Market.TransactOpts, sellerHiddenPrice, validUntil, teamId, buyerHiddenPrice, sig, sigV, isOffer2StartAuction)
}

// CompleteTeamAuction is a paid mutator transaction binding the contract method 0xa51dcaa3.
//
// Solidity: function completeTeamAuction(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32 buyerHiddenPrice, bytes32[3] sig, uint8 sigV, bool isOffer2StartAuction) returns()
func (_Market *MarketTransactorSession) CompleteTeamAuction(sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, buyerHiddenPrice [32]byte, sig [3][32]byte, sigV uint8, isOffer2StartAuction bool) (*types.Transaction, error) {
	return _Market.Contract.CompleteTeamAuction(&_Market.TransactOpts, sellerHiddenPrice, validUntil, teamId, buyerHiddenPrice, sig, sigV, isOffer2StartAuction)
}

// FreezePlayer is a paid mutator transaction binding the contract method 0x5e8f77e7.
//
// Solidity: function freezePlayer(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32[3] sig, uint8 sigV) returns()
func (_Market *MarketTransactor) FreezePlayer(opts *bind.TransactOpts, sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, sig [3][32]byte, sigV uint8) (*types.Transaction, error) {
	return _Market.contract.Transact(opts, "freezePlayer", sellerHiddenPrice, validUntil, playerId, sig, sigV)
}

// FreezePlayer is a paid mutator transaction binding the contract method 0x5e8f77e7.
//
// Solidity: function freezePlayer(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32[3] sig, uint8 sigV) returns()
func (_Market *MarketSession) FreezePlayer(sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, sig [3][32]byte, sigV uint8) (*types.Transaction, error) {
	return _Market.Contract.FreezePlayer(&_Market.TransactOpts, sellerHiddenPrice, validUntil, playerId, sig, sigV)
}

// FreezePlayer is a paid mutator transaction binding the contract method 0x5e8f77e7.
//
// Solidity: function freezePlayer(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 playerId, bytes32[3] sig, uint8 sigV) returns()
func (_Market *MarketTransactorSession) FreezePlayer(sellerHiddenPrice [32]byte, validUntil *big.Int, playerId *big.Int, sig [3][32]byte, sigV uint8) (*types.Transaction, error) {
	return _Market.Contract.FreezePlayer(&_Market.TransactOpts, sellerHiddenPrice, validUntil, playerId, sig, sigV)
}

// FreezeTeam is a paid mutator transaction binding the contract method 0x75a55471.
//
// Solidity: function freezeTeam(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32[3] sig, uint8 sigV) returns()
func (_Market *MarketTransactor) FreezeTeam(opts *bind.TransactOpts, sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, sig [3][32]byte, sigV uint8) (*types.Transaction, error) {
	return _Market.contract.Transact(opts, "freezeTeam", sellerHiddenPrice, validUntil, teamId, sig, sigV)
}

// FreezeTeam is a paid mutator transaction binding the contract method 0x75a55471.
//
// Solidity: function freezeTeam(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32[3] sig, uint8 sigV) returns()
func (_Market *MarketSession) FreezeTeam(sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, sig [3][32]byte, sigV uint8) (*types.Transaction, error) {
	return _Market.Contract.FreezeTeam(&_Market.TransactOpts, sellerHiddenPrice, validUntil, teamId, sig, sigV)
}

// FreezeTeam is a paid mutator transaction binding the contract method 0x75a55471.
//
// Solidity: function freezeTeam(bytes32 sellerHiddenPrice, uint256 validUntil, uint256 teamId, bytes32[3] sig, uint8 sigV) returns()
func (_Market *MarketTransactorSession) FreezeTeam(sellerHiddenPrice [32]byte, validUntil *big.Int, teamId *big.Int, sig [3][32]byte, sigV uint8) (*types.Transaction, error) {
	return _Market.Contract.FreezeTeam(&_Market.TransactOpts, sellerHiddenPrice, validUntil, teamId, sig, sigV)
}

// SetAssetsAddress is a paid mutator transaction binding the contract method 0x1490a174.
//
// Solidity: function setAssetsAddress(address addr) returns()
func (_Market *MarketTransactor) SetAssetsAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _Market.contract.Transact(opts, "setAssetsAddress", addr)
}

// SetAssetsAddress is a paid mutator transaction binding the contract method 0x1490a174.
//
// Solidity: function setAssetsAddress(address addr) returns()
func (_Market *MarketSession) SetAssetsAddress(addr common.Address) (*types.Transaction, error) {
	return _Market.Contract.SetAssetsAddress(&_Market.TransactOpts, addr)
}

// SetAssetsAddress is a paid mutator transaction binding the contract method 0x1490a174.
//
// Solidity: function setAssetsAddress(address addr) returns()
func (_Market *MarketTransactorSession) SetAssetsAddress(addr common.Address) (*types.Transaction, error) {
	return _Market.Contract.SetAssetsAddress(&_Market.TransactOpts, addr)
}

// MarketPlayerFreezeIterator is returned from FilterPlayerFreeze and is used to iterate over the raw logs and unpacked data for PlayerFreeze events raised by the Market contract.
type MarketPlayerFreezeIterator struct {
	Event *MarketPlayerFreeze // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MarketPlayerFreezeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MarketPlayerFreeze)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MarketPlayerFreeze)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MarketPlayerFreezeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MarketPlayerFreezeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MarketPlayerFreeze represents a PlayerFreeze event raised by the Market contract.
type MarketPlayerFreeze struct {
	PlayerId    *big.Int
	AuctionData *big.Int
	Frozen      bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPlayerFreeze is a free log retrieval operation binding the contract event 0x44b30f34d8f6f8c3cc8737fe3476b3bfd6fff21e03ef669b48b787a213b1f110.
//
// Solidity: event PlayerFreeze(uint256 playerId, uint256 auctionData, bool frozen)
func (_Market *MarketFilterer) FilterPlayerFreeze(opts *bind.FilterOpts) (*MarketPlayerFreezeIterator, error) {

	logs, sub, err := _Market.contract.FilterLogs(opts, "PlayerFreeze")
	if err != nil {
		return nil, err
	}
	return &MarketPlayerFreezeIterator{contract: _Market.contract, event: "PlayerFreeze", logs: logs, sub: sub}, nil
}

// WatchPlayerFreeze is a free log subscription operation binding the contract event 0x44b30f34d8f6f8c3cc8737fe3476b3bfd6fff21e03ef669b48b787a213b1f110.
//
// Solidity: event PlayerFreeze(uint256 playerId, uint256 auctionData, bool frozen)
func (_Market *MarketFilterer) WatchPlayerFreeze(opts *bind.WatchOpts, sink chan<- *MarketPlayerFreeze) (event.Subscription, error) {

	logs, sub, err := _Market.contract.WatchLogs(opts, "PlayerFreeze")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MarketPlayerFreeze)
				if err := _Market.contract.UnpackLog(event, "PlayerFreeze", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MarketTeamFreezeIterator is returned from FilterTeamFreeze and is used to iterate over the raw logs and unpacked data for TeamFreeze events raised by the Market contract.
type MarketTeamFreezeIterator struct {
	Event *MarketTeamFreeze // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MarketTeamFreezeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MarketTeamFreeze)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MarketTeamFreeze)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MarketTeamFreezeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MarketTeamFreezeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MarketTeamFreeze represents a TeamFreeze event raised by the Market contract.
type MarketTeamFreeze struct {
	TeamId      *big.Int
	AuctionData *big.Int
	Frozen      bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTeamFreeze is a free log retrieval operation binding the contract event 0xbe5cb47c7008dd92757f872cf47b2f27eae3e4d3efbba0bd54969ed71c927d0e.
//
// Solidity: event TeamFreeze(uint256 teamId, uint256 auctionData, bool frozen)
func (_Market *MarketFilterer) FilterTeamFreeze(opts *bind.FilterOpts) (*MarketTeamFreezeIterator, error) {

	logs, sub, err := _Market.contract.FilterLogs(opts, "TeamFreeze")
	if err != nil {
		return nil, err
	}
	return &MarketTeamFreezeIterator{contract: _Market.contract, event: "TeamFreeze", logs: logs, sub: sub}, nil
}

// WatchTeamFreeze is a free log subscription operation binding the contract event 0xbe5cb47c7008dd92757f872cf47b2f27eae3e4d3efbba0bd54969ed71c927d0e.
//
// Solidity: event TeamFreeze(uint256 teamId, uint256 auctionData, bool frozen)
func (_Market *MarketFilterer) WatchTeamFreeze(opts *bind.WatchOpts, sink chan<- *MarketTeamFreeze) (event.Subscription, error) {

	logs, sub, err := _Market.contract.WatchLogs(opts, "TeamFreeze")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MarketTeamFreeze)
				if err := _Market.contract.UnpackLog(event, "TeamFreeze", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
