// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testutils

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

// TestutilsABI is the input ABI used to generate the binding from.
const TestutilsABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"increase\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// TestutilsBin is the compiled bytecode used for deploying new contracts.
const TestutilsBin = `608060405234801561001057600080fd5b506000808190555060bd806100266000396000f3fe6080604052348015600f57600080fd5b5060043610604f576000357c01000000000000000000000000000000000000000000000000000000009004806306661abd146054578063e8927fbc146070575b600080fd5b605a6078565b6040518082815260200191505060405180910390f35b6076607e565b005b60005481565b600080815480929190600101919050555056fea165627a7a72305820d666f1482e1bbc1c6626a0cf3f0de57d65af0685e5b1bd686e96758b6c5d15ae0029`

// DeployTestutils deploys a new Ethereum contract, binding an instance of Testutils to it.
func DeployTestutils(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Testutils, error) {
	parsed, err := abi.JSON(strings.NewReader(TestutilsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TestutilsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Testutils{TestutilsCaller: TestutilsCaller{contract: contract}, TestutilsTransactor: TestutilsTransactor{contract: contract}, TestutilsFilterer: TestutilsFilterer{contract: contract}}, nil
}

// Testutils is an auto generated Go binding around an Ethereum contract.
type Testutils struct {
	TestutilsCaller     // Read-only binding to the contract
	TestutilsTransactor // Write-only binding to the contract
	TestutilsFilterer   // Log filterer for contract events
}

// TestutilsCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestutilsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestutilsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestutilsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestutilsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestutilsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestutilsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestutilsSession struct {
	Contract     *Testutils        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestutilsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestutilsCallerSession struct {
	Contract *TestutilsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TestutilsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestutilsTransactorSession struct {
	Contract     *TestutilsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TestutilsRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestutilsRaw struct {
	Contract *Testutils // Generic contract binding to access the raw methods on
}

// TestutilsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestutilsCallerRaw struct {
	Contract *TestutilsCaller // Generic read-only contract binding to access the raw methods on
}

// TestutilsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestutilsTransactorRaw struct {
	Contract *TestutilsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestutils creates a new instance of Testutils, bound to a specific deployed contract.
func NewTestutils(address common.Address, backend bind.ContractBackend) (*Testutils, error) {
	contract, err := bindTestutils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Testutils{TestutilsCaller: TestutilsCaller{contract: contract}, TestutilsTransactor: TestutilsTransactor{contract: contract}, TestutilsFilterer: TestutilsFilterer{contract: contract}}, nil
}

// NewTestutilsCaller creates a new read-only instance of Testutils, bound to a specific deployed contract.
func NewTestutilsCaller(address common.Address, caller bind.ContractCaller) (*TestutilsCaller, error) {
	contract, err := bindTestutils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestutilsCaller{contract: contract}, nil
}

// NewTestutilsTransactor creates a new write-only instance of Testutils, bound to a specific deployed contract.
func NewTestutilsTransactor(address common.Address, transactor bind.ContractTransactor) (*TestutilsTransactor, error) {
	contract, err := bindTestutils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestutilsTransactor{contract: contract}, nil
}

// NewTestutilsFilterer creates a new log filterer instance of Testutils, bound to a specific deployed contract.
func NewTestutilsFilterer(address common.Address, filterer bind.ContractFilterer) (*TestutilsFilterer, error) {
	contract, err := bindTestutils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestutilsFilterer{contract: contract}, nil
}

// bindTestutils binds a generic wrapper to an already deployed contract.
func bindTestutils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestutilsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Testutils *TestutilsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Testutils.Contract.TestutilsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Testutils *TestutilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Testutils.Contract.TestutilsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Testutils *TestutilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Testutils.Contract.TestutilsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Testutils *TestutilsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Testutils.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Testutils *TestutilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Testutils.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Testutils *TestutilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Testutils.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() constant returns(uint256)
func (_Testutils *TestutilsCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Testutils.contract.Call(opts, out, "count")
	return *ret0, err
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() constant returns(uint256)
func (_Testutils *TestutilsSession) Count() (*big.Int, error) {
	return _Testutils.Contract.Count(&_Testutils.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() constant returns(uint256)
func (_Testutils *TestutilsCallerSession) Count() (*big.Int, error) {
	return _Testutils.Contract.Count(&_Testutils.CallOpts)
}

// Increase is a paid mutator transaction binding the contract method 0xe8927fbc.
//
// Solidity: function increase() returns()
func (_Testutils *TestutilsTransactor) Increase(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Testutils.contract.Transact(opts, "increase")
}

// Increase is a paid mutator transaction binding the contract method 0xe8927fbc.
//
// Solidity: function increase() returns()
func (_Testutils *TestutilsSession) Increase() (*types.Transaction, error) {
	return _Testutils.Contract.Increase(&_Testutils.TransactOpts)
}

// Increase is a paid mutator transaction binding the contract method 0xe8927fbc.
//
// Solidity: function increase() returns()
func (_Testutils *TestutilsTransactorSession) Increase() (*types.Transaction, error) {
	return _Testutils.Contract.Increase(&_Testutils.TransactOpts)
}
