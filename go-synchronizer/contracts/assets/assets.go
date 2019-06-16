// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package assets

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

// AssetsABI is the input ABI used to generate the binding from.
const AssetsABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"NUM_SKILLS\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"PLAYERS_PER_TEAM\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"playerState\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"teamName\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"teamId\",\"type\":\"uint256\"}],\"name\":\"TeamCreation\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamId\",\"type\":\"uint256\"}],\"name\":\"getTeamCreationTimestamp\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamId\",\"type\":\"uint256\"}],\"name\":\"getTeamCurrentHistory\",\"outputs\":[{\"name\":\"currentLeagueId\",\"type\":\"uint256\"},{\"name\":\"posInCurrentLeague\",\"type\":\"uint8\"},{\"name\":\"prevLeagueId\",\"type\":\"uint256\"},{\"name\":\"posInPrevLeague\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"playerId0\",\"type\":\"uint256\"},{\"name\":\"playerId1\",\"type\":\"uint256\"}],\"name\":\"exchangePlayersTeams\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"createTeam\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"teamId\",\"type\":\"uint256\"},{\"name\":\"leagueId\",\"type\":\"uint256\"},{\"name\":\"posInLeague\",\"type\":\"uint8\"}],\"name\":\"signToLeague\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getTeamOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"countTeams\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamId\",\"type\":\"uint256\"}],\"name\":\"getTeamName\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamId\",\"type\":\"uint256\"},{\"name\":\"posInTeam\",\"type\":\"uint8\"}],\"name\":\"getPlayerIdFromTeamIdAndPos\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamId\",\"type\":\"uint256\"}],\"name\":\"getTeamPlayerIds\",\"outputs\":[{\"name\":\"playerIds\",\"type\":\"uint256[11]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerId\",\"type\":\"uint256\"}],\"name\":\"getPlayerState\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Assets is an auto generated Go binding around an Ethereum contract.
type Assets struct {
	AssetsCaller     // Read-only binding to the contract
	AssetsTransactor // Write-only binding to the contract
	AssetsFilterer   // Log filterer for contract events
}

