// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package updates

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

// UpdatesABI is the input ABI used to generate the binding from.
const UpdatesABI = "[{\"inputs\":[],\"constant\":true,\"name\":\"SECS_BETWEEN_VERSES\",\"outputs\":[{\"type\":\"uint16\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"signature\":\"0x28116e59\",\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"timeZoneForRound1\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"signature\":\"0x61703da5\",\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"NULL_TIMEZONE\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"signature\":\"0xc52fd716\",\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"nextVerseTimestamp\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"signature\":\"0xe28d3a50\",\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"currentVerse\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"signature\":\"0xe97696b2\",\"type\":\"function\"},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"teamId\"},{\"indexed\":false,\"type\":\"address\",\"name\":\"to\"}],\"type\":\"event\",\"name\":\"TeamTransfer\",\"anonymous\":false,\"signature\":\"0x77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf38\"},{\"inputs\":[{\"indexed\":false,\"type\":\"uint8\",\"name\":\"timeZone\"},{\"indexed\":false,\"type\":\"uint8\",\"name\":\"day\"},{\"indexed\":false,\"type\":\"uint8\",\"name\":\"turnInDay\"},{\"indexed\":false,\"type\":\"bytes32\",\"name\":\"seed\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"submissionTime\"}],\"type\":\"event\",\"name\":\"ActionsSubmission\",\"anonymous\":false,\"signature\":\"0x3eaf8623b39cd8460e98a71c435b5fea0d21dc62cc5588b0e628edf98b5db685\"},{\"inputs\":[{\"indexed\":false,\"type\":\"uint8\",\"name\":\"timeZone\"},{\"indexed\":false,\"type\":\"bytes32\",\"name\":\"root\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"submissionTime\"}],\"type\":\"event\",\"name\":\"TimeZoneUpdate\",\"anonymous\":false,\"signature\":\"0x21910d8429bd1db0b36330d1ac1df9f94ff1eeae897e11db418e30bc0fb418b9\"},{\"inputs\":[{\"type\":\"address\",\"name\":\"addr\"}],\"constant\":false,\"name\":\"initUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"signature\":\"0x098ef280\",\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"getNow\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"signature\":\"0xbbe4fd50\",\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"actionsRoot\"}],\"constant\":false,\"name\":\"submitActionsRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"signature\":\"0xe8507051\",\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"root\"}],\"constant\":false,\"name\":\"updateTZ\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"signature\":\"0x16b5d047\",\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"nextTimeZoneToUpdate\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint8\",\"name\":\"day\"},{\"type\":\"uint8\",\"name\":\"turnInDay\"}],\"stateMutability\":\"view\",\"payable\":false,\"signature\":\"0x8a89148c\",\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"prevTimeZoneToUpdate\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint8\",\"name\":\"day\"},{\"type\":\"uint8\",\"name\":\"turnInDay\"}],\"stateMutability\":\"view\",\"payable\":false,\"signature\":\"0xc395d7a8\",\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"verse\"},{\"type\":\"uint8\",\"name\":\"TZForRound1\"}],\"constant\":true,\"name\":\"_timeZoneToUpdatePure\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint8\",\"name\":\"day\"},{\"type\":\"uint8\",\"name\":\"turnInDay\"}],\"stateMutability\":\"pure\",\"payable\":false,\"signature\":\"0xeaa9e98e\",\"type\":\"function\"},{\"inputs\":[{\"type\":\"bytes32\",\"name\":\"seed\"}],\"constant\":false,\"name\":\"setCurrentVerseSeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"signature\":\"0x50c0f2af\",\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"getCurrentVerseSeed\",\"outputs\":[{\"type\":\"bytes32\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"signature\":\"0x21eab316\",\"type\":\"function\"}]"

