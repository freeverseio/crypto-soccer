// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package leagues

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

// LeaguesABI is the input ABI used to generate the binding from.
const LeaguesABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getScores\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamIds\",\"type\":\"uint256[]\"},{\"name\":\"tactics\",\"type\":\"uint8[3][]\"},{\"name\":\"blocks\",\"type\":\"uint256[]\"}],\"name\":\"computeUsersAlongDataHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getNTeams\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"hasFinished\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"resetUpdater\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpdateBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"leaguesCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"scores\",\"type\":\"uint16[]\"},{\"name\":\"score\",\"type\":\"uint16\"}],\"name\":\"scoresAppend\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isVerified\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tactics\",\"type\":\"uint256[3][]\"}],\"name\":\"hashTactics\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"day\",\"type\":\"uint256\"}],\"name\":\"scoresGetDay\",\"outputs\":[{\"name\":\"dayScores\",\"type\":\"uint16[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"hasStarted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getTactics\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getDayStateHashes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"home\",\"type\":\"uint8\"},{\"name\":\"visitor\",\"type\":\"uint8\"}],\"name\":\"encodeScore\",\"outputs\":[{\"name\":\"score\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"leagueId\",\"type\":\"uint256\"},{\"name\":\"leagueDay\",\"type\":\"uint256\"},{\"name\":\"initLeagueState\",\"type\":\"uint256[]\"},{\"name\":\"tactics\",\"type\":\"uint8[3][]\"}],\"name\":\"computeDay\",\"outputs\":[{\"name\":\"scores\",\"type\":\"uint16[]\"},{\"name\":\"finalLeagueState\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"initBlock\",\"type\":\"uint256\"},{\"name\":\"step\",\"type\":\"uint256\"},{\"name\":\"teamIds\",\"type\":\"uint256[]\"},{\"name\":\"tactics\",\"type\":\"uint8[3][]\"}],\"name\":\"create\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getChallengePeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getStep\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"usersInitDataTeamIds\",\"type\":\"uint256[]\"},{\"name\":\"usersInitDataTactics\",\"type\":\"uint8[3][]\"},{\"name\":\"usersAlongDataTeamIds\",\"type\":\"uint256[]\"},{\"name\":\"usersAlongDataTactics\",\"type\":\"uint8[3][]\"},{\"name\":\"usersAlongDataBlocks\",\"type\":\"uint256[]\"},{\"name\":\"leagueDay\",\"type\":\"uint256\"},{\"name\":\"prevMatchdayStates\",\"type\":\"uint256[]\"}],\"name\":\"challengeMatchdayStates\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getMatchPerDay\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"matchday\",\"type\":\"uint256\"},{\"name\":\"matchIdx\",\"type\":\"uint256\"}],\"name\":\"getTeamsInMatch\",\"outputs\":[{\"name\":\"homeIdx\",\"type\":\"uint256\"},{\"name\":\"visitorIdx\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getEndBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"initStateHash\",\"type\":\"bytes32\"},{\"name\":\"dayStateHashes\",\"type\":\"bytes32[]\"},{\"name\":\"scores\",\"type\":\"uint16[]\"},{\"name\":\"isLie\",\"type\":\"bool\"}],\"name\":\"updateLeague\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getEngineContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUsersAlongDataHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"target\",\"type\":\"uint16[]\"},{\"name\":\"scores\",\"type\":\"uint16[]\"}],\"name\":\"scoresConcat\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUpdater\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"countLeagueDays\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"teamIds\",\"type\":\"uint256[]\"},{\"name\":\"tactics\",\"type\":\"uint8[3][]\"},{\"name\":\"dataToChallengeInitStates\",\"type\":\"uint256[]\"}],\"name\":\"challengeInitStates\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"teamIds\",\"type\":\"uint256[]\"},{\"name\":\"tactics\",\"type\":\"uint8[3][]\"},{\"name\":\"blocks\",\"type\":\"uint256[]\"}],\"name\":\"updateUsersAlongDataHash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getInitStateHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamIds\",\"type\":\"uint256[]\"},{\"name\":\"tactics\",\"type\":\"uint8[3][]\"}],\"name\":\"hashUsersInitData\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"day\",\"type\":\"uint256\"}],\"name\":\"getMatchDayBlockHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getUsersInitDataHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"stakersContract\",\"type\":\"address\"}],\"name\":\"setStakersContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256[]\"}],\"name\":\"hashInitState\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getLastChallengeBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256[]\"}],\"name\":\"hashDayState\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isUpdated\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getIsLie\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"scoresCreate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getInitBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"teamIds\",\"type\":\"uint256[]\"},{\"name\":\"tactics\",\"type\":\"uint8[3][]\"},{\"name\":\"dataToChallengeInitStates\",\"type\":\"uint256[]\"}],\"name\":\"getInitPlayerStates\",\"outputs\":[{\"name\":\"state\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getTeams\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"score\",\"type\":\"uint16\"}],\"name\":\"decodeScore\",\"outputs\":[{\"name\":\"home\",\"type\":\"uint8\"},{\"name\":\"visitor\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"engine\",\"type\":\"address\"},{\"name\":\"state\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"LeagueCreated\",\"type\":\"event\"}]"

// Leagues is an auto generated Go binding around an Ethereum contract.
type Leagues struct {
	LeaguesCaller     // Read-only binding to the contract
	LeaguesTransactor // Write-only binding to the contract
	LeaguesFilterer   // Log filterer for contract events
}

