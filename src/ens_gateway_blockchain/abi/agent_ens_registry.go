// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// KeyValue is an auto generated low-level Go binding around an user-defined struct.
type KeyValue struct {
	Key   string
	Value string
}

// AgentEnsRegistryMetaData contains all meta data concerning the AgentEnsRegistry contract.
var AgentEnsRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"managerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"parentNode_\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collection\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"AlreadyProvisioned\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"DuplicateKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"MissingRequiredKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collection\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"NotAgentOwnerOrWallet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collection\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"NotProvisioned\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"ReservedKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"StringsInsufficientHexLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StringsInvalidAddressFormat\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collection\",\"type\":\"address\"}],\"name\":\"UnknownCollection\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"parameter\",\"type\":\"string\"}],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"collection\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"AgentTextUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newManager\",\"type\":\"address\"}],\"name\":\"ManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newParentNode\",\"type\":\"bytes32\"}],\"name\":\"ParentNodeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"collection\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"}],\"name\":\"SubnodeProvisioned\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"addr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"agentRef\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"collection\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collection\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"internalType\":\"structKeyValue[]\",\"name\":\"records\",\"type\":\"tuple[]\"}],\"name\":\"initAgentTexts\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"manager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"parentNode\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"provisioned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collection\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"setAgentText\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newManager\",\"type\":\"address\"}],\"name\":\"setManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newParentNode\",\"type\":\"bytes32\"}],\"name\":\"setParentNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"text\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AgentEnsRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentEnsRegistryMetaData.ABI instead.
var AgentEnsRegistryABI = AgentEnsRegistryMetaData.ABI

// AgentEnsRegistry is an auto generated Go binding around an Ethereum contract.
type AgentEnsRegistry struct {
	AgentEnsRegistryCaller     // Read-only binding to the contract
	AgentEnsRegistryTransactor // Write-only binding to the contract
	AgentEnsRegistryFilterer   // Log filterer for contract events
}

// AgentEnsRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgentEnsRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentEnsRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentEnsRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentEnsRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentEnsRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentEnsRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentEnsRegistrySession struct {
	Contract     *AgentEnsRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AgentEnsRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentEnsRegistryCallerSession struct {
	Contract *AgentEnsRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// AgentEnsRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentEnsRegistryTransactorSession struct {
	Contract     *AgentEnsRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// AgentEnsRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgentEnsRegistryRaw struct {
	Contract *AgentEnsRegistry // Generic contract binding to access the raw methods on
}

// AgentEnsRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentEnsRegistryCallerRaw struct {
	Contract *AgentEnsRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// AgentEnsRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentEnsRegistryTransactorRaw struct {
	Contract *AgentEnsRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgentEnsRegistry creates a new instance of AgentEnsRegistry, bound to a specific deployed contract.
func NewAgentEnsRegistry(address common.Address, backend bind.ContractBackend) (*AgentEnsRegistry, error) {
	contract, err := bindAgentEnsRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistry{AgentEnsRegistryCaller: AgentEnsRegistryCaller{contract: contract}, AgentEnsRegistryTransactor: AgentEnsRegistryTransactor{contract: contract}, AgentEnsRegistryFilterer: AgentEnsRegistryFilterer{contract: contract}}, nil
}

// NewAgentEnsRegistryCaller creates a new read-only instance of AgentEnsRegistry, bound to a specific deployed contract.
func NewAgentEnsRegistryCaller(address common.Address, caller bind.ContractCaller) (*AgentEnsRegistryCaller, error) {
	contract, err := bindAgentEnsRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistryCaller{contract: contract}, nil
}

// NewAgentEnsRegistryTransactor creates a new write-only instance of AgentEnsRegistry, bound to a specific deployed contract.
func NewAgentEnsRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AgentEnsRegistryTransactor, error) {
	contract, err := bindAgentEnsRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistryTransactor{contract: contract}, nil
}