// AssetsCaller is an auto generated read-only Go binding around an Ethereum contract.
type AssetsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AssetsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AssetsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AssetsSession struct {
	Contract     *Assets           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssetsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AssetsCallerSession struct {
	Contract *AssetsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AssetsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AssetsTransactorSession struct {
	Contract     *AssetsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssetsRaw is an auto generated low-level Go binding around an Ethereum contract.
type AssetsRaw struct {
	Contract *Assets // Generic contract binding to access the raw methods on
}

// AssetsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AssetsCallerRaw struct {
	Contract *AssetsCaller // Generic read-only contract binding to access the raw methods on
}

// AssetsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AssetsTransactorRaw struct {
	Contract *AssetsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAssets creates a new instance of Assets, bound to a specific deployed contract.
func NewAssets(address common.Address, backend bind.ContractBackend) (*Assets, error) {
	contract, err := bindAssets(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Assets{AssetsCaller: AssetsCaller{contract: contract}, AssetsTransactor: AssetsTransactor{contract: contract}, AssetsFilterer: AssetsFilterer{contract: contract}}, nil
}

// NewAssetsCaller creates a new read-only instance of Assets, bound to a specific deployed contract.
func NewAssetsCaller(address common.Address, caller bind.ContractCaller) (*AssetsCaller, error) {
	contract, err := bindAssets(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AssetsCaller{contract: contract}, nil
}

// NewAssetsTransactor creates a new write-only instance of Assets, bound to a specific deployed contract.
func NewAssetsTransactor(address common.Address, transactor bind.ContractTransactor) (*AssetsTransactor, error) {
	contract, err := bindAssets(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AssetsTransactor{contract: contract}, nil
}

// NewAssetsFilterer creates a new log filterer instance of Assets, bound to a specific deployed contract.
func NewAssetsFilterer(address common.Address, filterer bind.ContractFilterer) (*AssetsFilterer, error) {
	contract, err := bindAssets(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AssetsFilterer{contract: contract}, nil
}

// bindAssets binds a generic wrapper to an already deployed contract.
func bindAssets(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AssetsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Assets *AssetsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Assets.Contract.AssetsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Assets *AssetsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assets.Contract.AssetsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Assets *AssetsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Assets.Contract.AssetsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Assets *AssetsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Assets.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Assets *AssetsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assets.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Assets *AssetsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Assets.Contract.contract.Transact(opts, method, params...)
}

// NUMSKILLS is a free data retrieval call binding the contract method 0x528afa3f.
//
// Solidity: function NUM_SKILLS() constant returns(uint8)
func (_Assets *AssetsCaller) NUMSKILLS(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "NUM_SKILLS")
	return *ret0, err
}

// NUMSKILLS is a free data retrieval call binding the contract method 0x528afa3f.
//
// Solidity: function NUM_SKILLS() constant returns(uint8)
func (_Assets *AssetsSession) NUMSKILLS() (uint8, error) {
	return _Assets.Contract.NUMSKILLS(&_Assets.CallOpts)
}

// NUMSKILLS is a free data retrieval call binding the contract method 0x528afa3f.
//
// Solidity: function NUM_SKILLS() constant returns(uint8)
func (_Assets *AssetsCallerSession) NUMSKILLS() (uint8, error) {
	return _Assets.Contract.NUMSKILLS(&_Assets.CallOpts)
}

// PLAYERSPERTEAM is a free data retrieval call binding the contract method 0xab6eeb00.
//
// Solidity: function PLAYERS_PER_TEAM() constant returns(uint8)
func (_Assets *AssetsCaller) PLAYERSPERTEAM(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "PLAYERS_PER_TEAM")
	return *ret0, err
}

// PLAYERSPERTEAM is a free data retrieval call binding the contract method 0xab6eeb00.
//
// Solidity: function PLAYERS_PER_TEAM() constant returns(uint8)
func (_Assets *AssetsSession) PLAYERSPERTEAM() (uint8, error) {
	return _Assets.Contract.PLAYERSPERTEAM(&_Assets.CallOpts)
}

// PLAYERSPERTEAM is a free data retrieval call binding the contract method 0xab6eeb00.
//
// Solidity: function PLAYERS_PER_TEAM() constant returns(uint8)
func (_Assets *AssetsCallerSession) PLAYERSPERTEAM() (uint8, error) {
	return _Assets.Contract.PLAYERSPERTEAM(&_Assets.CallOpts)
}

// CountTeams is a free data retrieval call binding the contract method 0x16cb9b9d.
//
// Solidity: function countTeams() constant returns(uint256)
func (_Assets *AssetsCaller) CountTeams(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "countTeams")
	return *ret0, err
}

// CountTeams is a free data retrieval call binding the contract method 0x16cb9b9d.
//
// Solidity: function countTeams() constant returns(uint256)
func (_Assets *AssetsSession) CountTeams() (*big.Int, error) {
	return _Assets.Contract.CountTeams(&_Assets.CallOpts)
}

// CountTeams is a free data retrieval call binding the contract method 0x16cb9b9d.
//
// Solidity: function countTeams() constant returns(uint256)
func (_Assets *AssetsCallerSession) CountTeams() (*big.Int, error) {
	return _Assets.Contract.CountTeams(&_Assets.CallOpts)
}

// GetPlayerIdFromTeamIdAndPos is a free data retrieval call binding the contract method 0x988be518.
//
// Solidity: function getPlayerIdFromTeamIdAndPos(uint256 teamId, uint8 posInTeam) constant returns(uint256)
func (_Assets *AssetsCaller) GetPlayerIdFromTeamIdAndPos(opts *bind.CallOpts, teamId *big.Int, posInTeam uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPlayerIdFromTeamIdAndPos", teamId, posInTeam)
	return *ret0, err
}

// GetPlayerIdFromTeamIdAndPos is a free data retrieval call binding the contract method 0x988be518.
//
// Solidity: function getPlayerIdFromTeamIdAndPos(uint256 teamId, uint8 posInTeam) constant returns(uint256)
func (_Assets *AssetsSession) GetPlayerIdFromTeamIdAndPos(teamId *big.Int, posInTeam uint8) (*big.Int, error) {
	return _Assets.Contract.GetPlayerIdFromTeamIdAndPos(&_Assets.CallOpts, teamId, posInTeam)
}

// GetPlayerIdFromTeamIdAndPos is a free data retrieval call binding the contract method 0x988be518.
//
// Solidity: function getPlayerIdFromTeamIdAndPos(uint256 teamId, uint8 posInTeam) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPlayerIdFromTeamIdAndPos(teamId *big.Int, posInTeam uint8) (*big.Int, error) {
	return _Assets.Contract.GetPlayerIdFromTeamIdAndPos(&_Assets.CallOpts, teamId, posInTeam)
}

// GetPlayerState is a free data retrieval call binding the contract method 0xec7ecec5.
//
// Solidity: function getPlayerState(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCaller) GetPlayerState(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPlayerState", playerId)
	return *ret0, err
}

// GetPlayerState is a free data retrieval call binding the contract method 0xec7ecec5.
//
// Solidity: function getPlayerState(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsSession) GetPlayerState(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerState(&_Assets.CallOpts, playerId)
}