// UpdatesBin is the compiled bytecode used for deploying new contracts.
const UpdatesBin = `0x60806040526001600460006101000a81548160ff02191690831515021790555034801561002b57600080fd5b50610d3a8061003b6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063bbe4fd501161008c578063e28d3a5011610066578063e28d3a50146102b5578063e8507051146102d3578063e97696b214610301578063eaa9e98e1461031f576100ea565b8063bbe4fd5014610235578063c395d7a814610253578063c52fd71614610291576100ea565b806328116e59116100c857806328116e591461017f57806350c0f2af146101a557806361703da5146101d35780638a89148c146101f7576100ea565b8063098ef280146100ef57806316b5d0471461013357806321eab31614610161575b600080fd5b6101316004803603602081101561010557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061038e565b005b61015f6004803603602081101561014957600080fd5b81019080803590602001909291905050506104ff565b005b6101696108b9565b6040518082815260200191505060405180910390f35b6101876108c3565b604051808261ffff1661ffff16815260200191505060405180910390f35b6101d1600480360360208110156101bb57600080fd5b81019080803590602001909291905050506108c9565b005b6101db6108d3565b604051808260ff1660ff16815260200191505060405180910390f35b6101ff6108e6565b604051808460ff1660ff1681526020018360ff1660ff1681526020018260ff1660ff168152602001935050505060405180910390f35b61023d610911565b6040518082815260200191505060405180910390f35b61025b610919565b604051808460ff1660ff1681526020018360ff1660ff1681526020018260ff1660ff168152602001935050505060405180910390f35b610299610968565b604051808260ff1660ff16815260200191505060405180910390f35b6102bd61096d565b6040518082815260200191505060405180910390f35b6102ff600480360360208110156102e957600080fd5b8101908080359060200190929190505050610973565b005b610309610b55565b6040518082815260200191505060405180910390f35b6103586004803603604081101561033557600080fd5b8101908080359060200190929190803560ff169060200190929190505050610b5b565b604051808460ff1660ff1681526020018360ff1660ff1681526020018260ff1660ff168152602001935050505060405180910390f35b60011515600460009054906101000a900460ff16151514610417576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f63616e6e6f7420696e697469616c697a6520747769636500000000000000000081525060200191505060405180910390fd5b61042081610c45565b600062015180428161042e57fe5b0690506000610e10828161043e57fe5b0490506000603c610e10830284038161045357fe5b0490506000603c8202610e1084028503039050602a8210156104a65782600101600160006101000a81548160ff021916908360ff16021790555080603c03603c83602c03024201016000819055506104dd565b82600201600160006101000a81548160ff021916908360ff160217905550610e1081603c03603c84602c0302420101016000819055505b6000600460006101000a81548160ff0219169083151502179055505050505050565b6000610509610919565b50509050600060ff168160ff16141561056d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526029815260200180610cb26029913960400191505060405180910390fd5b6000600460019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632d0e08fd836040518263ffffffff1660e01b8152600401808260ff1660ff16815260200191505060206040518083038186803b1580156105e857600080fd5b505afa1580156105fc573d6000803e3d6000fd5b505050506040513d602081101561061257600080fd5b810190808051906020019092919050505090506000600460019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fa80039b846040518263ffffffff1660e01b8152600401808260ff1660ff16815260200191505060206040518083038186803b1580156106a057600080fd5b505afa1580156106b4573d6000803e3d6000fd5b505050506040513d60208110156106ca57600080fd5b810190808051906020019092919050505090508082111561074857603c60ff1682014210610743576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526034815260200180610cdb6034913960400191505060405180910390fd5b6107a7565b603c60ff16810142106107a6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526034815260200180610cdb6034913960400191505060405180910390fd5b5b600460019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ec1c542384866040518363ffffffff1660e01b8152600401808360ff1660ff16815260200182815260200192505050602060405180830381600087803b15801561082a57600080fd5b505af115801561083e573d6000803e3d6000fd5b505050506040513d602081101561085457600080fd5b8101908080519060200190929190505050507f21910d8429bd1db0b36330d1ac1df9f94ff1eeae897e11db418e30bc0fb418b9838542604051808460ff1660ff168152602001838152602001828152602001935050505060405180910390a150505050565b6000600354905090565b61038481565b8060038190555050565b600160009054906101000a900460ff1681565b6000806000610906600254600160009054906101000a900460ff16610b5b565b925092509250909192565b600042905090565b600080600080600254141561093e576000806000819150809050925092509250610963565b61095c600160025403600160009054906101000a900460ff16610b5b565b9250925092505b909192565b600081565b60005481565b60008060006109806108e6565b925092509250600060ff168360ff161415610a0d5761099d610c89565b7f3eaf8623b39cd8460e98a71c435b5fea0d21dc62cc5588b0e628edf98b5db68560008060008042604051808660ff1660ff1681526020018560ff1681526020018460ff1681526020018360001b81526020018281526020019550505050505060405180910390a1505050610b52565b600460019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663dba6319e84866040518363ffffffff1660e01b8152600401808360ff1660ff16815260200182815260200192505050602060405180830381600087803b158015610a9057600080fd5b505af1158015610aa4573d6000803e3d6000fd5b505050506040513d6020811015610aba57600080fd5b810190808051906020019092919050505050610ad4610c89565b610ae160014303406108c9565b7f3eaf8623b39cd8460e98a71c435b5fea0d21dc62cc5588b0e628edf98b5db685838383600143034042604051808660ff1660ff1681526020018560ff1660ff1681526020018460ff1660ff1681526020018381526020018281526020019550505050505060405180910390a15050505b50565b60025481565b60008060008061060061ffff168681610b7057fe5b069050606060ff168161ffff161015610bc757601860048261ffff1681610b9357fe5b046001870360ff160161ffff1681610ba757fe5b0660010193506001925060048161ffff1681610bbf57fe5b069150610c3d565b606060ff168161ffff161415610be05760009350610c3c565b601860046001830361ffff1681610bf357fe5b046001870360ff160161ffff1681610c0757fe5b066001019350606060ff166001820361ffff1681610c2157fe5b04600101925060046001820361ffff1681610c3857fe5b0691505b5b509250925092565b80600460016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600160026000828254019250508190555061038461ffff16600080828254019250508190555056fe6e6f7468696e6720746f2075706461746520696e207468652063757272656e742074696d655a6f6e656368616c6c656e67696e6720706572696f6420697320616c7265616479206f76657220666f7220746869732074696d657a6f6e65a165627a7a723058207412877f6284018f2be93572c6bf026f35e554524147a0f548cb50293c84c2400029`