// NewAgentEnsRegistryFilterer creates a new log filterer instance of AgentEnsRegistry, bound to a specific deployed contract.
func NewAgentEnsRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AgentEnsRegistryFilterer, error) {
	contract, err := bindAgentEnsRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistryFilterer{contract: contract}, nil
}

// bindAgentEnsRegistry binds a generic wrapper to an already deployed contract.
func bindAgentEnsRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgentEnsRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentEnsRegistry *AgentEnsRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentEnsRegistry.Contract.AgentEnsRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentEnsRegistry *AgentEnsRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.AgentEnsRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentEnsRegistry *AgentEnsRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.AgentEnsRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentEnsRegistry *AgentEnsRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentEnsRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentEnsRegistry *AgentEnsRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentEnsRegistry *AgentEnsRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentEnsRegistry *AgentEnsRegistrySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentEnsRegistry.Contract.DEFAULTADMINROLE(&_AgentEnsRegistry.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentEnsRegistry.Contract.DEFAULTADMINROLE(&_AgentEnsRegistry.CallOpts)
}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
//
// Solidity: function addr(bytes32 node) view returns(address)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) Addr(opts *bind.CallOpts, node [32]byte) (common.Address, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "addr", node)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
//
// Solidity: function addr(bytes32 node) view returns(address)
func (_AgentEnsRegistry *AgentEnsRegistrySession) Addr(node [32]byte) (common.Address, error) {
	return _AgentEnsRegistry.Contract.Addr(&_AgentEnsRegistry.CallOpts, node)
}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
//
// Solidity: function addr(bytes32 node) view returns(address)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) Addr(node [32]byte) (common.Address, error) {
	return _AgentEnsRegistry.Contract.Addr(&_AgentEnsRegistry.CallOpts, node)
}

// AgentRef is a free data retrieval call binding the contract method 0x0dd323f7.
//
// Solidity: function agentRef(bytes32 ) view returns(address collection, uint256 tokenId)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) AgentRef(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Collection common.Address
	TokenId    *big.Int
}, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "agentRef", arg0)

	outstruct := new(struct {
		Collection common.Address
		TokenId    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Collection = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// AgentRef is a free data retrieval call binding the contract method 0x0dd323f7.
//
// Solidity: function agentRef(bytes32 ) view returns(address collection, uint256 tokenId)
func (_AgentEnsRegistry *AgentEnsRegistrySession) AgentRef(arg0 [32]byte) (struct {
	Collection common.Address
	TokenId    *big.Int
}, error) {
	return _AgentEnsRegistry.Contract.AgentRef(&_AgentEnsRegistry.CallOpts, arg0)
}

// AgentRef is a free data retrieval call binding the contract method 0x0dd323f7.
//
// Solidity: function agentRef(bytes32 ) view returns(address collection, uint256 tokenId)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) AgentRef(arg0 [32]byte) (struct {
	Collection common.Address
	TokenId    *big.Int
}, error) {
	return _AgentEnsRegistry.Contract.AgentRef(&_AgentEnsRegistry.CallOpts, arg0)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentEnsRegistry *AgentEnsRegistrySession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgentEnsRegistry.Contract.GetRoleAdmin(&_AgentEnsRegistry.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgentEnsRegistry.Contract.GetRoleAdmin(&_AgentEnsRegistry.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentEnsRegistry *AgentEnsRegistrySession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgentEnsRegistry.Contract.HasRole(&_AgentEnsRegistry.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgentEnsRegistry.Contract.HasRole(&_AgentEnsRegistry.CallOpts, role, account)
}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() view returns(address)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) Manager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "manager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() view returns(address)
func (_AgentEnsRegistry *AgentEnsRegistrySession) Manager() (common.Address, error) {
	return _AgentEnsRegistry.Contract.Manager(&_AgentEnsRegistry.CallOpts)
}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() view returns(address)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) Manager() (common.Address, error) {
	return _AgentEnsRegistry.Contract.Manager(&_AgentEnsRegistry.CallOpts)
}

// ParentNode is a free data retrieval call binding the contract method 0xf3068a00.
//
// Solidity: function parentNode() view returns(bytes32)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) ParentNode(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "parentNode")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ParentNode is a free data retrieval call binding the contract method 0xf3068a00.
//
// Solidity: function parentNode() view returns(bytes32)
func (_AgentEnsRegistry *AgentEnsRegistrySession) ParentNode() ([32]byte, error) {
	return _AgentEnsRegistry.Contract.ParentNode(&_AgentEnsRegistry.CallOpts)
}

// ParentNode is a free data retrieval call binding the contract method 0xf3068a00.
//
// Solidity: function parentNode() view returns(bytes32)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) ParentNode() ([32]byte, error) {
	return _AgentEnsRegistry.Contract.ParentNode(&_AgentEnsRegistry.CallOpts)
}