// LeaguesCaller is an auto generated read-only Go binding around an Ethereum contract.
type LeaguesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LeaguesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LeaguesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LeaguesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LeaguesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LeaguesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LeaguesSession struct {
	Contract     *Leagues          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LeaguesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LeaguesCallerSession struct {
	Contract *LeaguesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// LeaguesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LeaguesTransactorSession struct {
	Contract     *LeaguesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// LeaguesRaw is an auto generated low-level Go binding around an Ethereum contract.
type LeaguesRaw struct {
	Contract *Leagues // Generic contract binding to access the raw methods on
}

// LeaguesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LeaguesCallerRaw struct {
	Contract *LeaguesCaller // Generic read-only contract binding to access the raw methods on
}

// LeaguesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LeaguesTransactorRaw struct {
	Contract *LeaguesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLeagues creates a new instance of Leagues, bound to a specific deployed contract.
func NewLeagues(address common.Address, backend bind.ContractBackend) (*Leagues, error) {
	contract, err := bindLeagues(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Leagues{LeaguesCaller: LeaguesCaller{contract: contract}, LeaguesTransactor: LeaguesTransactor{contract: contract}, LeaguesFilterer: LeaguesFilterer{contract: contract}}, nil
}

// NewLeaguesCaller creates a new read-only instance of Leagues, bound to a specific deployed contract.
func NewLeaguesCaller(address common.Address, caller bind.ContractCaller) (*LeaguesCaller, error) {
	contract, err := bindLeagues(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LeaguesCaller{contract: contract}, nil
}

// NewLeaguesTransactor creates a new write-only instance of Leagues, bound to a specific deployed contract.
func NewLeaguesTransactor(address common.Address, transactor bind.ContractTransactor) (*LeaguesTransactor, error) {
	contract, err := bindLeagues(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LeaguesTransactor{contract: contract}, nil
}

// NewLeaguesFilterer creates a new log filterer instance of Leagues, bound to a specific deployed contract.
func NewLeaguesFilterer(address common.Address, filterer bind.ContractFilterer) (*LeaguesFilterer, error) {
	contract, err := bindLeagues(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LeaguesFilterer{contract: contract}, nil
}

// bindLeagues binds a generic wrapper to an already deployed contract.
func bindLeagues(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LeaguesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Leagues *LeaguesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Leagues.Contract.LeaguesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Leagues *LeaguesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Leagues.Contract.LeaguesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Leagues *LeaguesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Leagues.Contract.LeaguesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Leagues *LeaguesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Leagues.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Leagues *LeaguesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Leagues.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Leagues *LeaguesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Leagues.Contract.contract.Transact(opts, method, params...)
}

// ComputeDay is a free data retrieval call binding the contract method 0x77b7d5c0.
//
// Solidity: function computeDay(uint256 leagueId, uint256 leagueDay, uint256[] initLeagueState, uint8[3][] tactics) constant returns(uint16[] scores, uint256[] finalLeagueState)
func (_Leagues *LeaguesCaller) ComputeDay(opts *bind.CallOpts, leagueId *big.Int, leagueDay *big.Int, initLeagueState []*big.Int, tactics [][3]uint8) (struct {
	Scores           []uint16
	FinalLeagueState []*big.Int
}, error) {
	ret := new(struct {
		Scores           []uint16
		FinalLeagueState []*big.Int
	})
	out := ret
	err := _Leagues.contract.Call(opts, out, "computeDay", leagueId, leagueDay, initLeagueState, tactics)
	return *ret, err
}

// ComputeDay is a free data retrieval call binding the contract method 0x77b7d5c0.
//
// Solidity: function computeDay(uint256 leagueId, uint256 leagueDay, uint256[] initLeagueState, uint8[3][] tactics) constant returns(uint16[] scores, uint256[] finalLeagueState)
func (_Leagues *LeaguesSession) ComputeDay(leagueId *big.Int, leagueDay *big.Int, initLeagueState []*big.Int, tactics [][3]uint8) (struct {
	Scores           []uint16
	FinalLeagueState []*big.Int
}, error) {
	return _Leagues.Contract.ComputeDay(&_Leagues.CallOpts, leagueId, leagueDay, initLeagueState, tactics)
}

// ComputeDay is a free data retrieval call binding the contract method 0x77b7d5c0.
//
// Solidity: function computeDay(uint256 leagueId, uint256 leagueDay, uint256[] initLeagueState, uint8[3][] tactics) constant returns(uint16[] scores, uint256[] finalLeagueState)
func (_Leagues *LeaguesCallerSession) ComputeDay(leagueId *big.Int, leagueDay *big.Int, initLeagueState []*big.Int, tactics [][3]uint8) (struct {
	Scores           []uint16
	FinalLeagueState []*big.Int
}, error) {
	return _Leagues.Contract.ComputeDay(&_Leagues.CallOpts, leagueId, leagueDay, initLeagueState, tactics)
}

// ComputeUsersAlongDataHash is a free data retrieval call binding the contract method 0x06c6bb99.
//
// Solidity: function computeUsersAlongDataHash(uint256[] teamIds, uint8[3][] tactics, uint256[] blocks) constant returns(bytes32)
func (_Leagues *LeaguesCaller) ComputeUsersAlongDataHash(opts *bind.CallOpts, teamIds []*big.Int, tactics [][3]uint8, blocks []*big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "computeUsersAlongDataHash", teamIds, tactics, blocks)
	return *ret0, err
}

// ComputeUsersAlongDataHash is a free data retrieval call binding the contract method 0x06c6bb99.
//
// Solidity: function computeUsersAlongDataHash(uint256[] teamIds, uint8[3][] tactics, uint256[] blocks) constant returns(bytes32)
func (_Leagues *LeaguesSession) ComputeUsersAlongDataHash(teamIds []*big.Int, tactics [][3]uint8, blocks []*big.Int) ([32]byte, error) {
	return _Leagues.Contract.ComputeUsersAlongDataHash(&_Leagues.CallOpts, teamIds, tactics, blocks)
}

// ComputeUsersAlongDataHash is a free data retrieval call binding the contract method 0x06c6bb99.
//
// Solidity: function computeUsersAlongDataHash(uint256[] teamIds, uint8[3][] tactics, uint256[] blocks) constant returns(bytes32)
func (_Leagues *LeaguesCallerSession) ComputeUsersAlongDataHash(teamIds []*big.Int, tactics [][3]uint8, blocks []*big.Int) ([32]byte, error) {
	return _Leagues.Contract.ComputeUsersAlongDataHash(&_Leagues.CallOpts, teamIds, tactics, blocks)
}

// CountLeagueDays is a free data retrieval call binding the contract method 0xadac0d7f.
//
// Solidity: function countLeagueDays(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCaller) CountLeagueDays(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "countLeagueDays", id)
	return *ret0, err
}

// CountLeagueDays is a free data retrieval call binding the contract method 0xadac0d7f.
//
// Solidity: function countLeagueDays(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesSession) CountLeagueDays(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.CountLeagueDays(&_Leagues.CallOpts, id)
}

// CountLeagueDays is a free data retrieval call binding the contract method 0xadac0d7f.
//
// Solidity: function countLeagueDays(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) CountLeagueDays(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.CountLeagueDays(&_Leagues.CallOpts, id)
}

// DecodeScore is a free data retrieval call binding the contract method 0xfcb10b94.
//
// Solidity: function decodeScore(uint16 score) constant returns(uint8 home, uint8 visitor)
func (_Leagues *LeaguesCaller) DecodeScore(opts *bind.CallOpts, score uint16) (struct {
	Home    uint8
	Visitor uint8
}, error) {
	ret := new(struct {
		Home    uint8
		Visitor uint8
	})
	out := ret
	err := _Leagues.contract.Call(opts, out, "decodeScore", score)
	return *ret, err
}

// DecodeScore is a free data retrieval call binding the contract method 0xfcb10b94.
//
// Solidity: function decodeScore(uint16 score) constant returns(uint8 home, uint8 visitor)
func (_Leagues *LeaguesSession) DecodeScore(score uint16) (struct {
	Home    uint8
	Visitor uint8
}, error) {
	return _Leagues.Contract.DecodeScore(&_Leagues.CallOpts, score)
}

// DecodeScore is a free data retrieval call binding the contract method 0xfcb10b94.
//
// Solidity: function decodeScore(uint16 score) constant returns(uint8 home, uint8 visitor)
func (_Leagues *LeaguesCallerSession) DecodeScore(score uint16) (struct {
	Home    uint8
	Visitor uint8
}, error) {
	return _Leagues.Contract.DecodeScore(&_Leagues.CallOpts, score)
}

// EncodeScore is a free data retrieval call binding the contract method 0x723adff3.
//
// Solidity: function encodeScore(uint8 home, uint8 visitor) constant returns(uint16 score)
func (_Leagues *LeaguesCaller) EncodeScore(opts *bind.CallOpts, home uint8, visitor uint8) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "encodeScore", home, visitor)
	return *ret0, err
}

// EncodeScore is a free data retrieval call binding the contract method 0x723adff3.
//
// Solidity: function encodeScore(uint8 home, uint8 visitor) constant returns(uint16 score)
func (_Leagues *LeaguesSession) EncodeScore(home uint8, visitor uint8) (uint16, error) {
	return _Leagues.Contract.EncodeScore(&_Leagues.CallOpts, home, visitor)
}

// EncodeScore is a free data retrieval call binding the contract method 0x723adff3.
//
// Solidity: function encodeScore(uint8 home, uint8 visitor) constant returns(uint16 score)
func (_Leagues *LeaguesCallerSession) EncodeScore(home uint8, visitor uint8) (uint16, error) {
	return _Leagues.Contract.EncodeScore(&_Leagues.CallOpts, home, visitor)
}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() constant returns(uint256)
func (_Leagues *LeaguesCaller) GetChallengePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getChallengePeriod")
	return *ret0, err
}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() constant returns(uint256)
func (_Leagues *LeaguesSession) GetChallengePeriod() (*big.Int, error) {
	return _Leagues.Contract.GetChallengePeriod(&_Leagues.CallOpts)
}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetChallengePeriod() (*big.Int, error) {
	return _Leagues.Contract.GetChallengePeriod(&_Leagues.CallOpts)
}