// GetPlayerState is a free data retrieval call binding the contract method 0xec7ecec5.
//
// Solidity: function getPlayerState(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPlayerState(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerState(&_Assets.CallOpts, playerId)
}

// GetTeamCreationTimestamp is a free data retrieval call binding the contract method 0x93ae14fd.
//
// Solidity: function getTeamCreationTimestamp(uint256 teamId) constant returns(uint256)
func (_Assets *AssetsCaller) GetTeamCreationTimestamp(opts *bind.CallOpts, teamId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getTeamCreationTimestamp", teamId)
	return *ret0, err
}

// GetTeamCreationTimestamp is a free data retrieval call binding the contract method 0x93ae14fd.
//
// Solidity: function getTeamCreationTimestamp(uint256 teamId) constant returns(uint256)
func (_Assets *AssetsSession) GetTeamCreationTimestamp(teamId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetTeamCreationTimestamp(&_Assets.CallOpts, teamId)
}

// GetTeamCreationTimestamp is a free data retrieval call binding the contract method 0x93ae14fd.
//
// Solidity: function getTeamCreationTimestamp(uint256 teamId) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetTeamCreationTimestamp(teamId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetTeamCreationTimestamp(&_Assets.CallOpts, teamId)
}

// GetTeamCurrentHistory is a free data retrieval call binding the contract method 0xf97fa22f.
//
// Solidity: function getTeamCurrentHistory(uint256 teamId) constant returns(uint256 currentLeagueId, uint8 posInCurrentLeague, uint256 prevLeagueId, uint8 posInPrevLeague)
func (_Assets *AssetsCaller) GetTeamCurrentHistory(opts *bind.CallOpts, teamId *big.Int) (struct {
	CurrentLeagueId    *big.Int
	PosInCurrentLeague uint8
	PrevLeagueId       *big.Int
	PosInPrevLeague    uint8
}, error) {
	ret := new(struct {
		CurrentLeagueId    *big.Int
		PosInCurrentLeague uint8
		PrevLeagueId       *big.Int
		PosInPrevLeague    uint8
	})
	out := ret
	err := _Assets.contract.Call(opts, out, "getTeamCurrentHistory", teamId)
	return *ret, err
}

// GetTeamCurrentHistory is a free data retrieval call binding the contract method 0xf97fa22f.
//
// Solidity: function getTeamCurrentHistory(uint256 teamId) constant returns(uint256 currentLeagueId, uint8 posInCurrentLeague, uint256 prevLeagueId, uint8 posInPrevLeague)
func (_Assets *AssetsSession) GetTeamCurrentHistory(teamId *big.Int) (struct {
	CurrentLeagueId    *big.Int
	PosInCurrentLeague uint8
	PrevLeagueId       *big.Int
	PosInPrevLeague    uint8
}, error) {
	return _Assets.Contract.GetTeamCurrentHistory(&_Assets.CallOpts, teamId)
}

// GetTeamCurrentHistory is a free data retrieval call binding the contract method 0xf97fa22f.
//
// Solidity: function getTeamCurrentHistory(uint256 teamId) constant returns(uint256 currentLeagueId, uint8 posInCurrentLeague, uint256 prevLeagueId, uint8 posInPrevLeague)
func (_Assets *AssetsCallerSession) GetTeamCurrentHistory(teamId *big.Int) (struct {
	CurrentLeagueId    *big.Int
	PosInCurrentLeague uint8
	PrevLeagueId       *big.Int
	PosInPrevLeague    uint8
}, error) {
	return _Assets.Contract.GetTeamCurrentHistory(&_Assets.CallOpts, teamId)
}

// GetTeamName is a free data retrieval call binding the contract method 0x7ab60d48.
//
// Solidity: function getTeamName(uint256 teamId) constant returns(string)
func (_Assets *AssetsCaller) GetTeamName(opts *bind.CallOpts, teamId *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getTeamName", teamId)
	return *ret0, err
}