// Provisioned is a free data retrieval call binding the contract method 0x2d0e9c27.
//
// Solidity: function provisioned(bytes32 ) view returns(bool)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) Provisioned(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "provisioned", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Provisioned is a free data retrieval call binding the contract method 0x2d0e9c27.
//
// Solidity: function provisioned(bytes32 ) view returns(bool)
func (_AgentEnsRegistry *AgentEnsRegistrySession) Provisioned(arg0 [32]byte) (bool, error) {
	return _AgentEnsRegistry.Contract.Provisioned(&_AgentEnsRegistry.CallOpts, arg0)
}

// Provisioned is a free data retrieval call binding the contract method 0x2d0e9c27.
//
// Solidity: function provisioned(bytes32 ) view returns(bool)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) Provisioned(arg0 [32]byte) (bool, error) {
	return _AgentEnsRegistry.Contract.Provisioned(&_AgentEnsRegistry.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentEnsRegistry *AgentEnsRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentEnsRegistry.Contract.SupportsInterface(&_AgentEnsRegistry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentEnsRegistry.Contract.SupportsInterface(&_AgentEnsRegistry.CallOpts, interfaceId)
}

// Text is a free data retrieval call binding the contract method 0x59d1d43c.
//
// Solidity: function text(bytes32 node, string key) view returns(string)
func (_AgentEnsRegistry *AgentEnsRegistryCaller) Text(opts *bind.CallOpts, node [32]byte, key string) (string, error) {
	var out []interface{}
	err := _AgentEnsRegistry.contract.Call(opts, &out, "text", node, key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Text is a free data retrieval call binding the contract method 0x59d1d43c.
//
// Solidity: function text(bytes32 node, string key) view returns(string)
func (_AgentEnsRegistry *AgentEnsRegistrySession) Text(node [32]byte, key string) (string, error) {
	return _AgentEnsRegistry.Contract.Text(&_AgentEnsRegistry.CallOpts, node, key)
}

// Text is a free data retrieval call binding the contract method 0x59d1d43c.
//
// Solidity: function text(bytes32 node, string key) view returns(string)
func (_AgentEnsRegistry *AgentEnsRegistryCallerSession) Text(node [32]byte, key string) (string, error) {
	return _AgentEnsRegistry.Contract.Text(&_AgentEnsRegistry.CallOpts, node, key)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentEnsRegistry *AgentEnsRegistrySession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.GrantRole(&_AgentEnsRegistry.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.GrantRole(&_AgentEnsRegistry.TransactOpts, role, account)
}

// InitAgentTexts is a paid mutator transaction binding the contract method 0x41877a90.
//
// Solidity: function initAgentTexts(address collection, uint256 tokenId, (string,string)[] records) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactor) InitAgentTexts(opts *bind.TransactOpts, collection common.Address, tokenId *big.Int, records []KeyValue) (*types.Transaction, error) {
	return _AgentEnsRegistry.contract.Transact(opts, "initAgentTexts", collection, tokenId, records)
}

// InitAgentTexts is a paid mutator transaction binding the contract method 0x41877a90.
//
// Solidity: function initAgentTexts(address collection, uint256 tokenId, (string,string)[] records) returns()
func (_AgentEnsRegistry *AgentEnsRegistrySession) InitAgentTexts(collection common.Address, tokenId *big.Int, records []KeyValue) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.InitAgentTexts(&_AgentEnsRegistry.TransactOpts, collection, tokenId, records)
}

// InitAgentTexts is a paid mutator transaction binding the contract method 0x41877a90.
//
// Solidity: function initAgentTexts(address collection, uint256 tokenId, (string,string)[] records) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactorSession) InitAgentTexts(collection common.Address, tokenId *big.Int, records []KeyValue) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.InitAgentTexts(&_AgentEnsRegistry.TransactOpts, collection, tokenId, records)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentEnsRegistry *AgentEnsRegistrySession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.RenounceRole(&_AgentEnsRegistry.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.RenounceRole(&_AgentEnsRegistry.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentEnsRegistry *AgentEnsRegistrySession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.RevokeRole(&_AgentEnsRegistry.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.RevokeRole(&_AgentEnsRegistry.TransactOpts, role, account)
}

// SetAgentText is a paid mutator transaction binding the contract method 0xc6947c78.
//
// Solidity: function setAgentText(address collection, uint256 tokenId, string key, string value) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactor) SetAgentText(opts *bind.TransactOpts, collection common.Address, tokenId *big.Int, key string, value string) (*types.Transaction, error) {
	return _AgentEnsRegistry.contract.Transact(opts, "setAgentText", collection, tokenId, key, value)
}