// GetDayStateHashes is a free data retrieval call binding the contract method 0x620d4e8d.
//
// Solidity: function getDayStateHashes(uint256 id) constant returns(bytes32[])
func (_Leagues *LeaguesCaller) GetDayStateHashes(opts *bind.CallOpts, id *big.Int) ([][32]byte, error) {
	var (
		ret0 = new([][32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getDayStateHashes", id)
	return *ret0, err
}

// GetDayStateHashes is a free data retrieval call binding the contract method 0x620d4e8d.
//
// Solidity: function getDayStateHashes(uint256 id) constant returns(bytes32[])
func (_Leagues *LeaguesSession) GetDayStateHashes(id *big.Int) ([][32]byte, error) {
	return _Leagues.Contract.GetDayStateHashes(&_Leagues.CallOpts, id)
}

// GetDayStateHashes is a free data retrieval call binding the contract method 0x620d4e8d.
//
// Solidity: function getDayStateHashes(uint256 id) constant returns(bytes32[])
func (_Leagues *LeaguesCallerSession) GetDayStateHashes(id *big.Int) ([][32]byte, error) {
	return _Leagues.Contract.GetDayStateHashes(&_Leagues.CallOpts, id)
}

// GetEndBlock is a free data retrieval call binding the contract method 0x8c8ab7ad.
//
// Solidity: function getEndBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetEndBlock(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getEndBlock", id)
	return *ret0, err
}

// GetEndBlock is a free data retrieval call binding the contract method 0x8c8ab7ad.
//
// Solidity: function getEndBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesSession) GetEndBlock(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetEndBlock(&_Leagues.CallOpts, id)
}

// GetEndBlock is a free data retrieval call binding the contract method 0x8c8ab7ad.
//
// Solidity: function getEndBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetEndBlock(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetEndBlock(&_Leagues.CallOpts, id)
}

// GetEngineContract is a free data retrieval call binding the contract method 0x9c0536d8.
//
// Solidity: function getEngineContract() constant returns(address)
func (_Leagues *LeaguesCaller) GetEngineContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getEngineContract")
	return *ret0, err
}

// GetEngineContract is a free data retrieval call binding the contract method 0x9c0536d8.
//
// Solidity: function getEngineContract() constant returns(address)
func (_Leagues *LeaguesSession) GetEngineContract() (common.Address, error) {
	return _Leagues.Contract.GetEngineContract(&_Leagues.CallOpts)
}

// GetEngineContract is a free data retrieval call binding the contract method 0x9c0536d8.
//
// Solidity: function getEngineContract() constant returns(address)
func (_Leagues *LeaguesCallerSession) GetEngineContract() (common.Address, error) {
	return _Leagues.Contract.GetEngineContract(&_Leagues.CallOpts)
}

// GetInitBlock is a free data retrieval call binding the contract method 0xe878ffa3.
//
// Solidity: function getInitBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetInitBlock(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getInitBlock", id)
	return *ret0, err
}

// GetInitBlock is a free data retrieval call binding the contract method 0xe878ffa3.
//
// Solidity: function getInitBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesSession) GetInitBlock(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetInitBlock(&_Leagues.CallOpts, id)
}

// GetInitBlock is a free data retrieval call binding the contract method 0xe878ffa3.
//
// Solidity: function getInitBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetInitBlock(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetInitBlock(&_Leagues.CallOpts, id)
}

// GetInitStateHash is a free data retrieval call binding the contract method 0xb4c21295.
//
// Solidity: function getInitStateHash(uint256 id) constant returns(bytes32)
func (_Leagues *LeaguesCaller) GetInitStateHash(opts *bind.CallOpts, id *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getInitStateHash", id)
	return *ret0, err
}

// GetInitStateHash is a free data retrieval call binding the contract method 0xb4c21295.
//
// Solidity: function getInitStateHash(uint256 id) constant returns(bytes32)
func (_Leagues *LeaguesSession) GetInitStateHash(id *big.Int) ([32]byte, error) {
	return _Leagues.Contract.GetInitStateHash(&_Leagues.CallOpts, id)
}

// GetInitStateHash is a free data retrieval call binding the contract method 0xb4c21295.
//
// Solidity: function getInitStateHash(uint256 id) constant returns(bytes32)
func (_Leagues *LeaguesCallerSession) GetInitStateHash(id *big.Int) ([32]byte, error) {
	return _Leagues.Contract.GetInitStateHash(&_Leagues.CallOpts, id)
}

// GetIsLie is a free data retrieval call binding the contract method 0xe36decc3.
//
// Solidity: function getIsLie(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCaller) GetIsLie(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getIsLie", id)
	return *ret0, err
}

// GetIsLie is a free data retrieval call binding the contract method 0xe36decc3.
//
// Solidity: function getIsLie(uint256 id) constant returns(bool)
func (_Leagues *LeaguesSession) GetIsLie(id *big.Int) (bool, error) {
	return _Leagues.Contract.GetIsLie(&_Leagues.CallOpts, id)
}

// GetIsLie is a free data retrieval call binding the contract method 0xe36decc3.
//
// Solidity: function getIsLie(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCallerSession) GetIsLie(id *big.Int) (bool, error) {
	return _Leagues.Contract.GetIsLie(&_Leagues.CallOpts, id)
}

// GetLastChallengeBlock is a free data retrieval call binding the contract method 0xc85aed29.
//
// Solidity: function getLastChallengeBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetLastChallengeBlock(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getLastChallengeBlock", id)
	return *ret0, err
}

// GetLastChallengeBlock is a free data retrieval call binding the contract method 0xc85aed29.
//
// Solidity: function getLastChallengeBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesSession) GetLastChallengeBlock(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetLastChallengeBlock(&_Leagues.CallOpts, id)
}

// GetLastChallengeBlock is a free data retrieval call binding the contract method 0xc85aed29.
//
// Solidity: function getLastChallengeBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetLastChallengeBlock(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetLastChallengeBlock(&_Leagues.CallOpts, id)
}

// GetMatchDayBlockHash is a free data retrieval call binding the contract method 0xbdfe30ab.
//
// Solidity: function getMatchDayBlockHash(uint256 id, uint256 day) constant returns(bytes32)
func (_Leagues *LeaguesCaller) GetMatchDayBlockHash(opts *bind.CallOpts, id *big.Int, day *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getMatchDayBlockHash", id, day)
	return *ret0, err
}

// GetMatchDayBlockHash is a free data retrieval call binding the contract method 0xbdfe30ab.
//
// Solidity: function getMatchDayBlockHash(uint256 id, uint256 day) constant returns(bytes32)
func (_Leagues *LeaguesSession) GetMatchDayBlockHash(id *big.Int, day *big.Int) ([32]byte, error) {
	return _Leagues.Contract.GetMatchDayBlockHash(&_Leagues.CallOpts, id, day)
}

// GetMatchDayBlockHash is a free data retrieval call binding the contract method 0xbdfe30ab.
//
// Solidity: function getMatchDayBlockHash(uint256 id, uint256 day) constant returns(bytes32)
func (_Leagues *LeaguesCallerSession) GetMatchDayBlockHash(id *big.Int, day *big.Int) ([32]byte, error) {
	return _Leagues.Contract.GetMatchDayBlockHash(&_Leagues.CallOpts, id, day)
}

// GetMatchPerDay is a free data retrieval call binding the contract method 0x832e0081.
//
// Solidity: function getMatchPerDay(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetMatchPerDay(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getMatchPerDay", id)
	return *ret0, err
}

// GetMatchPerDay is a free data retrieval call binding the contract method 0x832e0081.
//
// Solidity: function getMatchPerDay(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesSession) GetMatchPerDay(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetMatchPerDay(&_Leagues.CallOpts, id)
}

// GetMatchPerDay is a free data retrieval call binding the contract method 0x832e0081.
//
// Solidity: function getMatchPerDay(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetMatchPerDay(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetMatchPerDay(&_Leagues.CallOpts, id)
}

