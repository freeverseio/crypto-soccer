// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package states

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

// StatesABI is the input ABI used to generate the binding from.
const StatesABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getSkills\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"}],\"name\":\"isValidPlayerState\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getEndurance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPrevLeagueId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getSpeed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getDefence\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"defence\",\"type\":\"uint256\"},{\"name\":\"speed\",\"type\":\"uint256\"},{\"name\":\"pass\",\"type\":\"uint256\"},{\"name\":\"shoot\",\"type\":\"uint256\"},{\"name\":\"endurance\",\"type\":\"uint256\"},{\"name\":\"monthOfBirthInUnixTime\",\"type\":\"uint256\"},{\"name\":\"playerId\",\"type\":\"uint256\"},{\"name\":\"currentTeamId\",\"type\":\"uint256\"},{\"name\":\"currentShirtNum\",\"type\":\"uint256\"},{\"name\":\"prevLeagueId\",\"type\":\"uint256\"},{\"name\":\"prevTeamPosInLeague\",\"type\":\"uint256\"},{\"name\":\"prevShirtNumInLeague\",\"type\":\"uint256\"},{\"name\":\"lastSaleBlock\",\"type\":\"uint256\"}],\"name\":\"playerStateCreate\",\"outputs\":[{\"name\":\"state\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPass\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPrevTeamPosInLeague\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getShoot\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPrevShirtNumInLeague\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPrevTeamPosInLeague\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getMonthOfBirthInUnixTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"},{\"name\":\"delta\",\"type\":\"uint16\"}],\"name\":\"playerStateEvolve\",\"outputs\":[{\"name\":\"evolvedState\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"},{\"name\":\"currentShirtNum\",\"type\":\"uint256\"}],\"name\":\"setCurrentShirtNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"},{\"name\":\"lastSaleBlock\",\"type\":\"uint256\"}],\"name\":\"setLastSaleBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"},{\"name\":\"teamId\",\"type\":\"uint256\"}],\"name\":\"setCurrentTeamId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getLastSaleBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getSkillsVec\",\"outputs\":[{\"name\":\"skills\",\"type\":\"uint16[5]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getCurrentTeamId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getCurrentShirtNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPlayerId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPrevLeagueId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"teamStateCreate\",\"outputs\":[{\"name\":\"state\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"},{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"teamStateAppend\",\"outputs\":[{\"name\":\"state\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"}],\"name\":\"teamStateSize\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"},{\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"teamStateAt\",\"outputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256[]\"}],\"name\":\"isValidTeamState\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"},{\"name\":\"delta\",\"type\":\"uint8\"}],\"name\":\"teamStateEvolve\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"}],\"name\":\"computeTeamRating\",\"outputs\":[{\"name\":\"rating\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// States is an auto generated Go binding around an Ethereum contract.
type States struct {
	StatesCaller     // Read-only binding to the contract
	StatesTransactor // Write-only binding to the contract
	StatesFilterer   // Log filterer for contract events
}