// SetAgentText is a paid mutator transaction binding the contract method 0xc6947c78.
//
// Solidity: function setAgentText(address collection, uint256 tokenId, string key, string value) returns()
func (_AgentEnsRegistry *AgentEnsRegistrySession) SetAgentText(collection common.Address, tokenId *big.Int, key string, value string) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.SetAgentText(&_AgentEnsRegistry.TransactOpts, collection, tokenId, key, value)
}

// SetAgentText is a paid mutator transaction binding the contract method 0xc6947c78.
//
// Solidity: function setAgentText(address collection, uint256 tokenId, string key, string value) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactorSession) SetAgentText(collection common.Address, tokenId *big.Int, key string, value string) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.SetAgentText(&_AgentEnsRegistry.TransactOpts, collection, tokenId, key, value)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address newManager) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactor) SetManager(opts *bind.TransactOpts, newManager common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.contract.Transact(opts, "setManager", newManager)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address newManager) returns()
func (_AgentEnsRegistry *AgentEnsRegistrySession) SetManager(newManager common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.SetManager(&_AgentEnsRegistry.TransactOpts, newManager)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address newManager) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactorSession) SetManager(newManager common.Address) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.SetManager(&_AgentEnsRegistry.TransactOpts, newManager)
}

// SetParentNode is a paid mutator transaction binding the contract method 0x52d32416.
//
// Solidity: function setParentNode(bytes32 newParentNode) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactor) SetParentNode(opts *bind.TransactOpts, newParentNode [32]byte) (*types.Transaction, error) {
	return _AgentEnsRegistry.contract.Transact(opts, "setParentNode", newParentNode)
}

// SetParentNode is a paid mutator transaction binding the contract method 0x52d32416.
//
// Solidity: function setParentNode(bytes32 newParentNode) returns()
func (_AgentEnsRegistry *AgentEnsRegistrySession) SetParentNode(newParentNode [32]byte) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.SetParentNode(&_AgentEnsRegistry.TransactOpts, newParentNode)
}

// SetParentNode is a paid mutator transaction binding the contract method 0x52d32416.
//
// Solidity: function setParentNode(bytes32 newParentNode) returns()
func (_AgentEnsRegistry *AgentEnsRegistryTransactorSession) SetParentNode(newParentNode [32]byte) (*types.Transaction, error) {
	return _AgentEnsRegistry.Contract.SetParentNode(&_AgentEnsRegistry.TransactOpts, newParentNode)
}