// GetNTeams is a free data retrieval call binding the contract method 0x07a6d222.
//
// Solidity: function getNTeams(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetNTeams(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getNTeams", id)
	return *ret0, err
}

// GetNTeams is a free data retrieval call binding the contract method 0x07a6d222.
//
// Solidity: function getNTeams(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesSession) GetNTeams(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetNTeams(&_Leagues.CallOpts, id)
}

// GetNTeams is a free data retrieval call binding the contract method 0x07a6d222.
//
// Solidity: function getNTeams(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetNTeams(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetNTeams(&_Leagues.CallOpts, id)
}

// GetScores is a free data retrieval call binding the contract method 0x04527d90.
//
// Solidity: function getScores(uint256 id) constant returns(uint16[])
func (_Leagues *LeaguesCaller) GetScores(opts *bind.CallOpts, id *big.Int) ([]uint16, error) {
	var (
		ret0 = new([]uint16)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getScores", id)
	return *ret0, err
}

// GetScores is a free data retrieval call binding the contract method 0x04527d90.
//
// Solidity: function getScores(uint256 id) constant returns(uint16[])
func (_Leagues *LeaguesSession) GetScores(id *big.Int) ([]uint16, error) {
	return _Leagues.Contract.GetScores(&_Leagues.CallOpts, id)
}

// GetScores is a free data retrieval call binding the contract method 0x04527d90.
//
// Solidity: function getScores(uint256 id) constant returns(uint16[])
func (_Leagues *LeaguesCallerSession) GetScores(id *big.Int) ([]uint16, error) {
	return _Leagues.Contract.GetScores(&_Leagues.CallOpts, id)
}

// GetStep is a free data retrieval call binding the contract method 0x7874888a.
//
// Solidity: function getStep(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetStep(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getStep", id)
	return *ret0, err
}

// GetStep is a free data retrieval call binding the contract method 0x7874888a.
//
// Solidity: function getStep(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesSession) GetStep(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetStep(&_Leagues.CallOpts, id)
}

// GetStep is a free data retrieval call binding the contract method 0x7874888a.
//
// Solidity: function getStep(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetStep(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetStep(&_Leagues.CallOpts, id)
}

// GetTactics is a free data retrieval call binding the contract method 0x5c106b7e.
//
// Solidity: function getTactics(uint256 id) constant returns(uint8[])
func (_Leagues *LeaguesCaller) GetTactics(opts *bind.CallOpts, id *big.Int) ([]uint8, error) {
	var (
		ret0 = new([]uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getTactics", id)
	return *ret0, err
}

// GetTactics is a free data retrieval call binding the contract method 0x5c106b7e.
//
// Solidity: function getTactics(uint256 id) constant returns(uint8[])
func (_Leagues *LeaguesSession) GetTactics(id *big.Int) ([]uint8, error) {
	return _Leagues.Contract.GetTactics(&_Leagues.CallOpts, id)
}

// GetTactics is a free data retrieval call binding the contract method 0x5c106b7e.
//
// Solidity: function getTactics(uint256 id) constant returns(uint8[])
func (_Leagues *LeaguesCallerSession) GetTactics(id *big.Int) ([]uint8, error) {
	return _Leagues.Contract.GetTactics(&_Leagues.CallOpts, id)
}

// GetTeams is a free data retrieval call binding the contract method 0xf6c8d8de.
//
// Solidity: function getTeams(uint256 id) constant returns(uint256[])
func (_Leagues *LeaguesCaller) GetTeams(opts *bind.CallOpts, id *big.Int) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getTeams", id)
	return *ret0, err
}

// GetTeams is a free data retrieval call binding the contract method 0xf6c8d8de.
//
// Solidity: function getTeams(uint256 id) constant returns(uint256[])
func (_Leagues *LeaguesSession) GetTeams(id *big.Int) ([]*big.Int, error) {
	return _Leagues.Contract.GetTeams(&_Leagues.CallOpts, id)
}

// GetTeams is a free data retrieval call binding the contract method 0xf6c8d8de.
//
// Solidity: function getTeams(uint256 id) constant returns(uint256[])
func (_Leagues *LeaguesCallerSession) GetTeams(id *big.Int) ([]*big.Int, error) {
	return _Leagues.Contract.GetTeams(&_Leagues.CallOpts, id)
}

// GetTeamsInMatch is a free data retrieval call binding the contract method 0x84ec783c.
//
// Solidity: function getTeamsInMatch(uint256 id, uint256 matchday, uint256 matchIdx) constant returns(uint256 homeIdx, uint256 visitorIdx)
func (_Leagues *LeaguesCaller) GetTeamsInMatch(opts *bind.CallOpts, id *big.Int, matchday *big.Int, matchIdx *big.Int) (struct {
	HomeIdx    *big.Int
	VisitorIdx *big.Int
}, error) {
	ret := new(struct {
		HomeIdx    *big.Int
		VisitorIdx *big.Int
	})
	out := ret
	err := _Leagues.contract.Call(opts, out, "getTeamsInMatch", id, matchday, matchIdx)
	return *ret, err
}

// GetTeamsInMatch is a free data retrieval call binding the contract method 0x84ec783c.
//
// Solidity: function getTeamsInMatch(uint256 id, uint256 matchday, uint256 matchIdx) constant returns(uint256 homeIdx, uint256 visitorIdx)
func (_Leagues *LeaguesSession) GetTeamsInMatch(id *big.Int, matchday *big.Int, matchIdx *big.Int) (struct {
	HomeIdx    *big.Int
	VisitorIdx *big.Int
}, error) {
	return _Leagues.Contract.GetTeamsInMatch(&_Leagues.CallOpts, id, matchday, matchIdx)
}

// GetTeamsInMatch is a free data retrieval call binding the contract method 0x84ec783c.
//
// Solidity: function getTeamsInMatch(uint256 id, uint256 matchday, uint256 matchIdx) constant returns(uint256 homeIdx, uint256 visitorIdx)
func (_Leagues *LeaguesCallerSession) GetTeamsInMatch(id *big.Int, matchday *big.Int, matchIdx *big.Int) (struct {
	HomeIdx    *big.Int
	VisitorIdx *big.Int
}, error) {
	return _Leagues.Contract.GetTeamsInMatch(&_Leagues.CallOpts, id, matchday, matchIdx)
}

// GetUpdateBlock is a free data retrieval call binding the contract method 0x2b866890.
//
// Solidity: function getUpdateBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetUpdateBlock(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getUpdateBlock", id)
	return *ret0, err
}

// GetUpdateBlock is a free data retrieval call binding the contract method 0x2b866890.
//
// Solidity: function getUpdateBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesSession) GetUpdateBlock(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetUpdateBlock(&_Leagues.CallOpts, id)
}

// GetUpdateBlock is a free data retrieval call binding the contract method 0x2b866890.
//
// Solidity: function getUpdateBlock(uint256 id) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetUpdateBlock(id *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetUpdateBlock(&_Leagues.CallOpts, id)
}

// GetUpdater is a free data retrieval call binding the contract method 0xa8541e8e.
//
// Solidity: function getUpdater(uint256 id) constant returns(address)
func (_Leagues *LeaguesCaller) GetUpdater(opts *bind.CallOpts, id *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getUpdater", id)
	return *ret0, err
}

// GetUpdater is a free data retrieval call binding the contract method 0xa8541e8e.
//
// Solidity: function getUpdater(uint256 id) constant returns(address)
func (_Leagues *LeaguesSession) GetUpdater(id *big.Int) (common.Address, error) {
	return _Leagues.Contract.GetUpdater(&_Leagues.CallOpts, id)
}

// GetUpdater is a free data retrieval call binding the contract method 0xa8541e8e.
//
// Solidity: function getUpdater(uint256 id) constant returns(address)
func (_Leagues *LeaguesCallerSession) GetUpdater(id *big.Int) (common.Address, error) {
	return _Leagues.Contract.GetUpdater(&_Leagues.CallOpts, id)
}