// DeployUpdates deploys a new Ethereum contract, binding an instance of Updates to it.
func DeployUpdates(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Updates, error) {
	parsed, err := abi.JSON(strings.NewReader(UpdatesABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(UpdatesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Updates{UpdatesCaller: UpdatesCaller{contract: contract}, UpdatesTransactor: UpdatesTransactor{contract: contract}, UpdatesFilterer: UpdatesFilterer{contract: contract}}, nil
}

// Updates is an auto generated Go binding around an Ethereum contract.
type Updates struct {
	UpdatesCaller     // Read-only binding to the contract
	UpdatesTransactor // Write-only binding to the contract
	UpdatesFilterer   // Log filterer for contract events
}

// UpdatesCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpdatesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpdatesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpdatesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpdatesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpdatesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpdatesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpdatesSession struct {
	Contract     *Updates          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpdatesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpdatesCallerSession struct {
	Contract *UpdatesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// UpdatesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpdatesTransactorSession struct {
	Contract     *UpdatesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// UpdatesRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpdatesRaw struct {
	Contract *Updates // Generic contract binding to access the raw methods on
}

// UpdatesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpdatesCallerRaw struct {
	Contract *UpdatesCaller // Generic read-only contract binding to access the raw methods on
}

// UpdatesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpdatesTransactorRaw struct {
	Contract *UpdatesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpdates creates a new instance of Updates, bound to a specific deployed contract.
func NewUpdates(address common.Address, backend bind.ContractBackend) (*Updates, error) {
	contract, err := bindUpdates(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Updates{UpdatesCaller: UpdatesCaller{contract: contract}, UpdatesTransactor: UpdatesTransactor{contract: contract}, UpdatesFilterer: UpdatesFilterer{contract: contract}}, nil
}

// NewUpdatesCaller creates a new read-only instance of Updates, bound to a specific deployed contract.
func NewUpdatesCaller(address common.Address, caller bind.ContractCaller) (*UpdatesCaller, error) {
	contract, err := bindUpdates(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpdatesCaller{contract: contract}, nil
}

// NewUpdatesTransactor creates a new write-only instance of Updates, bound to a specific deployed contract.
func NewUpdatesTransactor(address common.Address, transactor bind.ContractTransactor) (*UpdatesTransactor, error) {
	contract, err := bindUpdates(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpdatesTransactor{contract: contract}, nil
}

// NewUpdatesFilterer creates a new log filterer instance of Updates, bound to a specific deployed contract.
func NewUpdatesFilterer(address common.Address, filterer bind.ContractFilterer) (*UpdatesFilterer, error) {
	contract, err := bindUpdates(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpdatesFilterer{contract: contract}, nil
}

// bindUpdates binds a generic wrapper to an already deployed contract.
func bindUpdates(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpdatesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Updates *UpdatesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Updates.Contract.UpdatesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Updates *UpdatesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Updates.Contract.UpdatesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Updates *UpdatesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Updates.Contract.UpdatesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Updates *UpdatesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Updates.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Updates *UpdatesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Updates.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Updates *UpdatesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Updates.Contract.contract.Transact(opts, method, params...)
}

// NULLTIMEZONE is a free data retrieval call binding the contract method 0xc52fd716.
//
// Solidity: function NULL_TIMEZONE() constant returns(uint8)
func (_Updates *UpdatesCaller) NULLTIMEZONE(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Updates.contract.Call(opts, out, "NULL_TIMEZONE")
	return *ret0, err
}

// NULLTIMEZONE is a free data retrieval call binding the contract method 0xc52fd716.
//
// Solidity: function NULL_TIMEZONE() constant returns(uint8)
func (_Updates *UpdatesSession) NULLTIMEZONE() (uint8, error) {
	return _Updates.Contract.NULLTIMEZONE(&_Updates.CallOpts)
}

// NULLTIMEZONE is a free data retrieval call binding the contract method 0xc52fd716.
//
// Solidity: function NULL_TIMEZONE() constant returns(uint8)
func (_Updates *UpdatesCallerSession) NULLTIMEZONE() (uint8, error) {
	return _Updates.Contract.NULLTIMEZONE(&_Updates.CallOpts)
}

// SECSBETWEENVERSES is a free data retrieval call binding the contract method 0x28116e59.
//
// Solidity: function SECS_BETWEEN_VERSES() constant returns(uint16)
func (_Updates *UpdatesCaller) SECSBETWEENVERSES(opts *bind.CallOpts) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _Updates.contract.Call(opts, out, "SECS_BETWEEN_VERSES")
	return *ret0, err
}

// SECSBETWEENVERSES is a free data retrieval call binding the contract method 0x28116e59.
//
// Solidity: function SECS_BETWEEN_VERSES() constant returns(uint16)
func (_Updates *UpdatesSession) SECSBETWEENVERSES() (uint16, error) {
	return _Updates.Contract.SECSBETWEENVERSES(&_Updates.CallOpts)
}

// SECSBETWEENVERSES is a free data retrieval call binding the contract method 0x28116e59.
//
// Solidity: function SECS_BETWEEN_VERSES() constant returns(uint16)
func (_Updates *UpdatesCallerSession) SECSBETWEENVERSES() (uint16, error) {
	return _Updates.Contract.SECSBETWEENVERSES(&_Updates.CallOpts)
}

// TimeZoneToUpdatePure is a free data retrieval call binding the contract method 0xeaa9e98e.
//
// Solidity: function _timeZoneToUpdatePure(uint256 verse, uint8 TZForRound1) constant returns(uint8 timeZone, uint8 day, uint8 turnInDay)
func (_Updates *UpdatesCaller) TimeZoneToUpdatePure(opts *bind.CallOpts, verse *big.Int, TZForRound1 uint8) (struct {
	TimeZone  uint8
	Day       uint8
	TurnInDay uint8
}, error) {
	ret := new(struct {
		TimeZone  uint8
		Day       uint8
		TurnInDay uint8
	})
	out := ret
	err := _Updates.contract.Call(opts, out, "_timeZoneToUpdatePure", verse, TZForRound1)
	return *ret, err
}

// TimeZoneToUpdatePure is a free data retrieval call binding the contract method 0xeaa9e98e.
//
// Solidity: function _timeZoneToUpdatePure(uint256 verse, uint8 TZForRound1) constant returns(uint8 timeZone, uint8 day, uint8 turnInDay)
func (_Updates *UpdatesSession) TimeZoneToUpdatePure(verse *big.Int, TZForRound1 uint8) (struct {
	TimeZone  uint8
	Day       uint8
	TurnInDay uint8
}, error) {
	return _Updates.Contract.TimeZoneToUpdatePure(&_Updates.CallOpts, verse, TZForRound1)
}

// TimeZoneToUpdatePure is a free data retrieval call binding the contract method 0xeaa9e98e.
//
// Solidity: function _timeZoneToUpdatePure(uint256 verse, uint8 TZForRound1) constant returns(uint8 timeZone, uint8 day, uint8 turnInDay)
func (_Updates *UpdatesCallerSession) TimeZoneToUpdatePure(verse *big.Int, TZForRound1 uint8) (struct {
	TimeZone  uint8
	Day       uint8
	TurnInDay uint8
}, error) {
	return _Updates.Contract.TimeZoneToUpdatePure(&_Updates.CallOpts, verse, TZForRound1)
}

// CurrentVerse is a free data retrieval call binding the contract method 0xe97696b2.
//
// Solidity: function currentVerse() constant returns(uint256)
func (_Updates *UpdatesCaller) CurrentVerse(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Updates.contract.Call(opts, out, "currentVerse")
	return *ret0, err
}

// CurrentVerse is a free data retrieval call binding the contract method 0xe97696b2.
//
// Solidity: function currentVerse() constant returns(uint256)
func (_Updates *UpdatesSession) CurrentVerse() (*big.Int, error) {
	return _Updates.Contract.CurrentVerse(&_Updates.CallOpts)
}

// CurrentVerse is a free data retrieval call binding the contract method 0xe97696b2.
//
// Solidity: function currentVerse() constant returns(uint256)
func (_Updates *UpdatesCallerSession) CurrentVerse() (*big.Int, error) {
	return _Updates.Contract.CurrentVerse(&_Updates.CallOpts)
}

// GetCurrentVerseSeed is a free data retrieval call binding the contract method 0x21eab316.
//
// Solidity: function getCurrentVerseSeed() constant returns(bytes32)
func (_Updates *UpdatesCaller) GetCurrentVerseSeed(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Updates.contract.Call(opts, out, "getCurrentVerseSeed")
	return *ret0, err
}

// GetCurrentVerseSeed is a free data retrieval call binding the contract method 0x21eab316.
//
// Solidity: function getCurrentVerseSeed() constant returns(bytes32)
func (_Updates *UpdatesSession) GetCurrentVerseSeed() ([32]byte, error) {
	return _Updates.Contract.GetCurrentVerseSeed(&_Updates.CallOpts)
}

// GetCurrentVerseSeed is a free data retrieval call binding the contract method 0x21eab316.
//
// Solidity: function getCurrentVerseSeed() constant returns(bytes32)
func (_Updates *UpdatesCallerSession) GetCurrentVerseSeed() ([32]byte, error) {
	return _Updates.Contract.GetCurrentVerseSeed(&_Updates.CallOpts)
}

// GetNow is a free data retrieval call binding the contract method 0xbbe4fd50.
//
// Solidity: function getNow() constant returns(uint256)
func (_Updates *UpdatesCaller) GetNow(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Updates.contract.Call(opts, out, "getNow")
	return *ret0, err
}

// GetNow is a free data retrieval call binding the contract method 0xbbe4fd50.
//
// Solidity: function getNow() constant returns(uint256)
func (_Updates *UpdatesSession) GetNow() (*big.Int, error) {
	return _Updates.Contract.GetNow(&_Updates.CallOpts)
}

// GetNow is a free data retrieval call binding the contract method 0xbbe4fd50.
//
// Solidity: function getNow() constant returns(uint256)
func (_Updates *UpdatesCallerSession) GetNow() (*big.Int, error) {
	return _Updates.Contract.GetNow(&_Updates.CallOpts)
}

// NextTimeZoneToUpdate is a free data retrieval call binding the contract method 0x8a89148c.
//
// Solidity: function nextTimeZoneToUpdate() constant returns(uint8 timeZone, uint8 day, uint8 turnInDay)
func (_Updates *UpdatesCaller) NextTimeZoneToUpdate(opts *bind.CallOpts) (struct {
	TimeZone  uint8
	Day       uint8
	TurnInDay uint8
}, error) {
	ret := new(struct {
		TimeZone  uint8
		Day       uint8
		TurnInDay uint8
	})
	out := ret
	err := _Updates.contract.Call(opts, out, "nextTimeZoneToUpdate")
	return *ret, err
}

// NextTimeZoneToUpdate is a free data retrieval call binding the contract method 0x8a89148c.
//
// Solidity: function nextTimeZoneToUpdate() constant returns(uint8 timeZone, uint8 day, uint8 turnInDay)
func (_Updates *UpdatesSession) NextTimeZoneToUpdate() (struct {
	TimeZone  uint8
	Day       uint8
	TurnInDay uint8
}, error) {
	return _Updates.Contract.NextTimeZoneToUpdate(&_Updates.CallOpts)
}

// NextTimeZoneToUpdate is a free data retrieval call binding the contract method 0x8a89148c.
//
// Solidity: function nextTimeZoneToUpdate() constant returns(uint8 timeZone, uint8 day, uint8 turnInDay)
func (_Updates *UpdatesCallerSession) NextTimeZoneToUpdate() (struct {
	TimeZone  uint8
	Day       uint8
	TurnInDay uint8
}, error) {
	return _Updates.Contract.NextTimeZoneToUpdate(&_Updates.CallOpts)
}

// NextVerseTimestamp is a free data retrieval call binding the contract method 0xe28d3a50.
//
// Solidity: function nextVerseTimestamp() constant returns(uint256)
func (_Updates *UpdatesCaller) NextVerseTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Updates.contract.Call(opts, out, "nextVerseTimestamp")
	return *ret0, err
}

// NextVerseTimestamp is a free data retrieval call binding the contract method 0xe28d3a50.
//
// Solidity: function nextVerseTimestamp() constant returns(uint256)
func (_Updates *UpdatesSession) NextVerseTimestamp() (*big.Int, error) {
	return _Updates.Contract.NextVerseTimestamp(&_Updates.CallOpts)
}

// NextVerseTimestamp is a free data retrieval call binding the contract method 0xe28d3a50.
//
// Solidity: function nextVerseTimestamp() constant returns(uint256)
func (_Updates *UpdatesCallerSession) NextVerseTimestamp() (*big.Int, error) {
	return _Updates.Contract.NextVerseTimestamp(&_Updates.CallOpts)
}

// PrevTimeZoneToUpdate is a free data retrieval call binding the contract method 0xc395d7a8.
//
// Solidity: function prevTimeZoneToUpdate() constant returns(uint8 timeZone, uint8 day, uint8 turnInDay)
func (_Updates *UpdatesCaller) PrevTimeZoneToUpdate(opts *bind.CallOpts) (struct {
	TimeZone  uint8
	Day       uint8
	TurnInDay uint8
}, error) {
	ret := new(struct {
		TimeZone  uint8
		Day       uint8
		TurnInDay uint8
	})
	out := ret
	err := _Updates.contract.Call(opts, out, "prevTimeZoneToUpdate")
	return *ret, err
}

// PrevTimeZoneToUpdate is a free data retrieval call binding the contract method 0xc395d7a8.
//
// Solidity: function prevTimeZoneToUpdate() constant returns(uint8 timeZone, uint8 day, uint8 turnInDay)
func (_Updates *UpdatesSession) PrevTimeZoneToUpdate() (struct {
	TimeZone  uint8
	Day       uint8
	TurnInDay uint8
}, error) {
	return _Updates.Contract.PrevTimeZoneToUpdate(&_Updates.CallOpts)
}

// PrevTimeZoneToUpdate is a free data retrieval call binding the contract method 0xc395d7a8.
//
// Solidity: function prevTimeZoneToUpdate() constant returns(uint8 timeZone, uint8 day, uint8 turnInDay)
func (_Updates *UpdatesCallerSession) PrevTimeZoneToUpdate() (struct {
	TimeZone  uint8
	Day       uint8
	TurnInDay uint8
}, error) {
	return _Updates.Contract.PrevTimeZoneToUpdate(&_Updates.CallOpts)
}

// TimeZoneForRound1 is a free data retrieval call binding the contract method 0x61703da5.
//
// Solidity: function timeZoneForRound1() constant returns(uint8)
func (_Updates *UpdatesCaller) TimeZoneForRound1(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Updates.contract.Call(opts, out, "timeZoneForRound1")
	return *ret0, err
}

// TimeZoneForRound1 is a free data retrieval call binding the contract method 0x61703da5.
//
// Solidity: function timeZoneForRound1() constant returns(uint8)
func (_Updates *UpdatesSession) TimeZoneForRound1() (uint8, error) {
	return _Updates.Contract.TimeZoneForRound1(&_Updates.CallOpts)
}

// TimeZoneForRound1 is a free data retrieval call binding the contract method 0x61703da5.
//
// Solidity: function timeZoneForRound1() constant returns(uint8)
func (_Updates *UpdatesCallerSession) TimeZoneForRound1() (uint8, error) {
	return _Updates.Contract.TimeZoneForRound1(&_Updates.CallOpts)
}

// InitUpdates is a paid mutator transaction binding the contract method 0x098ef280.
//
// Solidity: function initUpdates(address addr) returns()
func (_Updates *UpdatesTransactor) InitUpdates(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _Updates.contract.Transact(opts, "initUpdates", addr)
}

// InitUpdates is a paid mutator transaction binding the contract method 0x098ef280.
//
// Solidity: function initUpdates(address addr) returns()
func (_Updates *UpdatesSession) InitUpdates(addr common.Address) (*types.Transaction, error) {
	return _Updates.Contract.InitUpdates(&_Updates.TransactOpts, addr)
}

// InitUpdates is a paid mutator transaction binding the contract method 0x098ef280.
//
// Solidity: function initUpdates(address addr) returns()
func (_Updates *UpdatesTransactorSession) InitUpdates(addr common.Address) (*types.Transaction, error) {
	return _Updates.Contract.InitUpdates(&_Updates.TransactOpts, addr)
}

// SetCurrentVerseSeed is a paid mutator transaction binding the contract method 0x50c0f2af.
//
// Solidity: function setCurrentVerseSeed(bytes32 seed) returns()
func (_Updates *UpdatesTransactor) SetCurrentVerseSeed(opts *bind.TransactOpts, seed [32]byte) (*types.Transaction, error) {
	return _Updates.contract.Transact(opts, "setCurrentVerseSeed", seed)
}

// SetCurrentVerseSeed is a paid mutator transaction binding the contract method 0x50c0f2af.
//
// Solidity: function setCurrentVerseSeed(bytes32 seed) returns()
func (_Updates *UpdatesSession) SetCurrentVerseSeed(seed [32]byte) (*types.Transaction, error) {
	return _Updates.Contract.SetCurrentVerseSeed(&_Updates.TransactOpts, seed)
}

// SetCurrentVerseSeed is a paid mutator transaction binding the contract method 0x50c0f2af.
//
// Solidity: function setCurrentVerseSeed(bytes32 seed) returns()
func (_Updates *UpdatesTransactorSession) SetCurrentVerseSeed(seed [32]byte) (*types.Transaction, error) {
	return _Updates.Contract.SetCurrentVerseSeed(&_Updates.TransactOpts, seed)
}

// SubmitActionsRoot is a paid mutator transaction binding the contract method 0xe8507051.
//
// Solidity: function submitActionsRoot(bytes32 actionsRoot) returns()
func (_Updates *UpdatesTransactor) SubmitActionsRoot(opts *bind.TransactOpts, actionsRoot [32]byte) (*types.Transaction, error) {
	return _Updates.contract.Transact(opts, "submitActionsRoot", actionsRoot)
}

// SubmitActionsRoot is a paid mutator transaction binding the contract method 0xe8507051.
//
// Solidity: function submitActionsRoot(bytes32 actionsRoot) returns()
func (_Updates *UpdatesSession) SubmitActionsRoot(actionsRoot [32]byte) (*types.Transaction, error) {
	return _Updates.Contract.SubmitActionsRoot(&_Updates.TransactOpts, actionsRoot)
}

// SubmitActionsRoot is a paid mutator transaction binding the contract method 0xe8507051.
//
// Solidity: function submitActionsRoot(bytes32 actionsRoot) returns()
func (_Updates *UpdatesTransactorSession) SubmitActionsRoot(actionsRoot [32]byte) (*types.Transaction, error) {
	return _Updates.Contract.SubmitActionsRoot(&_Updates.TransactOpts, actionsRoot)
}

// UpdateTZ is a paid mutator transaction binding the contract method 0x16b5d047.
//
// Solidity: function updateTZ(bytes32 root) returns()
func (_Updates *UpdatesTransactor) UpdateTZ(opts *bind.TransactOpts, root [32]byte) (*types.Transaction, error) {
	return _Updates.contract.Transact(opts, "updateTZ", root)
}

// UpdateTZ is a paid mutator transaction binding the contract method 0x16b5d047.
//
// Solidity: function updateTZ(bytes32 root) returns()
func (_Updates *UpdatesSession) UpdateTZ(root [32]byte) (*types.Transaction, error) {
	return _Updates.Contract.UpdateTZ(&_Updates.TransactOpts, root)
}

// UpdateTZ is a paid mutator transaction binding the contract method 0x16b5d047.
//
// Solidity: function updateTZ(bytes32 root) returns()
func (_Updates *UpdatesTransactorSession) UpdateTZ(root [32]byte) (*types.Transaction, error) {
	return _Updates.Contract.UpdateTZ(&_Updates.TransactOpts, root)
}

// UpdatesActionsSubmissionIterator is returned from FilterActionsSubmission and is used to iterate over the raw logs and unpacked data for ActionsSubmission events raised by the Updates contract.
type UpdatesActionsSubmissionIterator struct {
	Event *UpdatesActionsSubmission // Event containing the contract specifics and raw log

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
func (it *UpdatesActionsSubmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpdatesActionsSubmission)
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
		it.Event = new(UpdatesActionsSubmission)
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
func (it *UpdatesActionsSubmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpdatesActionsSubmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpdatesActionsSubmission represents a ActionsSubmission event raised by the Updates contract.
type UpdatesActionsSubmission struct {
	TimeZone       uint8
	Day            uint8
	TurnInDay      uint8
	Seed           [32]byte
	SubmissionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterActionsSubmission is a free log retrieval operation binding the contract event 0x3eaf8623b39cd8460e98a71c435b5fea0d21dc62cc5588b0e628edf98b5db685.
//
// Solidity: event ActionsSubmission(uint8 timeZone, uint8 day, uint8 turnInDay, bytes32 seed, uint256 submissionTime)
func (_Updates *UpdatesFilterer) FilterActionsSubmission(opts *bind.FilterOpts) (*UpdatesActionsSubmissionIterator, error) {

	logs, sub, err := _Updates.contract.FilterLogs(opts, "ActionsSubmission")
	if err != nil {
		return nil, err
	}
	return &UpdatesActionsSubmissionIterator{contract: _Updates.contract, event: "ActionsSubmission", logs: logs, sub: sub}, nil
}

// WatchActionsSubmission is a free log subscription operation binding the contract event 0x3eaf8623b39cd8460e98a71c435b5fea0d21dc62cc5588b0e628edf98b5db685.
//
// Solidity: event ActionsSubmission(uint8 timeZone, uint8 day, uint8 turnInDay, bytes32 seed, uint256 submissionTime)
func (_Updates *UpdatesFilterer) WatchActionsSubmission(opts *bind.WatchOpts, sink chan<- *UpdatesActionsSubmission) (event.Subscription, error) {

	logs, sub, err := _Updates.contract.WatchLogs(opts, "ActionsSubmission")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpdatesActionsSubmission)
				if err := _Updates.contract.UnpackLog(event, "ActionsSubmission", log); err != nil {
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

// UpdatesTeamTransferIterator is returned from FilterTeamTransfer and is used to iterate over the raw logs and unpacked data for TeamTransfer events raised by the Updates contract.
type UpdatesTeamTransferIterator struct {
	Event *UpdatesTeamTransfer // Event containing the contract specifics and raw log

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
func (it *UpdatesTeamTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpdatesTeamTransfer)
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
		it.Event = new(UpdatesTeamTransfer)
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
func (it *UpdatesTeamTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpdatesTeamTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpdatesTeamTransfer represents a TeamTransfer event raised by the Updates contract.
type UpdatesTeamTransfer struct {
	TeamId *big.Int
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTeamTransfer is a free log retrieval operation binding the contract event 0x77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf38.
//
// Solidity: event TeamTransfer(uint256 teamId, address to)
func (_Updates *UpdatesFilterer) FilterTeamTransfer(opts *bind.FilterOpts) (*UpdatesTeamTransferIterator, error) {

	logs, sub, err := _Updates.contract.FilterLogs(opts, "TeamTransfer")
	if err != nil {
		return nil, err
	}
	return &UpdatesTeamTransferIterator{contract: _Updates.contract, event: "TeamTransfer", logs: logs, sub: sub}, nil
}

// WatchTeamTransfer is a free log subscription operation binding the contract event 0x77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf38.
//
// Solidity: event TeamTransfer(uint256 teamId, address to)
func (_Updates *UpdatesFilterer) WatchTeamTransfer(opts *bind.WatchOpts, sink chan<- *UpdatesTeamTransfer) (event.Subscription, error) {

	logs, sub, err := _Updates.contract.WatchLogs(opts, "TeamTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpdatesTeamTransfer)
				if err := _Updates.contract.UnpackLog(event, "TeamTransfer", log); err != nil {
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

// UpdatesTimeZoneUpdateIterator is returned from FilterTimeZoneUpdate and is used to iterate over the raw logs and unpacked data for TimeZoneUpdate events raised by the Updates contract.
type UpdatesTimeZoneUpdateIterator struct {
	Event *UpdatesTimeZoneUpdate // Event containing the contract specifics and raw log

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
func (it *UpdatesTimeZoneUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpdatesTimeZoneUpdate)
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
		it.Event = new(UpdatesTimeZoneUpdate)
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
func (it *UpdatesTimeZoneUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpdatesTimeZoneUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpdatesTimeZoneUpdate represents a TimeZoneUpdate event raised by the Updates contract.
type UpdatesTimeZoneUpdate struct {
	TimeZone       uint8
	Root           [32]byte
	SubmissionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterTimeZoneUpdate is a free log retrieval operation binding the contract event 0x21910d8429bd1db0b36330d1ac1df9f94ff1eeae897e11db418e30bc0fb418b9.
//
// Solidity: event TimeZoneUpdate(uint8 timeZone, bytes32 root, uint256 submissionTime)
func (_Updates *UpdatesFilterer) FilterTimeZoneUpdate(opts *bind.FilterOpts) (*UpdatesTimeZoneUpdateIterator, error) {

	logs, sub, err := _Updates.contract.FilterLogs(opts, "TimeZoneUpdate")
	if err != nil {
		return nil, err
	}
	return &UpdatesTimeZoneUpdateIterator{contract: _Updates.contract, event: "TimeZoneUpdate", logs: logs, sub: sub}, nil
}

// WatchTimeZoneUpdate is a free log subscription operation binding the contract event 0x21910d8429bd1db0b36330d1ac1df9f94ff1eeae897e11db418e30bc0fb418b9.
//
// Solidity: event TimeZoneUpdate(uint8 timeZone, bytes32 root, uint256 submissionTime)
func (_Updates *UpdatesFilterer) WatchTimeZoneUpdate(opts *bind.WatchOpts, sink chan<- *UpdatesTimeZoneUpdate) (event.Subscription, error) {

	logs, sub, err := _Updates.contract.WatchLogs(opts, "TimeZoneUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpdatesTimeZoneUpdate)
				if err := _Updates.contract.UnpackLog(event, "TimeZoneUpdate", log); err != nil {
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