// AgentEnsRegistryAgentTextUpdatedIterator is returned from FilterAgentTextUpdated and is used to iterate over the raw logs and unpacked data for AgentTextUpdated events raised by the AgentEnsRegistry contract.
type AgentEnsRegistryAgentTextUpdatedIterator struct {
	Event *AgentEnsRegistryAgentTextUpdated // Event containing the contract specifics and raw log

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
func (it *AgentEnsRegistryAgentTextUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsRegistryAgentTextUpdated)
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
		it.Event = new(AgentEnsRegistryAgentTextUpdated)
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
func (it *AgentEnsRegistryAgentTextUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsRegistryAgentTextUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsRegistryAgentTextUpdated represents a AgentTextUpdated event raised by the AgentEnsRegistry contract.
type AgentEnsRegistryAgentTextUpdated struct {
	Collection common.Address
	TokenId    *big.Int
	Key        string
	Value      string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAgentTextUpdated is a free log retrieval operation binding the contract event 0x02a3afca70f705104b53a1f0081265d9f3214cdd9b3ce3ca14bbb16638ed93b0.
//
// Solidity: event AgentTextUpdated(address indexed collection, uint256 indexed tokenId, string key, string value)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) FilterAgentTextUpdated(opts *bind.FilterOpts, collection []common.Address, tokenId []*big.Int) (*AgentEnsRegistryAgentTextUpdatedIterator, error) {

	var collectionRule []interface{}
	for _, collectionItem := range collection {
		collectionRule = append(collectionRule, collectionItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.FilterLogs(opts, "AgentTextUpdated", collectionRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistryAgentTextUpdatedIterator{contract: _AgentEnsRegistry.contract, event: "AgentTextUpdated", logs: logs, sub: sub}, nil
}

// WatchAgentTextUpdated is a free log subscription operation binding the contract event 0x02a3afca70f705104b53a1f0081265d9f3214cdd9b3ce3ca14bbb16638ed93b0.
//
// Solidity: event AgentTextUpdated(address indexed collection, uint256 indexed tokenId, string key, string value)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) WatchAgentTextUpdated(opts *bind.WatchOpts, sink chan<- *AgentEnsRegistryAgentTextUpdated, collection []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var collectionRule []interface{}
	for _, collectionItem := range collection {
		collectionRule = append(collectionRule, collectionItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.WatchLogs(opts, "AgentTextUpdated", collectionRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsRegistryAgentTextUpdated)
				if err := _AgentEnsRegistry.contract.UnpackLog(event, "AgentTextUpdated", log); err != nil {
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

// ParseAgentTextUpdated is a log parse operation binding the contract event 0x02a3afca70f705104b53a1f0081265d9f3214cdd9b3ce3ca14bbb16638ed93b0.
//
// Solidity: event AgentTextUpdated(address indexed collection, uint256 indexed tokenId, string key, string value)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) ParseAgentTextUpdated(log types.Log) (*AgentEnsRegistryAgentTextUpdated, error) {
	event := new(AgentEnsRegistryAgentTextUpdated)
	if err := _AgentEnsRegistry.contract.UnpackLog(event, "AgentTextUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsRegistryManagerUpdatedIterator is returned from FilterManagerUpdated and is used to iterate over the raw logs and unpacked data for ManagerUpdated events raised by the AgentEnsRegistry contract.
type AgentEnsRegistryManagerUpdatedIterator struct {
	Event *AgentEnsRegistryManagerUpdated // Event containing the contract specifics and raw log

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
func (it *AgentEnsRegistryManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsRegistryManagerUpdated)
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
		it.Event = new(AgentEnsRegistryManagerUpdated)
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
func (it *AgentEnsRegistryManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsRegistryManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsRegistryManagerUpdated represents a ManagerUpdated event raised by the AgentEnsRegistry contract.
type AgentEnsRegistryManagerUpdated struct {
	NewManager common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterManagerUpdated is a free log retrieval operation binding the contract event 0x2c1c11af44aa5608f1dca38c00275c30ea091e02417d36e70e9a1538689c433d.
//
// Solidity: event ManagerUpdated(address indexed newManager)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) FilterManagerUpdated(opts *bind.FilterOpts, newManager []common.Address) (*AgentEnsRegistryManagerUpdatedIterator, error) {

	var newManagerRule []interface{}
	for _, newManagerItem := range newManager {
		newManagerRule = append(newManagerRule, newManagerItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.FilterLogs(opts, "ManagerUpdated", newManagerRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistryManagerUpdatedIterator{contract: _AgentEnsRegistry.contract, event: "ManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchManagerUpdated is a free log subscription operation binding the contract event 0x2c1c11af44aa5608f1dca38c00275c30ea091e02417d36e70e9a1538689c433d.
//
// Solidity: event ManagerUpdated(address indexed newManager)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) WatchManagerUpdated(opts *bind.WatchOpts, sink chan<- *AgentEnsRegistryManagerUpdated, newManager []common.Address) (event.Subscription, error) {

	var newManagerRule []interface{}
	for _, newManagerItem := range newManager {
		newManagerRule = append(newManagerRule, newManagerItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.WatchLogs(opts, "ManagerUpdated", newManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsRegistryManagerUpdated)
				if err := _AgentEnsRegistry.contract.UnpackLog(event, "ManagerUpdated", log); err != nil {
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

// ParseManagerUpdated is a log parse operation binding the contract event 0x2c1c11af44aa5608f1dca38c00275c30ea091e02417d36e70e9a1538689c433d.
//
// Solidity: event ManagerUpdated(address indexed newManager)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) ParseManagerUpdated(log types.Log) (*AgentEnsRegistryManagerUpdated, error) {
	event := new(AgentEnsRegistryManagerUpdated)
	if err := _AgentEnsRegistry.contract.UnpackLog(event, "ManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsRegistryParentNodeUpdatedIterator is returned from FilterParentNodeUpdated and is used to iterate over the raw logs and unpacked data for ParentNodeUpdated events raised by the AgentEnsRegistry contract.
type AgentEnsRegistryParentNodeUpdatedIterator struct {
	Event *AgentEnsRegistryParentNodeUpdated // Event containing the contract specifics and raw log

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
func (it *AgentEnsRegistryParentNodeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsRegistryParentNodeUpdated)
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
		it.Event = new(AgentEnsRegistryParentNodeUpdated)
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
func (it *AgentEnsRegistryParentNodeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsRegistryParentNodeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsRegistryParentNodeUpdated represents a ParentNodeUpdated event raised by the AgentEnsRegistry contract.
type AgentEnsRegistryParentNodeUpdated struct {
	NewParentNode [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterParentNodeUpdated is a free log retrieval operation binding the contract event 0xe6817a10601362205ecdb69e5da819694f3f12c10bcdbf3f6755bac5050ab60a.
//
// Solidity: event ParentNodeUpdated(bytes32 indexed newParentNode)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) FilterParentNodeUpdated(opts *bind.FilterOpts, newParentNode [][32]byte) (*AgentEnsRegistryParentNodeUpdatedIterator, error) {

	var newParentNodeRule []interface{}
	for _, newParentNodeItem := range newParentNode {
		newParentNodeRule = append(newParentNodeRule, newParentNodeItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.FilterLogs(opts, "ParentNodeUpdated", newParentNodeRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistryParentNodeUpdatedIterator{contract: _AgentEnsRegistry.contract, event: "ParentNodeUpdated", logs: logs, sub: sub}, nil
}

// WatchParentNodeUpdated is a free log subscription operation binding the contract event 0xe6817a10601362205ecdb69e5da819694f3f12c10bcdbf3f6755bac5050ab60a.
//
// Solidity: event ParentNodeUpdated(bytes32 indexed newParentNode)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) WatchParentNodeUpdated(opts *bind.WatchOpts, sink chan<- *AgentEnsRegistryParentNodeUpdated, newParentNode [][32]byte) (event.Subscription, error) {

	var newParentNodeRule []interface{}
	for _, newParentNodeItem := range newParentNode {
		newParentNodeRule = append(newParentNodeRule, newParentNodeItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.WatchLogs(opts, "ParentNodeUpdated", newParentNodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsRegistryParentNodeUpdated)
				if err := _AgentEnsRegistry.contract.UnpackLog(event, "ParentNodeUpdated", log); err != nil {
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

// ParseParentNodeUpdated is a log parse operation binding the contract event 0xe6817a10601362205ecdb69e5da819694f3f12c10bcdbf3f6755bac5050ab60a.
//
// Solidity: event ParentNodeUpdated(bytes32 indexed newParentNode)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) ParseParentNodeUpdated(log types.Log) (*AgentEnsRegistryParentNodeUpdated, error) {
	event := new(AgentEnsRegistryParentNodeUpdated)
	if err := _AgentEnsRegistry.contract.UnpackLog(event, "ParentNodeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsRegistryRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the AgentEnsRegistry contract.
type AgentEnsRegistryRoleAdminChangedIterator struct {
	Event *AgentEnsRegistryRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *AgentEnsRegistryRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsRegistryRoleAdminChanged)
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
		it.Event = new(AgentEnsRegistryRoleAdminChanged)
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
func (it *AgentEnsRegistryRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsRegistryRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsRegistryRoleAdminChanged represents a RoleAdminChanged event raised by the AgentEnsRegistry contract.
type AgentEnsRegistryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*AgentEnsRegistryRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistryRoleAdminChangedIterator{contract: _AgentEnsRegistry.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *AgentEnsRegistryRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsRegistryRoleAdminChanged)
				if err := _AgentEnsRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) ParseRoleAdminChanged(log types.Log) (*AgentEnsRegistryRoleAdminChanged, error) {
	event := new(AgentEnsRegistryRoleAdminChanged)
	if err := _AgentEnsRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsRegistryRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the AgentEnsRegistry contract.
type AgentEnsRegistryRoleGrantedIterator struct {
	Event *AgentEnsRegistryRoleGranted // Event containing the contract specifics and raw log

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
func (it *AgentEnsRegistryRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsRegistryRoleGranted)
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
		it.Event = new(AgentEnsRegistryRoleGranted)
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
func (it *AgentEnsRegistryRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsRegistryRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsRegistryRoleGranted represents a RoleGranted event raised by the AgentEnsRegistry contract.
type AgentEnsRegistryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgentEnsRegistryRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistryRoleGrantedIterator{contract: _AgentEnsRegistry.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *AgentEnsRegistryRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsRegistryRoleGranted)
				if err := _AgentEnsRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) ParseRoleGranted(log types.Log) (*AgentEnsRegistryRoleGranted, error) {
	event := new(AgentEnsRegistryRoleGranted)
	if err := _AgentEnsRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsRegistryRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the AgentEnsRegistry contract.
type AgentEnsRegistryRoleRevokedIterator struct {
	Event *AgentEnsRegistryRoleRevoked // Event containing the contract specifics and raw log

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
func (it *AgentEnsRegistryRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsRegistryRoleRevoked)
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
		it.Event = new(AgentEnsRegistryRoleRevoked)
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
func (it *AgentEnsRegistryRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsRegistryRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsRegistryRoleRevoked represents a RoleRevoked event raised by the AgentEnsRegistry contract.
type AgentEnsRegistryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgentEnsRegistryRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistryRoleRevokedIterator{contract: _AgentEnsRegistry.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *AgentEnsRegistryRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsRegistryRoleRevoked)
				if err := _AgentEnsRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) ParseRoleRevoked(log types.Log) (*AgentEnsRegistryRoleRevoked, error) {
	event := new(AgentEnsRegistryRoleRevoked)
	if err := _AgentEnsRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentEnsRegistrySubnodeProvisionedIterator is returned from FilterSubnodeProvisioned and is used to iterate over the raw logs and unpacked data for SubnodeProvisioned events raised by the AgentEnsRegistry contract.
type AgentEnsRegistrySubnodeProvisionedIterator struct {
	Event *AgentEnsRegistrySubnodeProvisioned // Event containing the contract specifics and raw log

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
func (it *AgentEnsRegistrySubnodeProvisionedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentEnsRegistrySubnodeProvisioned)
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
		it.Event = new(AgentEnsRegistrySubnodeProvisioned)
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
func (it *AgentEnsRegistrySubnodeProvisionedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentEnsRegistrySubnodeProvisionedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentEnsRegistrySubnodeProvisioned represents a SubnodeProvisioned event raised by the AgentEnsRegistry contract.
type AgentEnsRegistrySubnodeProvisioned struct {
	Collection common.Address
	TokenId    *big.Int
	Node       [32]byte
	Label      string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSubnodeProvisioned is a free log retrieval operation binding the contract event 0x43e2b90d5c5d23c397f43f751926e1a36e9b6386154990b2a4969d7291d7b3d7.
//
// Solidity: event SubnodeProvisioned(address indexed collection, uint256 indexed tokenId, bytes32 indexed node, string label)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) FilterSubnodeProvisioned(opts *bind.FilterOpts, collection []common.Address, tokenId []*big.Int, node [][32]byte) (*AgentEnsRegistrySubnodeProvisionedIterator, error) {

	var collectionRule []interface{}
	for _, collectionItem := range collection {
		collectionRule = append(collectionRule, collectionItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.FilterLogs(opts, "SubnodeProvisioned", collectionRule, tokenIdRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return &AgentEnsRegistrySubnodeProvisionedIterator{contract: _AgentEnsRegistry.contract, event: "SubnodeProvisioned", logs: logs, sub: sub}, nil
}

// WatchSubnodeProvisioned is a free log subscription operation binding the contract event 0x43e2b90d5c5d23c397f43f751926e1a36e9b6386154990b2a4969d7291d7b3d7.
//
// Solidity: event SubnodeProvisioned(address indexed collection, uint256 indexed tokenId, bytes32 indexed node, string label)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) WatchSubnodeProvisioned(opts *bind.WatchOpts, sink chan<- *AgentEnsRegistrySubnodeProvisioned, collection []common.Address, tokenId []*big.Int, node [][32]byte) (event.Subscription, error) {

	var collectionRule []interface{}
	for _, collectionItem := range collection {
		collectionRule = append(collectionRule, collectionItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _AgentEnsRegistry.contract.WatchLogs(opts, "SubnodeProvisioned", collectionRule, tokenIdRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentEnsRegistrySubnodeProvisioned)
				if err := _AgentEnsRegistry.contract.UnpackLog(event, "SubnodeProvisioned", log); err != nil {
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

// ParseSubnodeProvisioned is a log parse operation binding the contract event 0x43e2b90d5c5d23c397f43f751926e1a36e9b6386154990b2a4969d7291d7b3d7.
//
// Solidity: event SubnodeProvisioned(address indexed collection, uint256 indexed tokenId, bytes32 indexed node, string label)
func (_AgentEnsRegistry *AgentEnsRegistryFilterer) ParseSubnodeProvisioned(log types.Log) (*AgentEnsRegistrySubnodeProvisioned, error) {
	event := new(AgentEnsRegistrySubnodeProvisioned)
	if err := _AgentEnsRegistry.contract.UnpackLog(event, "SubnodeProvisioned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// --- RawLog() methods added by go generate ---
func (e AgentEnsRegistryAgentTextUpdated) RawLog() types.Log   { return e.Raw }
func (e AgentEnsRegistryManagerUpdated) RawLog() types.Log     { return e.Raw }
func (e AgentEnsRegistryParentNodeUpdated) RawLog() types.Log  { return e.Raw }
func (e AgentEnsRegistryRoleAdminChanged) RawLog() types.Log   { return e.Raw }
func (e AgentEnsRegistryRoleGranted) RawLog() types.Log        { return e.Raw }
func (e AgentEnsRegistryRoleRevoked) RawLog() types.Log        { return e.Raw }
func (e AgentEnsRegistrySubnodeProvisioned) RawLog() types.Log { return e.Raw }