// GetUsersAlongDataHash is a free data retrieval call binding the contract method 0x9e33c8a1.
//
// Solidity: function getUsersAlongDataHash(uint256 id) constant returns(bytes32)
func (_Leagues *LeaguesCaller) GetUsersAlongDataHash(opts *bind.CallOpts, id *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getUsersAlongDataHash", id)
	return *ret0, err
}

// GetUsersAlongDataHash is a free data retrieval call binding the contract method 0x9e33c8a1.
//
// Solidity: function getUsersAlongDataHash(uint256 id) constant returns(bytes32)
func (_Leagues *LeaguesSession) GetUsersAlongDataHash(id *big.Int) ([32]byte, error) {
	return _Leagues.Contract.GetUsersAlongDataHash(&_Leagues.CallOpts, id)
}

// GetUsersAlongDataHash is a free data retrieval call binding the contract method 0x9e33c8a1.
//
// Solidity: function getUsersAlongDataHash(uint256 id) constant returns(bytes32)
func (_Leagues *LeaguesCallerSession) GetUsersAlongDataHash(id *big.Int) ([32]byte, error) {
	return _Leagues.Contract.GetUsersAlongDataHash(&_Leagues.CallOpts, id)
}

// GetUsersInitDataHash is a free data retrieval call binding the contract method 0xbeec5522.
//
// Solidity: function getUsersInitDataHash(uint256 id) constant returns(bytes32)
func (_Leagues *LeaguesCaller) GetUsersInitDataHash(opts *bind.CallOpts, id *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getUsersInitDataHash", id)
	return *ret0, err
}

// GetUsersInitDataHash is a free data retrieval call binding the contract method 0xbeec5522.
//
// Solidity: function getUsersInitDataHash(uint256 id) constant returns(bytes32)
func (_Leagues *LeaguesSession) GetUsersInitDataHash(id *big.Int) ([32]byte, error) {
	return _Leagues.Contract.GetUsersInitDataHash(&_Leagues.CallOpts, id)
}

// GetUsersInitDataHash is a free data retrieval call binding the contract method 0xbeec5522.
//
// Solidity: function getUsersInitDataHash(uint256 id) constant returns(bytes32)
func (_Leagues *LeaguesCallerSession) GetUsersInitDataHash(id *big.Int) ([32]byte, error) {
	return _Leagues.Contract.GetUsersInitDataHash(&_Leagues.CallOpts, id)
}

// HasFinished is a free data retrieval call binding the contract method 0x126f648c.
//
// Solidity: function hasFinished(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCaller) HasFinished(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "hasFinished", id)
	return *ret0, err
}

// HasFinished is a free data retrieval call binding the contract method 0x126f648c.
//
// Solidity: function hasFinished(uint256 id) constant returns(bool)
func (_Leagues *LeaguesSession) HasFinished(id *big.Int) (bool, error) {
	return _Leagues.Contract.HasFinished(&_Leagues.CallOpts, id)
}

// HasFinished is a free data retrieval call binding the contract method 0x126f648c.
//
// Solidity: function hasFinished(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCallerSession) HasFinished(id *big.Int) (bool, error) {
	return _Leagues.Contract.HasFinished(&_Leagues.CallOpts, id)
}

// HasStarted is a free data retrieval call binding the contract method 0x51f41c09.
//
// Solidity: function hasStarted(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCaller) HasStarted(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "hasStarted", id)
	return *ret0, err
}

// HasStarted is a free data retrieval call binding the contract method 0x51f41c09.
//
// Solidity: function hasStarted(uint256 id) constant returns(bool)
func (_Leagues *LeaguesSession) HasStarted(id *big.Int) (bool, error) {
	return _Leagues.Contract.HasStarted(&_Leagues.CallOpts, id)
}

// HasStarted is a free data retrieval call binding the contract method 0x51f41c09.
//
// Solidity: function hasStarted(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCallerSession) HasStarted(id *big.Int) (bool, error) {
	return _Leagues.Contract.HasStarted(&_Leagues.CallOpts, id)
}

// HashDayState is a free data retrieval call binding the contract method 0xcc540c8f.
//
// Solidity: function hashDayState(uint256[] state) constant returns(bytes32)
func (_Leagues *LeaguesCaller) HashDayState(opts *bind.CallOpts, state []*big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "hashDayState", state)
	return *ret0, err
}

// HashDayState is a free data retrieval call binding the contract method 0xcc540c8f.
//
// Solidity: function hashDayState(uint256[] state) constant returns(bytes32)
func (_Leagues *LeaguesSession) HashDayState(state []*big.Int) ([32]byte, error) {
	return _Leagues.Contract.HashDayState(&_Leagues.CallOpts, state)
}

// HashDayState is a free data retrieval call binding the contract method 0xcc540c8f.
//
// Solidity: function hashDayState(uint256[] state) constant returns(bytes32)
func (_Leagues *LeaguesCallerSession) HashDayState(state []*big.Int) ([32]byte, error) {
	return _Leagues.Contract.HashDayState(&_Leagues.CallOpts, state)
}

// HashInitState is a free data retrieval call binding the contract method 0xc46a6939.
//
// Solidity: function hashInitState(uint256[] state) constant returns(bytes32)
func (_Leagues *LeaguesCaller) HashInitState(opts *bind.CallOpts, state []*big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "hashInitState", state)
	return *ret0, err
}

// HashInitState is a free data retrieval call binding the contract method 0xc46a6939.
//
// Solidity: function hashInitState(uint256[] state) constant returns(bytes32)
func (_Leagues *LeaguesSession) HashInitState(state []*big.Int) ([32]byte, error) {
	return _Leagues.Contract.HashInitState(&_Leagues.CallOpts, state)
}

// HashInitState is a free data retrieval call binding the contract method 0xc46a6939.
//
// Solidity: function hashInitState(uint256[] state) constant returns(bytes32)
func (_Leagues *LeaguesCallerSession) HashInitState(state []*big.Int) ([32]byte, error) {
	return _Leagues.Contract.HashInitState(&_Leagues.CallOpts, state)
}

// HashTactics is a free data retrieval call binding the contract method 0x49b5bda9.
//
// Solidity: function hashTactics(uint256[3][] tactics) constant returns(bytes32)
func (_Leagues *LeaguesCaller) HashTactics(opts *bind.CallOpts, tactics [][3]*big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "hashTactics", tactics)
	return *ret0, err
}

// HashTactics is a free data retrieval call binding the contract method 0x49b5bda9.
//
// Solidity: function hashTactics(uint256[3][] tactics) constant returns(bytes32)
func (_Leagues *LeaguesSession) HashTactics(tactics [][3]*big.Int) ([32]byte, error) {
	return _Leagues.Contract.HashTactics(&_Leagues.CallOpts, tactics)
}

// HashTactics is a free data retrieval call binding the contract method 0x49b5bda9.
//
// Solidity: function hashTactics(uint256[3][] tactics) constant returns(bytes32)
func (_Leagues *LeaguesCallerSession) HashTactics(tactics [][3]*big.Int) ([32]byte, error) {
	return _Leagues.Contract.HashTactics(&_Leagues.CallOpts, tactics)
}

// HashUsersInitData is a free data retrieval call binding the contract method 0xb524951f.
//
// Solidity: function hashUsersInitData(uint256[] teamIds, uint8[3][] tactics) constant returns(bytes32)
func (_Leagues *LeaguesCaller) HashUsersInitData(opts *bind.CallOpts, teamIds []*big.Int, tactics [][3]uint8) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "hashUsersInitData", teamIds, tactics)
	return *ret0, err
}

// HashUsersInitData is a free data retrieval call binding the contract method 0xb524951f.
//
// Solidity: function hashUsersInitData(uint256[] teamIds, uint8[3][] tactics) constant returns(bytes32)
func (_Leagues *LeaguesSession) HashUsersInitData(teamIds []*big.Int, tactics [][3]uint8) ([32]byte, error) {
	return _Leagues.Contract.HashUsersInitData(&_Leagues.CallOpts, teamIds, tactics)
}

