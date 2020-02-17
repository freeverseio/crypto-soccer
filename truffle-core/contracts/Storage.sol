pragma solidity ^0.6.0;
// pragma solidity >=0.5.12 <0.6.2;

/**
* @title Manages the state variables of a DelegateProxy
*/
contract DelegateProxySlotStorage {

    // MAGIC=keccac256("freeverse.proxy.rnd")
    bytes32 constant private MAGIC = 0x5f91a51f585f2d4491bce3e4c2d81799aa0dfc271c36675dbd936650723b29b9;
    uint256 constant private FWD_GAS_LIMIT = 10000;

    uint256[2**10] _slotReserve;
    // TODO: move to a "proposed new owner" + "accept" instead of stright "set net owner"
    address private _storageOwner;
    bytes4[] private _allFunctions;
    uint256[] private _allContracts;
    mapping (bytes32 => uint256) private _functionToContractId;
    mapping (uint256 => ContractInfo) private _contractIdToInfo;

    // ContractInfo: address & name
    //      It's good to store the name of the contract to keep track (and query from outside)
    //      about what they are supposed to fulfil. Examples: name = "Market".
    struct ContractInfo {
        address addr;
        string name;
    }

    constructor() public {
        _storageOwner = msg.sender;
    }
    
    modifier onlyOwner() 
    {
        require(msg.sender == _storageOwner, "Only owner is authorized.");
        _;
    }


    /**
    * @dev execute a delegate call via fallback function
    */
    fallback () external payable {
        address contractAddress = _contractIdToInfo[_functionToContractId[msg.sig]].addr;
        require(contractAddress != address(0), "function is non-declared or not assigned to a valid contract");
        delegate(contractAddress, msg.data);
    } 
    
    /**
    * @dev Performs a delegatecall and returns whatever the delegatecall returned
    *      (entire context execution will return!)
    * @dev NOTE: does not check if the implementation (code) address is a contract,
    *      so having an incorrect implementation could lead to unexpected results
    * @param _dst Destination address to perform the delegatecall
    * @param _calldata Calldata for the delegatecall
    */
    function delegate(address _dst, bytes memory _calldata) internal {
        uint256 fwdGasLimit = FWD_GAS_LIMIT;
        assembly {
            let result := delegatecall(sub(gas(), fwdGasLimit), _dst, add(_calldata, 0x20), mload(_calldata), 0, 0)
            let size := returndatasize()
            let ptr := mload(0x40)
            returndatacopy(ptr, 0, size)

            // revert instead of invalid() bc if the underlying call failed with invalid() it already wasted gas.
            // if the call returned error data, forward it
            switch result case 0 { revert(ptr, size) }
            default { return(ptr, size) }
        }
    }
    
    function setStorageOwner(address newOwner) public onlyOwner {
        _storageOwner = newOwner;
    }
    
    function addNewFunction(bytes4 selector, uint256 contractId) public {
        require(_functionToContractId[selector] == 0, "function with same selector already exists");
        require(_contractExists(contractId), "function cannot point to a non-specified contract");
        _functionToContractId[selector] = contractId;
        _allFunctions.push(selector);
    }

    function deleteFunction(bytes4 selector) public {
        _functionToContractId[selector] = 0;
    }
    
    function addNewContract(uint256 contractId, address addr, string memory name) public {
        require(!_contractExists(contractId), "contractId already exists");
        require(!_stringIsEmpty(name), "cannot create a contract without name");
        require(addr != address(0), "cannot create a contract with null address");
        ContractInfo memory info;
        info.addr = addr;
        info.name = name;
        _contractIdToInfo[contractId] = info;
        _allContracts.push(contractId);
    }

    function changeContractAddr(uint256 contractId, address addr) public {
        require(addr != address(0), "cannot set a contract with null address");
        _contractIdToInfo[contractId].addr = addr;
    }

    function changeContractName(uint256 contractId, string memory name) public {
        require(!_stringIsEmpty(name), "cannot set a contract without name");
        _contractIdToInfo[contractId].name = name;
    }
    
    function deleteContract(uint256 contractId) public {
        delete _contractIdToInfo[contractId];
    }

    function getContractInfo(uint256 contractId) public view returns (address, string memory) {
        return (
            _contractIdToInfo[contractId].addr,
            _contractIdToInfo[contractId].name
        );
    }

    function _contractExists(uint256 contractId) internal view returns (bool) {
        return _contractIdToInfo[contractId].addr != address(0);
    }
    
    function _stringIsEmpty(string memory str) internal pure returns (bool) {
        return bytes(str).length == 0;
    }

  


}