// GetTeamName is a free data retrieval call binding the contract method 0x7ab60d48.
//
// Solidity: function getTeamName(uint256 teamId) constant returns(string)
func (_Assets *AssetsSession) GetTeamName(teamId *big.Int) (string, error) {
	return _Assets.Contract.GetTeamName(&_Assets.CallOpts, teamId)
}

// GetTeamName is a free data retrieval call binding the contract method 0x7ab60d48.
//
// Solidity: function getTeamName(uint256 teamId) constant returns(string)
func (_Assets *AssetsCallerSession) GetTeamName(teamId *big.Int) (string, error) {
	return _Assets.Contract.GetTeamName(&_Assets.CallOpts, teamId)
}

// GetTeamOwner is a free data retrieval call binding the contract method 0xb2a93a1f.
//
// Solidity: function getTeamOwner(string name) constant returns(address)
func (_Assets *AssetsCaller) GetTeamOwner(opts *bind.CallOpts, name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getTeamOwner", name)
	return *ret0, err
}

// GetTeamOwner is a free data retrieval call binding the contract method 0xb2a93a1f.
//
// Solidity: function getTeamOwner(string name) constant returns(address)
func (_Assets *AssetsSession) GetTeamOwner(name string) (common.Address, error) {
	return _Assets.Contract.GetTeamOwner(&_Assets.CallOpts, name)
}

// GetTeamOwner is a free data retrieval call binding the contract method 0xb2a93a1f.
//
// Solidity: function getTeamOwner(string name) constant returns(address)
func (_Assets *AssetsCallerSession) GetTeamOwner(name string) (common.Address, error) {
	return _Assets.Contract.GetTeamOwner(&_Assets.CallOpts, name)
}

// GetTeamPlayerIds is a free data retrieval call binding the contract method 0x03908478.
//
// Solidity: function getTeamPlayerIds(uint256 teamId) constant returns(uint256[11] playerIds)
func (_Assets *AssetsCaller) GetTeamPlayerIds(opts *bind.CallOpts, teamId *big.Int) ([11]*big.Int, error) {
	var (
		ret0 = new([11]*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getTeamPlayerIds", teamId)
	return *ret0, err
}

// GetTeamPlayerIds is a free data retrieval call binding the contract method 0x03908478.
//
// Solidity: function getTeamPlayerIds(uint256 teamId) constant returns(uint256[11] playerIds)
func (_Assets *AssetsSession) GetTeamPlayerIds(teamId *big.Int) ([11]*big.Int, error) {
	return _Assets.Contract.GetTeamPlayerIds(&_Assets.CallOpts, teamId)
}

// GetTeamPlayerIds is a free data retrieval call binding the contract method 0x03908478.
//
// Solidity: function getTeamPlayerIds(uint256 teamId) constant returns(uint256[11] playerIds)
func (_Assets *AssetsCallerSession) GetTeamPlayerIds(teamId *big.Int) ([11]*big.Int, error) {
	return _Assets.Contract.GetTeamPlayerIds(&_Assets.CallOpts, teamId)
}

// CreateTeam is a paid mutator transaction binding the contract method 0xa66b1e78.
//
// Solidity: function createTeam(string name, address owner) returns()
func (_Assets *AssetsTransactor) CreateTeam(opts *bind.TransactOpts, name string, owner common.Address) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "createTeam", name, owner)
}

// CreateTeam is a paid mutator transaction binding the contract method 0xa66b1e78.
//
// Solidity: function createTeam(string name, address owner) returns()
func (_Assets *AssetsSession) CreateTeam(name string, owner common.Address) (*types.Transaction, error) {
	return _Assets.Contract.CreateTeam(&_Assets.TransactOpts, name, owner)
}

// CreateTeam is a paid mutator transaction binding the contract method 0xa66b1e78.
//
// Solidity: function createTeam(string name, address owner) returns()
func (_Assets *AssetsTransactorSession) CreateTeam(name string, owner common.Address) (*types.Transaction, error) {
	return _Assets.Contract.CreateTeam(&_Assets.TransactOpts, name, owner)
}

// ExchangePlayersTeams is a paid mutator transaction binding the contract method 0x789875b5.
//
// Solidity: function exchangePlayersTeams(uint256 playerId0, uint256 playerId1) returns()
func (_Assets *AssetsTransactor) ExchangePlayersTeams(opts *bind.TransactOpts, playerId0 *big.Int, playerId1 *big.Int) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "exchangePlayersTeams", playerId0, playerId1)
}