// HashUsersInitData is a free data retrieval call binding the contract method 0xb524951f.
//
// Solidity: function hashUsersInitData(uint256[] teamIds, uint8[3][] tactics) constant returns(bytes32)
func (_Leagues *LeaguesCallerSession) HashUsersInitData(teamIds []*big.Int, tactics [][3]uint8) ([32]byte, error) {
	return _Leagues.Contract.HashUsersInitData(&_Leagues.CallOpts, teamIds, tactics)
}

// IsUpdated is a free data retrieval call binding the contract method 0xd39daf7e.
//
// Solidity: function isUpdated(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCaller) IsUpdated(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "isUpdated", id)
	return *ret0, err
}

// IsUpdated is a free data retrieval call binding the contract method 0xd39daf7e.
//
// Solidity: function isUpdated(uint256 id) constant returns(bool)
func (_Leagues *LeaguesSession) IsUpdated(id *big.Int) (bool, error) {
	return _Leagues.Contract.IsUpdated(&_Leagues.CallOpts, id)
}

// IsUpdated is a free data retrieval call binding the contract method 0xd39daf7e.
//
// Solidity: function isUpdated(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCallerSession) IsUpdated(id *big.Int) (bool, error) {
	return _Leagues.Contract.IsUpdated(&_Leagues.CallOpts, id)
}

// IsVerified is a free data retrieval call binding the contract method 0x37b6d96b.
//
// Solidity: function isVerified(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCaller) IsVerified(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "isVerified", id)
	return *ret0, err
}

// IsVerified is a free data retrieval call binding the contract method 0x37b6d96b.
//
// Solidity: function isVerified(uint256 id) constant returns(bool)
func (_Leagues *LeaguesSession) IsVerified(id *big.Int) (bool, error) {
	return _Leagues.Contract.IsVerified(&_Leagues.CallOpts, id)
}

// IsVerified is a free data retrieval call binding the contract method 0x37b6d96b.
//
// Solidity: function isVerified(uint256 id) constant returns(bool)
func (_Leagues *LeaguesCallerSession) IsVerified(id *big.Int) (bool, error) {
	return _Leagues.Contract.IsVerified(&_Leagues.CallOpts, id)
}

// LeaguesCount is a free data retrieval call binding the contract method 0x336b5a65.
//
// Solidity: function leaguesCount() constant returns(uint256)
func (_Leagues *LeaguesCaller) LeaguesCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "leaguesCount")
	return *ret0, err
}

// LeaguesCount is a free data retrieval call binding the contract method 0x336b5a65.
//
// Solidity: function leaguesCount() constant returns(uint256)
func (_Leagues *LeaguesSession) LeaguesCount() (*big.Int, error) {
	return _Leagues.Contract.LeaguesCount(&_Leagues.CallOpts)
}

// LeaguesCount is a free data retrieval call binding the contract method 0x336b5a65.
//
// Solidity: function leaguesCount() constant returns(uint256)
func (_Leagues *LeaguesCallerSession) LeaguesCount() (*big.Int, error) {
	return _Leagues.Contract.LeaguesCount(&_Leagues.CallOpts)
}

// ScoresAppend is a free data retrieval call binding the contract method 0x33a4aba6.
//
// Solidity: function scoresAppend(uint16[] scores, uint16 score) constant returns(uint16[])
func (_Leagues *LeaguesCaller) ScoresAppend(opts *bind.CallOpts, scores []uint16, score uint16) ([]uint16, error) {
	var (
		ret0 = new([]uint16)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "scoresAppend", scores, score)
	return *ret0, err
}

// ScoresAppend is a free data retrieval call binding the contract method 0x33a4aba6.
//
// Solidity: function scoresAppend(uint16[] scores, uint16 score) constant returns(uint16[])
func (_Leagues *LeaguesSession) ScoresAppend(scores []uint16, score uint16) ([]uint16, error) {
	return _Leagues.Contract.ScoresAppend(&_Leagues.CallOpts, scores, score)
}

// ScoresAppend is a free data retrieval call binding the contract method 0x33a4aba6.
//
// Solidity: function scoresAppend(uint16[] scores, uint16 score) constant returns(uint16[])
func (_Leagues *LeaguesCallerSession) ScoresAppend(scores []uint16, score uint16) ([]uint16, error) {
	return _Leagues.Contract.ScoresAppend(&_Leagues.CallOpts, scores, score)
}

// ScoresConcat is a free data retrieval call binding the contract method 0xa338f237.
//
// Solidity: function scoresConcat(uint16[] target, uint16[] scores) constant returns(uint16[])
func (_Leagues *LeaguesCaller) ScoresConcat(opts *bind.CallOpts, target []uint16, scores []uint16) ([]uint16, error) {
	var (
		ret0 = new([]uint16)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "scoresConcat", target, scores)
	return *ret0, err
}

// ScoresConcat is a free data retrieval call binding the contract method 0xa338f237.
//
// Solidity: function scoresConcat(uint16[] target, uint16[] scores) constant returns(uint16[])
func (_Leagues *LeaguesSession) ScoresConcat(target []uint16, scores []uint16) ([]uint16, error) {
	return _Leagues.Contract.ScoresConcat(&_Leagues.CallOpts, target, scores)
}

// ScoresConcat is a free data retrieval call binding the contract method 0xa338f237.
//
// Solidity: function scoresConcat(uint16[] target, uint16[] scores) constant returns(uint16[])
func (_Leagues *LeaguesCallerSession) ScoresConcat(target []uint16, scores []uint16) ([]uint16, error) {
	return _Leagues.Contract.ScoresConcat(&_Leagues.CallOpts, target, scores)
}

// ScoresCreate is a free data retrieval call binding the contract method 0xe41f70bd.
//
// Solidity: function scoresCreate() constant returns(uint16[])
func (_Leagues *LeaguesCaller) ScoresCreate(opts *bind.CallOpts) ([]uint16, error) {
	var (
		ret0 = new([]uint16)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "scoresCreate")
	return *ret0, err
}

// ScoresCreate is a free data retrieval call binding the contract method 0xe41f70bd.
//
// Solidity: function scoresCreate() constant returns(uint16[])
func (_Leagues *LeaguesSession) ScoresCreate() ([]uint16, error) {
	return _Leagues.Contract.ScoresCreate(&_Leagues.CallOpts)
}

// ScoresCreate is a free data retrieval call binding the contract method 0xe41f70bd.
//
// Solidity: function scoresCreate() constant returns(uint16[])
func (_Leagues *LeaguesCallerSession) ScoresCreate() ([]uint16, error) {
	return _Leagues.Contract.ScoresCreate(&_Leagues.CallOpts)
}

// ScoresGetDay is a free data retrieval call binding the contract method 0x4a52e3ec.
//
// Solidity: function scoresGetDay(uint256 id, uint256 day) constant returns(uint16[] dayScores)
func (_Leagues *LeaguesCaller) ScoresGetDay(opts *bind.CallOpts, id *big.Int, day *big.Int) ([]uint16, error) {
	var (
		ret0 = new([]uint16)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "scoresGetDay", id, day)
	return *ret0, err
}

// ScoresGetDay is a free data retrieval call binding the contract method 0x4a52e3ec.
//
// Solidity: function scoresGetDay(uint256 id, uint256 day) constant returns(uint16[] dayScores)
func (_Leagues *LeaguesSession) ScoresGetDay(id *big.Int, day *big.Int) ([]uint16, error) {
	return _Leagues.Contract.ScoresGetDay(&_Leagues.CallOpts, id, day)
}

// ScoresGetDay is a free data retrieval call binding the contract method 0x4a52e3ec.
//
// Solidity: function scoresGetDay(uint256 id, uint256 day) constant returns(uint16[] dayScores)
func (_Leagues *LeaguesCallerSession) ScoresGetDay(id *big.Int, day *big.Int) ([]uint16, error) {
	return _Leagues.Contract.ScoresGetDay(&_Leagues.CallOpts, id, day)
}