// StatesCaller is an auto generated read-only Go binding around an Ethereum contract.
type StatesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StatesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StatesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StatesSession struct {
	Contract     *States           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StatesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StatesCallerSession struct {
	Contract *StatesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StatesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StatesTransactorSession struct {
	Contract     *StatesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StatesRaw is an auto generated low-level Go binding around an Ethereum contract.
type StatesRaw struct {
	Contract *States // Generic contract binding to access the raw methods on
}

// StatesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StatesCallerRaw struct {
	Contract *StatesCaller // Generic read-only contract binding to access the raw methods on
}

// StatesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StatesTransactorRaw struct {
	Contract *StatesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStates creates a new instance of States, bound to a specific deployed contract.
func NewStates(address common.Address, backend bind.ContractBackend) (*States, error) {
	contract, err := bindStates(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &States{StatesCaller: StatesCaller{contract: contract}, StatesTransactor: StatesTransactor{contract: contract}, StatesFilterer: StatesFilterer{contract: contract}}, nil
}

// NewStatesCaller creates a new read-only instance of States, bound to a specific deployed contract.
func NewStatesCaller(address common.Address, caller bind.ContractCaller) (*StatesCaller, error) {
	contract, err := bindStates(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StatesCaller{contract: contract}, nil
}

// NewStatesTransactor creates a new write-only instance of States, bound to a specific deployed contract.
func NewStatesTransactor(address common.Address, transactor bind.ContractTransactor) (*StatesTransactor, error) {
	contract, err := bindStates(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StatesTransactor{contract: contract}, nil
}

// NewStatesFilterer creates a new log filterer instance of States, bound to a specific deployed contract.
func NewStatesFilterer(address common.Address, filterer bind.ContractFilterer) (*StatesFilterer, error) {
	contract, err := bindStates(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StatesFilterer{contract: contract}, nil
}

// bindStates binds a generic wrapper to an already deployed contract.
func bindStates(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StatesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_States *StatesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _States.Contract.StatesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_States *StatesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _States.Contract.StatesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_States *StatesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _States.Contract.StatesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_States *StatesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _States.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_States *StatesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _States.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_States *StatesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _States.Contract.contract.Transact(opts, method, params...)
}

// ComputeTeamRating is a free data retrieval call binding the contract method 0xbd64d4fa.
//
// Solidity: function computeTeamRating(uint256[] teamState) constant returns(uint256 rating)
func (_States *StatesCaller) ComputeTeamRating(opts *bind.CallOpts, teamState []*big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "computeTeamRating", teamState)
	return *ret0, err
}

// ComputeTeamRating is a free data retrieval call binding the contract method 0xbd64d4fa.
//
// Solidity: function computeTeamRating(uint256[] teamState) constant returns(uint256 rating)
func (_States *StatesSession) ComputeTeamRating(teamState []*big.Int) (*big.Int, error) {
	return _States.Contract.ComputeTeamRating(&_States.CallOpts, teamState)
}

// ComputeTeamRating is a free data retrieval call binding the contract method 0xbd64d4fa.
//
// Solidity: function computeTeamRating(uint256[] teamState) constant returns(uint256 rating)
func (_States *StatesCallerSession) ComputeTeamRating(teamState []*big.Int) (*big.Int, error) {
	return _States.Contract.ComputeTeamRating(&_States.CallOpts, teamState)
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetCurrentShirtNum(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getCurrentShirtNum", playerState)
	return *ret0, err
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetCurrentShirtNum(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetCurrentShirtNum(&_States.CallOpts, playerState)
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetCurrentShirtNum(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetCurrentShirtNum(&_States.CallOpts, playerState)
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetCurrentTeamId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getCurrentTeamId", playerState)
	return *ret0, err
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetCurrentTeamId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetCurrentTeamId(&_States.CallOpts, playerState)
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetCurrentTeamId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetCurrentTeamId(&_States.CallOpts, playerState)
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetDefence(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getDefence", playerState)
	return *ret0, err
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetDefence(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetDefence(&_States.CallOpts, playerState)
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetDefence(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetDefence(&_States.CallOpts, playerState)
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetEndurance(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getEndurance", playerState)
	return *ret0, err
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetEndurance(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetEndurance(&_States.CallOpts, playerState)
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetEndurance(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetEndurance(&_States.CallOpts, playerState)
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetLastSaleBlock(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getLastSaleBlock", playerState)
	return *ret0, err
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetLastSaleBlock(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetLastSaleBlock(&_States.CallOpts, playerState)
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetLastSaleBlock(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetLastSaleBlock(&_States.CallOpts, playerState)
}

// GetMonthOfBirthInUnixTime is a free data retrieval call binding the contract method 0x85053566.
//
// Solidity: function getMonthOfBirthInUnixTime(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetMonthOfBirthInUnixTime(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getMonthOfBirthInUnixTime", playerState)
	return *ret0, err
}

// GetMonthOfBirthInUnixTime is a free data retrieval call binding the contract method 0x85053566.
//
// Solidity: function getMonthOfBirthInUnixTime(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetMonthOfBirthInUnixTime(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetMonthOfBirthInUnixTime(&_States.CallOpts, playerState)
}

// GetMonthOfBirthInUnixTime is a free data retrieval call binding the contract method 0x85053566.
//
// Solidity: function getMonthOfBirthInUnixTime(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetMonthOfBirthInUnixTime(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetMonthOfBirthInUnixTime(&_States.CallOpts, playerState)
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPass(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPass", playerState)
	return *ret0, err
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPass(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPass(&_States.CallOpts, playerState)
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPass(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPass(&_States.CallOpts, playerState)
}

// GetPlayerId is a free data retrieval call binding the contract method 0xf4385912.
//
// Solidity: function getPlayerId(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPlayerId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPlayerId", playerState)
	return *ret0, err
}

// GetPlayerId is a free data retrieval call binding the contract method 0xf4385912.
//
// Solidity: function getPlayerId(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPlayerId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPlayerId(&_States.CallOpts, playerState)
}

// GetPlayerId is a free data retrieval call binding the contract method 0xf4385912.
//
// Solidity: function getPlayerId(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPlayerId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPlayerId(&_States.CallOpts, playerState)
}

// GetPrevLeagueId is a free data retrieval call binding the contract method 0xf8bd3e6e.
//
// Solidity: function getPrevLeagueId(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPrevLeagueId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPrevLeagueId", playerState)
	return *ret0, err
}

// GetPrevLeagueId is a free data retrieval call binding the contract method 0xf8bd3e6e.
//
// Solidity: function getPrevLeagueId(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPrevLeagueId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevLeagueId(&_States.CallOpts, playerState)
}

// GetPrevLeagueId is a free data retrieval call binding the contract method 0xf8bd3e6e.
//
// Solidity: function getPrevLeagueId(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPrevLeagueId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevLeagueId(&_States.CallOpts, playerState)
}

// GetPrevShirtNumInLeague is a free data retrieval call binding the contract method 0x666d0070.
//
// Solidity: function getPrevShirtNumInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPrevShirtNumInLeague(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPrevShirtNumInLeague", playerState)
	return *ret0, err
}

// GetPrevShirtNumInLeague is a free data retrieval call binding the contract method 0x666d0070.
//
// Solidity: function getPrevShirtNumInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPrevShirtNumInLeague(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevShirtNumInLeague(&_States.CallOpts, playerState)
}

// GetPrevShirtNumInLeague is a free data retrieval call binding the contract method 0x666d0070.
//
// Solidity: function getPrevShirtNumInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPrevShirtNumInLeague(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevShirtNumInLeague(&_States.CallOpts, playerState)
}

// GetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x58a7a46a.
//
// Solidity: function getPrevTeamPosInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPrevTeamPosInLeague(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPrevTeamPosInLeague", playerState)
	return *ret0, err
}

// GetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x58a7a46a.
//
// Solidity: function getPrevTeamPosInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPrevTeamPosInLeague(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevTeamPosInLeague(&_States.CallOpts, playerState)
}

// GetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x58a7a46a.
//
// Solidity: function getPrevTeamPosInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPrevTeamPosInLeague(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevTeamPosInLeague(&_States.CallOpts, playerState)
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetShoot(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getShoot", playerState)
	return *ret0, err
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetShoot(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetShoot(&_States.CallOpts, playerState)
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetShoot(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetShoot(&_States.CallOpts, playerState)
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetSkills(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getSkills", playerState)
	return *ret0, err
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetSkills(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetSkills(&_States.CallOpts, playerState)
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetSkills(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetSkills(&_States.CallOpts, playerState)
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 playerState) constant returns(uint16[5] skills)
func (_States *StatesCaller) GetSkillsVec(opts *bind.CallOpts, playerState *big.Int) ([5]uint16, error) {
	var (
		ret0 = new([5]uint16)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getSkillsVec", playerState)
	return *ret0, err
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 playerState) constant returns(uint16[5] skills)
func (_States *StatesSession) GetSkillsVec(playerState *big.Int) ([5]uint16, error) {
	return _States.Contract.GetSkillsVec(&_States.CallOpts, playerState)
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 playerState) constant returns(uint16[5] skills)
func (_States *StatesCallerSession) GetSkillsVec(playerState *big.Int) ([5]uint16, error) {
	return _States.Contract.GetSkillsVec(&_States.CallOpts, playerState)
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetSpeed(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getSpeed", playerState)
	return *ret0, err
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetSpeed(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetSpeed(&_States.CallOpts, playerState)
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetSpeed(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetSpeed(&_States.CallOpts, playerState)
}

// IsValidPlayerState is a free data retrieval call binding the contract method 0x19a4860c.
//
// Solidity: function isValidPlayerState(uint256 state) constant returns(bool)
func (_States *StatesCaller) IsValidPlayerState(opts *bind.CallOpts, state *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "isValidPlayerState", state)
	return *ret0, err
}

// IsValidPlayerState is a free data retrieval call binding the contract method 0x19a4860c.
//
// Solidity: function isValidPlayerState(uint256 state) constant returns(bool)
func (_States *StatesSession) IsValidPlayerState(state *big.Int) (bool, error) {
	return _States.Contract.IsValidPlayerState(&_States.CallOpts, state)
}

// IsValidPlayerState is a free data retrieval call binding the contract method 0x19a4860c.
//
// Solidity: function isValidPlayerState(uint256 state) constant returns(bool)
func (_States *StatesCallerSession) IsValidPlayerState(state *big.Int) (bool, error) {
	return _States.Contract.IsValidPlayerState(&_States.CallOpts, state)
}

// IsValidTeamState is a free data retrieval call binding the contract method 0xaf26723a.
//
// Solidity: function isValidTeamState(uint256[] state) constant returns(bool)
func (_States *StatesCaller) IsValidTeamState(opts *bind.CallOpts, state []*big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "isValidTeamState", state)
	return *ret0, err
}

// IsValidTeamState is a free data retrieval call binding the contract method 0xaf26723a.
//
// Solidity: function isValidTeamState(uint256[] state) constant returns(bool)
func (_States *StatesSession) IsValidTeamState(state []*big.Int) (bool, error) {
	return _States.Contract.IsValidTeamState(&_States.CallOpts, state)
}

// IsValidTeamState is a free data retrieval call binding the contract method 0xaf26723a.
//
// Solidity: function isValidTeamState(uint256[] state) constant returns(bool)
func (_States *StatesCallerSession) IsValidTeamState(state []*big.Int) (bool, error) {
	return _States.Contract.IsValidTeamState(&_States.CallOpts, state)
}

// PlayerStateCreate is a free data retrieval call binding the contract method 0x530f63d6.
//
// Solidity: function playerStateCreate(uint256 defence, uint256 speed, uint256 pass, uint256 shoot, uint256 endurance, uint256 monthOfBirthInUnixTime, uint256 playerId, uint256 currentTeamId, uint256 currentShirtNum, uint256 prevLeagueId, uint256 prevTeamPosInLeague, uint256 prevShirtNumInLeague, uint256 lastSaleBlock) constant returns(uint256 state)
func (_States *StatesCaller) PlayerStateCreate(opts *bind.CallOpts, defence *big.Int, speed *big.Int, pass *big.Int, shoot *big.Int, endurance *big.Int, monthOfBirthInUnixTime *big.Int, playerId *big.Int, currentTeamId *big.Int, currentShirtNum *big.Int, prevLeagueId *big.Int, prevTeamPosInLeague *big.Int, prevShirtNumInLeague *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "playerStateCreate", defence, speed, pass, shoot, endurance, monthOfBirthInUnixTime, playerId, currentTeamId, currentShirtNum, prevLeagueId, prevTeamPosInLeague, prevShirtNumInLeague, lastSaleBlock)
	return *ret0, err
}

// PlayerStateCreate is a free data retrieval call binding the contract method 0x530f63d6.
//
// Solidity: function playerStateCreate(uint256 defence, uint256 speed, uint256 pass, uint256 shoot, uint256 endurance, uint256 monthOfBirthInUnixTime, uint256 playerId, uint256 currentTeamId, uint256 currentShirtNum, uint256 prevLeagueId, uint256 prevTeamPosInLeague, uint256 prevShirtNumInLeague, uint256 lastSaleBlock) constant returns(uint256 state)
func (_States *StatesSession) PlayerStateCreate(defence *big.Int, speed *big.Int, pass *big.Int, shoot *big.Int, endurance *big.Int, monthOfBirthInUnixTime *big.Int, playerId *big.Int, currentTeamId *big.Int, currentShirtNum *big.Int, prevLeagueId *big.Int, prevTeamPosInLeague *big.Int, prevShirtNumInLeague *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _States.Contract.PlayerStateCreate(&_States.CallOpts, defence, speed, pass, shoot, endurance, monthOfBirthInUnixTime, playerId, currentTeamId, currentShirtNum, prevLeagueId, prevTeamPosInLeague, prevShirtNumInLeague, lastSaleBlock)
}

// PlayerStateCreate is a free data retrieval call binding the contract method 0x530f63d6.
//
// Solidity: function playerStateCreate(uint256 defence, uint256 speed, uint256 pass, uint256 shoot, uint256 endurance, uint256 monthOfBirthInUnixTime, uint256 playerId, uint256 currentTeamId, uint256 currentShirtNum, uint256 prevLeagueId, uint256 prevTeamPosInLeague, uint256 prevShirtNumInLeague, uint256 lastSaleBlock) constant returns(uint256 state)
func (_States *StatesCallerSession) PlayerStateCreate(defence *big.Int, speed *big.Int, pass *big.Int, shoot *big.Int, endurance *big.Int, monthOfBirthInUnixTime *big.Int, playerId *big.Int, currentTeamId *big.Int, currentShirtNum *big.Int, prevLeagueId *big.Int, prevTeamPosInLeague *big.Int, prevShirtNumInLeague *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _States.Contract.PlayerStateCreate(&_States.CallOpts, defence, speed, pass, shoot, endurance, monthOfBirthInUnixTime, playerId, currentTeamId, currentShirtNum, prevLeagueId, prevTeamPosInLeague, prevShirtNumInLeague, lastSaleBlock)
}

// PlayerStateEvolve is a free data retrieval call binding the contract method 0x8d216b52.
//
// Solidity: function playerStateEvolve(uint256 playerState, uint16 delta) constant returns(uint256 evolvedState)
func (_States *StatesCaller) PlayerStateEvolve(opts *bind.CallOpts, playerState *big.Int, delta uint16) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "playerStateEvolve", playerState, delta)
	return *ret0, err
}

// PlayerStateEvolve is a free data retrieval call binding the contract method 0x8d216b52.
//
// Solidity: function playerStateEvolve(uint256 playerState, uint16 delta) constant returns(uint256 evolvedState)
func (_States *StatesSession) PlayerStateEvolve(playerState *big.Int, delta uint16) (*big.Int, error) {
	return _States.Contract.PlayerStateEvolve(&_States.CallOpts, playerState, delta)
}

// PlayerStateEvolve is a free data retrieval call binding the contract method 0x8d216b52.
//
// Solidity: function playerStateEvolve(uint256 playerState, uint16 delta) constant returns(uint256 evolvedState)
func (_States *StatesCallerSession) PlayerStateEvolve(playerState *big.Int, delta uint16) (*big.Int, error) {
	return _States.Contract.PlayerStateEvolve(&_States.CallOpts, playerState, delta)
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0xa95e858b.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint256 currentShirtNum) constant returns(uint256)
func (_States *StatesCaller) SetCurrentShirtNum(opts *bind.CallOpts, state *big.Int, currentShirtNum *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setCurrentShirtNum", state, currentShirtNum)
	return *ret0, err
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0xa95e858b.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint256 currentShirtNum) constant returns(uint256)
func (_States *StatesSession) SetCurrentShirtNum(state *big.Int, currentShirtNum *big.Int) (*big.Int, error) {
	return _States.Contract.SetCurrentShirtNum(&_States.CallOpts, state, currentShirtNum)
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0xa95e858b.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint256 currentShirtNum) constant returns(uint256)
func (_States *StatesCallerSession) SetCurrentShirtNum(state *big.Int, currentShirtNum *big.Int) (*big.Int, error) {
	return _States.Contract.SetCurrentShirtNum(&_States.CallOpts, state, currentShirtNum)
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_States *StatesCaller) SetCurrentTeamId(opts *bind.CallOpts, playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setCurrentTeamId", playerState, teamId)
	return *ret0, err
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_States *StatesSession) SetCurrentTeamId(playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	return _States.Contract.SetCurrentTeamId(&_States.CallOpts, playerState, teamId)
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_States *StatesCallerSession) SetCurrentTeamId(playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	return _States.Contract.SetCurrentTeamId(&_States.CallOpts, playerState, teamId)
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_States *StatesCaller) SetLastSaleBlock(opts *bind.CallOpts, state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setLastSaleBlock", state, lastSaleBlock)
	return *ret0, err
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_States *StatesSession) SetLastSaleBlock(state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _States.Contract.SetLastSaleBlock(&_States.CallOpts, state, lastSaleBlock)
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_States *StatesCallerSession) SetLastSaleBlock(state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _States.Contract.SetLastSaleBlock(&_States.CallOpts, state, lastSaleBlock)
}

// SetPrevLeagueId is a free data retrieval call binding the contract method 0x47f3d716.
//
// Solidity: function setPrevLeagueId(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesCaller) SetPrevLeagueId(opts *bind.CallOpts, state *big.Int, value *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setPrevLeagueId", state, value)
	return *ret0, err
}

// SetPrevLeagueId is a free data retrieval call binding the contract method 0x47f3d716.
//
// Solidity: function setPrevLeagueId(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesSession) SetPrevLeagueId(state *big.Int, value *big.Int) (*big.Int, error) {
	return _States.Contract.SetPrevLeagueId(&_States.CallOpts, state, value)
}

// SetPrevLeagueId is a free data retrieval call binding the contract method 0x47f3d716.
//
// Solidity: function setPrevLeagueId(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesCallerSession) SetPrevLeagueId(state *big.Int, value *big.Int) (*big.Int, error) {
	return _States.Contract.SetPrevLeagueId(&_States.CallOpts, state, value)
}

// SetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x7ee0aebc.
//
// Solidity: function setPrevTeamPosInLeague(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesCaller) SetPrevTeamPosInLeague(opts *bind.CallOpts, state *big.Int, value *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setPrevTeamPosInLeague", state, value)
	return *ret0, err
}

// SetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x7ee0aebc.
//
// Solidity: function setPrevTeamPosInLeague(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesSession) SetPrevTeamPosInLeague(state *big.Int, value *big.Int) (*big.Int, error) {
	return _States.Contract.SetPrevTeamPosInLeague(&_States.CallOpts, state, value)
}

// SetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x7ee0aebc.
//
// Solidity: function setPrevTeamPosInLeague(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesCallerSession) SetPrevTeamPosInLeague(state *big.Int, value *big.Int) (*big.Int, error) {
	return _States.Contract.SetPrevTeamPosInLeague(&_States.CallOpts, state, value)
}

// TeamStateAppend is a free data retrieval call binding the contract method 0xe06da0bf.
//
// Solidity: function teamStateAppend(uint256[] teamState, uint256 playerState) constant returns(uint256[] state)
func (_States *StatesCaller) TeamStateAppend(opts *bind.CallOpts, teamState []*big.Int, playerState *big.Int) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateAppend", teamState, playerState)
	return *ret0, err
}

// TeamStateAppend is a free data retrieval call binding the contract method 0xe06da0bf.
//
// Solidity: function teamStateAppend(uint256[] teamState, uint256 playerState) constant returns(uint256[] state)
func (_States *StatesSession) TeamStateAppend(teamState []*big.Int, playerState *big.Int) ([]*big.Int, error) {
	return _States.Contract.TeamStateAppend(&_States.CallOpts, teamState, playerState)
}

// TeamStateAppend is a free data retrieval call binding the contract method 0xe06da0bf.
//
// Solidity: function teamStateAppend(uint256[] teamState, uint256 playerState) constant returns(uint256[] state)
func (_States *StatesCallerSession) TeamStateAppend(teamState []*big.Int, playerState *big.Int) ([]*big.Int, error) {
	return _States.Contract.TeamStateAppend(&_States.CallOpts, teamState, playerState)
}

// TeamStateAt is a free data retrieval call binding the contract method 0x44328d1a.
//
// Solidity: function teamStateAt(uint256[] teamState, uint256 idx) constant returns(uint256 playerState)
func (_States *StatesCaller) TeamStateAt(opts *bind.CallOpts, teamState []*big.Int, idx *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateAt", teamState, idx)
	return *ret0, err
}

// TeamStateAt is a free data retrieval call binding the contract method 0x44328d1a.
//
// Solidity: function teamStateAt(uint256[] teamState, uint256 idx) constant returns(uint256 playerState)
func (_States *StatesSession) TeamStateAt(teamState []*big.Int, idx *big.Int) (*big.Int, error) {
	return _States.Contract.TeamStateAt(&_States.CallOpts, teamState, idx)
}

// TeamStateAt is a free data retrieval call binding the contract method 0x44328d1a.
//
// Solidity: function teamStateAt(uint256[] teamState, uint256 idx) constant returns(uint256 playerState)
func (_States *StatesCallerSession) TeamStateAt(teamState []*big.Int, idx *big.Int) (*big.Int, error) {
	return _States.Contract.TeamStateAt(&_States.CallOpts, teamState, idx)
}

// TeamStateCreate is a free data retrieval call binding the contract method 0x74833f73.
//
// Solidity: function teamStateCreate() constant returns(uint256[] state)
func (_States *StatesCaller) TeamStateCreate(opts *bind.CallOpts) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateCreate")
	return *ret0, err
}

// TeamStateCreate is a free data retrieval call binding the contract method 0x74833f73.
//
// Solidity: function teamStateCreate() constant returns(uint256[] state)
func (_States *StatesSession) TeamStateCreate() ([]*big.Int, error) {
	return _States.Contract.TeamStateCreate(&_States.CallOpts)
}

// TeamStateCreate is a free data retrieval call binding the contract method 0x74833f73.
//
// Solidity: function teamStateCreate() constant returns(uint256[] state)
func (_States *StatesCallerSession) TeamStateCreate() ([]*big.Int, error) {
	return _States.Contract.TeamStateCreate(&_States.CallOpts)
}

// TeamStateEvolve is a free data retrieval call binding the contract method 0x5c010782.
//
// Solidity: function teamStateEvolve(uint256[] teamState, uint8 delta) constant returns(uint256[])
func (_States *StatesCaller) TeamStateEvolve(opts *bind.CallOpts, teamState []*big.Int, delta uint8) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateEvolve", teamState, delta)
	return *ret0, err
}

// TeamStateEvolve is a free data retrieval call binding the contract method 0x5c010782.
//
// Solidity: function teamStateEvolve(uint256[] teamState, uint8 delta) constant returns(uint256[])
func (_States *StatesSession) TeamStateEvolve(teamState []*big.Int, delta uint8) ([]*big.Int, error) {
	return _States.Contract.TeamStateEvolve(&_States.CallOpts, teamState, delta)
}

// TeamStateEvolve is a free data retrieval call binding the contract method 0x5c010782.
//
// Solidity: function teamStateEvolve(uint256[] teamState, uint8 delta) constant returns(uint256[])
func (_States *StatesCallerSession) TeamStateEvolve(teamState []*big.Int, delta uint8) ([]*big.Int, error) {
	return _States.Contract.TeamStateEvolve(&_States.CallOpts, teamState, delta)
}

// TeamStateSize is a free data retrieval call binding the contract method 0x4444a283.
//
// Solidity: function teamStateSize(uint256[] teamState) constant returns(uint256 count)
func (_States *StatesCaller) TeamStateSize(opts *bind.CallOpts, teamState []*big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateSize", teamState)
	return *ret0, err
}

// TeamStateSize is a free data retrieval call binding the contract method 0x4444a283.
//
// Solidity: function teamStateSize(uint256[] teamState) constant returns(uint256 count)
func (_States *StatesSession) TeamStateSize(teamState []*big.Int) (*big.Int, error) {
	return _States.Contract.TeamStateSize(&_States.CallOpts, teamState)
}

// TeamStateSize is a free data retrieval call binding the contract method 0x4444a283.
//
// Solidity: function teamStateSize(uint256[] teamState) constant returns(uint256 count)
func (_States *StatesCallerSession) TeamStateSize(teamState []*big.Int) (*big.Int, error) {
	return _States.Contract.TeamStateSize(&_States.CallOpts, teamState)
}