// ExchangePlayersTeams is a paid mutator transaction binding the contract method 0x789875b5.
//
// Solidity: function exchangePlayersTeams(uint256 playerId0, uint256 playerId1) returns()
func (_Assets *AssetsSession) ExchangePlayersTeams(playerId0 *big.Int, playerId1 *big.Int) (*types.Transaction, error) {
	return _Assets.Contract.ExchangePlayersTeams(&_Assets.TransactOpts, playerId0, playerId1)
}

// ExchangePlayersTeams is a paid mutator transaction binding the contract method 0x789875b5.
//
// Solidity: function exchangePlayersTeams(uint256 playerId0, uint256 playerId1) returns()
func (_Assets *AssetsTransactorSession) ExchangePlayersTeams(playerId0 *big.Int, playerId1 *big.Int) (*types.Transaction, error) {
	return _Assets.Contract.ExchangePlayersTeams(&_Assets.TransactOpts, playerId0, playerId1)
}

// SignToLeague is a paid mutator transaction binding the contract method 0x02154127.
//
// Solidity: function signToLeague(uint256 teamId, uint256 leagueId, uint8 posInLeague) returns()
func (_Assets *AssetsTransactor) SignToLeague(opts *bind.TransactOpts, teamId *big.Int, leagueId *big.Int, posInLeague uint8) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "signToLeague", teamId, leagueId, posInLeague)
}

// SignToLeague is a paid mutator transaction binding the contract method 0x02154127.
//
// Solidity: function signToLeague(uint256 teamId, uint256 leagueId, uint8 posInLeague) returns()
func (_Assets *AssetsSession) SignToLeague(teamId *big.Int, leagueId *big.Int, posInLeague uint8) (*types.Transaction, error) {
	return _Assets.Contract.SignToLeague(&_Assets.TransactOpts, teamId, leagueId, posInLeague)
}

// SignToLeague is a paid mutator transaction binding the contract method 0x02154127.
//
// Solidity: function signToLeague(uint256 teamId, uint256 leagueId, uint8 posInLeague) returns()
func (_Assets *AssetsTransactorSession) SignToLeague(teamId *big.Int, leagueId *big.Int, posInLeague uint8) (*types.Transaction, error) {
	return _Assets.Contract.SignToLeague(&_Assets.TransactOpts, teamId, leagueId, posInLeague)
}

// AssetsTeamCreationIterator is returned from FilterTeamCreation and is used to iterate over the raw logs and unpacked data for TeamCreation events raised by the Assets contract.
type AssetsTeamCreationIterator struct {
	Event *AssetsTeamCreation // Event containing the contract specifics and raw log

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
func (it *AssetsTeamCreationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AssetsTeamCreation)
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
		it.Event = new(AssetsTeamCreation)
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
func (it *AssetsTeamCreationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AssetsTeamCreationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AssetsTeamCreation represents a TeamCreation event raised by the Assets contract.
type AssetsTeamCreation struct {
	TeamName string
	TeamId   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTeamCreation is a free log retrieval operation binding the contract event 0x46d7278a8c77d35e56c1a7bf3699d639520893a20903fe88510aca1341b283a0.
//
// Solidity: event TeamCreation(string teamName, uint256 teamId)
func (_Assets *AssetsFilterer) FilterTeamCreation(opts *bind.FilterOpts) (*AssetsTeamCreationIterator, error) {

	logs, sub, err := _Assets.contract.FilterLogs(opts, "TeamCreation")
	if err != nil {
		return nil, err
	}
	return &AssetsTeamCreationIterator{contract: _Assets.contract, event: "TeamCreation", logs: logs, sub: sub}, nil
}

// WatchTeamCreation is a free log subscription operation binding the contract event 0x46d7278a8c77d35e56c1a7bf3699d639520893a20903fe88510aca1341b283a0.
//
// Solidity: event TeamCreation(string teamName, uint256 teamId)
func (_Assets *AssetsFilterer) WatchTeamCreation(opts *bind.WatchOpts, sink chan<- *AssetsTeamCreation) (event.Subscription, error) {

	logs, sub, err := _Assets.contract.WatchLogs(opts, "TeamCreation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AssetsTeamCreation)
				if err := _Assets.contract.UnpackLog(event, "TeamCreation", log); err != nil {
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