// ChallengeInitStates is a paid mutator transaction binding the contract method 0xb1f65f39.
//
// Solidity: function challengeInitStates(uint256 id, uint256[] teamIds, uint8[3][] tactics, uint256[] dataToChallengeInitStates) returns()
func (_Leagues *LeaguesTransactor) ChallengeInitStates(opts *bind.TransactOpts, id *big.Int, teamIds []*big.Int, tactics [][3]uint8, dataToChallengeInitStates []*big.Int) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "challengeInitStates", id, teamIds, tactics, dataToChallengeInitStates)
}

// ChallengeInitStates is a paid mutator transaction binding the contract method 0xb1f65f39.
//
// Solidity: function challengeInitStates(uint256 id, uint256[] teamIds, uint8[3][] tactics, uint256[] dataToChallengeInitStates) returns()
func (_Leagues *LeaguesSession) ChallengeInitStates(id *big.Int, teamIds []*big.Int, tactics [][3]uint8, dataToChallengeInitStates []*big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.ChallengeInitStates(&_Leagues.TransactOpts, id, teamIds, tactics, dataToChallengeInitStates)
}

// ChallengeInitStates is a paid mutator transaction binding the contract method 0xb1f65f39.
//
// Solidity: function challengeInitStates(uint256 id, uint256[] teamIds, uint8[3][] tactics, uint256[] dataToChallengeInitStates) returns()
func (_Leagues *LeaguesTransactorSession) ChallengeInitStates(id *big.Int, teamIds []*big.Int, tactics [][3]uint8, dataToChallengeInitStates []*big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.ChallengeInitStates(&_Leagues.TransactOpts, id, teamIds, tactics, dataToChallengeInitStates)
}

// ChallengeMatchdayStates is a paid mutator transaction binding the contract method 0x82b36ed8.
//
// Solidity: function challengeMatchdayStates(uint256 id, uint256[] usersInitDataTeamIds, uint8[3][] usersInitDataTactics, uint256[] usersAlongDataTeamIds, uint8[3][] usersAlongDataTactics, uint256[] usersAlongDataBlocks, uint256 leagueDay, uint256[] prevMatchdayStates) returns()
func (_Leagues *LeaguesTransactor) ChallengeMatchdayStates(opts *bind.TransactOpts, id *big.Int, usersInitDataTeamIds []*big.Int, usersInitDataTactics [][3]uint8, usersAlongDataTeamIds []*big.Int, usersAlongDataTactics [][3]uint8, usersAlongDataBlocks []*big.Int, leagueDay *big.Int, prevMatchdayStates []*big.Int) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "challengeMatchdayStates", id, usersInitDataTeamIds, usersInitDataTactics, usersAlongDataTeamIds, usersAlongDataTactics, usersAlongDataBlocks, leagueDay, prevMatchdayStates)
}

// ChallengeMatchdayStates is a paid mutator transaction binding the contract method 0x82b36ed8.
//
// Solidity: function challengeMatchdayStates(uint256 id, uint256[] usersInitDataTeamIds, uint8[3][] usersInitDataTactics, uint256[] usersAlongDataTeamIds, uint8[3][] usersAlongDataTactics, uint256[] usersAlongDataBlocks, uint256 leagueDay, uint256[] prevMatchdayStates) returns()
func (_Leagues *LeaguesSession) ChallengeMatchdayStates(id *big.Int, usersInitDataTeamIds []*big.Int, usersInitDataTactics [][3]uint8, usersAlongDataTeamIds []*big.Int, usersAlongDataTactics [][3]uint8, usersAlongDataBlocks []*big.Int, leagueDay *big.Int, prevMatchdayStates []*big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.ChallengeMatchdayStates(&_Leagues.TransactOpts, id, usersInitDataTeamIds, usersInitDataTactics, usersAlongDataTeamIds, usersAlongDataTactics, usersAlongDataBlocks, leagueDay, prevMatchdayStates)
}

// ChallengeMatchdayStates is a paid mutator transaction binding the contract method 0x82b36ed8.
//
// Solidity: function challengeMatchdayStates(uint256 id, uint256[] usersInitDataTeamIds, uint8[3][] usersInitDataTactics, uint256[] usersAlongDataTeamIds, uint8[3][] usersAlongDataTactics, uint256[] usersAlongDataBlocks, uint256 leagueDay, uint256[] prevMatchdayStates) returns()
func (_Leagues *LeaguesTransactorSession) ChallengeMatchdayStates(id *big.Int, usersInitDataTeamIds []*big.Int, usersInitDataTactics [][3]uint8, usersAlongDataTeamIds []*big.Int, usersAlongDataTactics [][3]uint8, usersAlongDataBlocks []*big.Int, leagueDay *big.Int, prevMatchdayStates []*big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.ChallengeMatchdayStates(&_Leagues.TransactOpts, id, usersInitDataTeamIds, usersInitDataTactics, usersAlongDataTeamIds, usersAlongDataTactics, usersAlongDataBlocks, leagueDay, prevMatchdayStates)
}

// Create is a paid mutator transaction binding the contract method 0x77f2a891.
//
// Solidity: function create(uint256 id, uint256 initBlock, uint256 step, uint256[] teamIds, uint8[3][] tactics) returns()
func (_Leagues *LeaguesTransactor) Create(opts *bind.TransactOpts, id *big.Int, initBlock *big.Int, step *big.Int, teamIds []*big.Int, tactics [][3]uint8) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "create", id, initBlock, step, teamIds, tactics)
}

// Create is a paid mutator transaction binding the contract method 0x77f2a891.
//
// Solidity: function create(uint256 id, uint256 initBlock, uint256 step, uint256[] teamIds, uint8[3][] tactics) returns()
func (_Leagues *LeaguesSession) Create(id *big.Int, initBlock *big.Int, step *big.Int, teamIds []*big.Int, tactics [][3]uint8) (*types.Transaction, error) {
	return _Leagues.Contract.Create(&_Leagues.TransactOpts, id, initBlock, step, teamIds, tactics)
}

// Create is a paid mutator transaction binding the contract method 0x77f2a891.
//
// Solidity: function create(uint256 id, uint256 initBlock, uint256 step, uint256[] teamIds, uint8[3][] tactics) returns()
func (_Leagues *LeaguesTransactorSession) Create(id *big.Int, initBlock *big.Int, step *big.Int, teamIds []*big.Int, tactics [][3]uint8) (*types.Transaction, error) {
	return _Leagues.Contract.Create(&_Leagues.TransactOpts, id, initBlock, step, teamIds, tactics)
}

// GetInitPlayerStates is a paid mutator transaction binding the contract method 0xf025f7f4.
//
// Solidity: function getInitPlayerStates(uint256 id, uint256[] teamIds, uint8[3][] tactics, uint256[] dataToChallengeInitStates) returns(uint256[] state)
func (_Leagues *LeaguesTransactor) GetInitPlayerStates(opts *bind.TransactOpts, id *big.Int, teamIds []*big.Int, tactics [][3]uint8, dataToChallengeInitStates []*big.Int) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "getInitPlayerStates", id, teamIds, tactics, dataToChallengeInitStates)
}

// GetInitPlayerStates is a paid mutator transaction binding the contract method 0xf025f7f4.
//
// Solidity: function getInitPlayerStates(uint256 id, uint256[] teamIds, uint8[3][] tactics, uint256[] dataToChallengeInitStates) returns(uint256[] state)
func (_Leagues *LeaguesSession) GetInitPlayerStates(id *big.Int, teamIds []*big.Int, tactics [][3]uint8, dataToChallengeInitStates []*big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.GetInitPlayerStates(&_Leagues.TransactOpts, id, teamIds, tactics, dataToChallengeInitStates)
}

// GetInitPlayerStates is a paid mutator transaction binding the contract method 0xf025f7f4.
//
// Solidity: function getInitPlayerStates(uint256 id, uint256[] teamIds, uint8[3][] tactics, uint256[] dataToChallengeInitStates) returns(uint256[] state)
func (_Leagues *LeaguesTransactorSession) GetInitPlayerStates(id *big.Int, teamIds []*big.Int, tactics [][3]uint8, dataToChallengeInitStates []*big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.GetInitPlayerStates(&_Leagues.TransactOpts, id, teamIds, tactics, dataToChallengeInitStates)
}

// ResetUpdater is a paid mutator transaction binding the contract method 0x28d3be79.
//
// Solidity: function resetUpdater(uint256 id) returns()
func (_Leagues *LeaguesTransactor) ResetUpdater(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "resetUpdater", id)
}

// ResetUpdater is a paid mutator transaction binding the contract method 0x28d3be79.
//
// Solidity: function resetUpdater(uint256 id) returns()
func (_Leagues *LeaguesSession) ResetUpdater(id *big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.ResetUpdater(&_Leagues.TransactOpts, id)
}

// ResetUpdater is a paid mutator transaction binding the contract method 0x28d3be79.
//
// Solidity: function resetUpdater(uint256 id) returns()
func (_Leagues *LeaguesTransactorSession) ResetUpdater(id *big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.ResetUpdater(&_Leagues.TransactOpts, id)
}

// SetStakersContract is a paid mutator transaction binding the contract method 0xc1fb64ca.
//
// Solidity: function setStakersContract(address stakersContract) returns()
func (_Leagues *LeaguesTransactor) SetStakersContract(opts *bind.TransactOpts, stakersContract common.Address) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "setStakersContract", stakersContract)
}

// SetStakersContract is a paid mutator transaction binding the contract method 0xc1fb64ca.
//
// Solidity: function setStakersContract(address stakersContract) returns()
func (_Leagues *LeaguesSession) SetStakersContract(stakersContract common.Address) (*types.Transaction, error) {
	return _Leagues.Contract.SetStakersContract(&_Leagues.TransactOpts, stakersContract)
}

// SetStakersContract is a paid mutator transaction binding the contract method 0xc1fb64ca.
//
// Solidity: function setStakersContract(address stakersContract) returns()
func (_Leagues *LeaguesTransactorSession) SetStakersContract(stakersContract common.Address) (*types.Transaction, error) {
	return _Leagues.Contract.SetStakersContract(&_Leagues.TransactOpts, stakersContract)
}

// UpdateLeague is a paid mutator transaction binding the contract method 0x94f37021.
//
// Solidity: function updateLeague(uint256 id, bytes32 initStateHash, bytes32[] dayStateHashes, uint16[] scores, bool isLie) returns()
func (_Leagues *LeaguesTransactor) UpdateLeague(opts *bind.TransactOpts, id *big.Int, initStateHash [32]byte, dayStateHashes [][32]byte, scores []uint16, isLie bool) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "updateLeague", id, initStateHash, dayStateHashes, scores, isLie)
}

// UpdateLeague is a paid mutator transaction binding the contract method 0x94f37021.
//
// Solidity: function updateLeague(uint256 id, bytes32 initStateHash, bytes32[] dayStateHashes, uint16[] scores, bool isLie) returns()
func (_Leagues *LeaguesSession) UpdateLeague(id *big.Int, initStateHash [32]byte, dayStateHashes [][32]byte, scores []uint16, isLie bool) (*types.Transaction, error) {
	return _Leagues.Contract.UpdateLeague(&_Leagues.TransactOpts, id, initStateHash, dayStateHashes, scores, isLie)
}

// UpdateLeague is a paid mutator transaction binding the contract method 0x94f37021.
//
// Solidity: function updateLeague(uint256 id, bytes32 initStateHash, bytes32[] dayStateHashes, uint16[] scores, bool isLie) returns()
func (_Leagues *LeaguesTransactorSession) UpdateLeague(id *big.Int, initStateHash [32]byte, dayStateHashes [][32]byte, scores []uint16, isLie bool) (*types.Transaction, error) {
	return _Leagues.Contract.UpdateLeague(&_Leagues.TransactOpts, id, initStateHash, dayStateHashes, scores, isLie)
}

// UpdateUsersAlongDataHash is a paid mutator transaction binding the contract method 0xb2df2c57.
//
// Solidity: function updateUsersAlongDataHash(uint256 id, uint256[] teamIds, uint8[3][] tactics, uint256[] blocks) returns()
func (_Leagues *LeaguesTransactor) UpdateUsersAlongDataHash(opts *bind.TransactOpts, id *big.Int, teamIds []*big.Int, tactics [][3]uint8, blocks []*big.Int) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "updateUsersAlongDataHash", id, teamIds, tactics, blocks)
}

// UpdateUsersAlongDataHash is a paid mutator transaction binding the contract method 0xb2df2c57.
//
// Solidity: function updateUsersAlongDataHash(uint256 id, uint256[] teamIds, uint8[3][] tactics, uint256[] blocks) returns()
func (_Leagues *LeaguesSession) UpdateUsersAlongDataHash(id *big.Int, teamIds []*big.Int, tactics [][3]uint8, blocks []*big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.UpdateUsersAlongDataHash(&_Leagues.TransactOpts, id, teamIds, tactics, blocks)
}

// UpdateUsersAlongDataHash is a paid mutator transaction binding the contract method 0xb2df2c57.
//
// Solidity: function updateUsersAlongDataHash(uint256 id, uint256[] teamIds, uint8[3][] tactics, uint256[] blocks) returns()
func (_Leagues *LeaguesTransactorSession) UpdateUsersAlongDataHash(id *big.Int, teamIds []*big.Int, tactics [][3]uint8, blocks []*big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.UpdateUsersAlongDataHash(&_Leagues.TransactOpts, id, teamIds, tactics, blocks)
}

// LeaguesLeagueCreatedIterator is returned from FilterLeagueCreated and is used to iterate over the raw logs and unpacked data for LeagueCreated events raised by the Leagues contract.
type LeaguesLeagueCreatedIterator struct {
	Event *LeaguesLeagueCreated // Event containing the contract specifics and raw log

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
func (it *LeaguesLeagueCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LeaguesLeagueCreated)
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
		it.Event = new(LeaguesLeagueCreated)
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
func (it *LeaguesLeagueCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LeaguesLeagueCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LeaguesLeagueCreated represents a LeagueCreated event raised by the Leagues contract.
type LeaguesLeagueCreated struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLeagueCreated is a free log retrieval operation binding the contract event 0x5d69f37aa0f3d80654d5e87cc70b8464663e1e43be29aae3a06dcecef8471906.
//
// Solidity: event LeagueCreated(uint256 id)
func (_Leagues *LeaguesFilterer) FilterLeagueCreated(opts *bind.FilterOpts) (*LeaguesLeagueCreatedIterator, error) {

	logs, sub, err := _Leagues.contract.FilterLogs(opts, "LeagueCreated")
	if err != nil {
		return nil, err
	}
	return &LeaguesLeagueCreatedIterator{contract: _Leagues.contract, event: "LeagueCreated", logs: logs, sub: sub}, nil
}

// WatchLeagueCreated is a free log subscription operation binding the contract event 0x5d69f37aa0f3d80654d5e87cc70b8464663e1e43be29aae3a06dcecef8471906.
//
// Solidity: event LeagueCreated(uint256 id)
func (_Leagues *LeaguesFilterer) WatchLeagueCreated(opts *bind.WatchOpts, sink chan<- *LeaguesLeagueCreated) (event.Subscription, error) {

	logs, sub, err := _Leagues.contract.WatchLogs(opts, "LeagueCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LeaguesLeagueCreated)
				if err := _Leagues.contract.UnpackLog(event, "LeagueCreated", log); err != nil {
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
