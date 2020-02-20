pragma solidity >=0.5.12 <0.6.2;

import "./Storage.sol";

/**
* @title Manages the state variables of a DelegateProxy
*/
contract StorageProxy is Storage {

    event AddContract(uint256 contactId, string name);

    uint256 constant private FWD_GAS_LIMIT = 10000; 

    constructor() public {
        _storageOwner = msg.sender;
        _contractIdToInfo.push(ContractInfo(address(0),"Dummy"));
    }
    
    modifier onlyOwner() 
    {
        require(msg.sender == _storageOwner, "Only owner is authorized.");
        _;
    }

    /**
    * @dev execute a delegate call via fallback function
    */
    function () external {
        address contractAddress = _contractIdToInfo[_selectorToContractId[msg.sig]].addr;
        require(contractAddress != address(0), "function selector is non-declared or not assigned to a valid contract");
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
    
    function addNewSelectors(bytes4[] memory selector, uint256 contractId) public onlyOwner {
        require(_contractExists(contractId), "selector cannot point to a non-specified contract");
        for (uint256 s = 0; s < selector.length; s++) {
            // If selector has never been declared before, add to _allSelectors to keep track.
            if(_selectorToContractId[selector[s]] == 0) { _allSelectors.push(selector[s]); }
            _selectorToContractId[selector[s]] = contractId;
           
        }
    }

    function deleteSelectors(bytes4[] memory selectors) public onlyOwner {
        for (uint256 s = 0; s < selectors.length; s++) {
            _selectorToContractId[selectors[s]] = 0;
        }
    }
    
    function addNewContract(address addr, string memory name) public onlyOwner {
        require(!_stringIsEmpty(name), "cannot create a contract without name");
        require(addr != address(0), "cannot create a contract with null address");
        ContractInfo memory info;
        info.addr = addr;
        info.name = name;
        _contractIdToInfo.push(info);
        emit AddContract(_contractIdToInfo.length - 1, name);
    }

    function changeContractAddr(uint256 contractId, address addr) public onlyOwner {
        require(addr != address(0), "cannot set a contract with null address");
        _contractIdToInfo[contractId].addr = addr;
    }

    function changeContractName(uint256 contractId, string memory name) public onlyOwner {
        require(!_stringIsEmpty(name), "cannot set a contract without name");
        _contractIdToInfo[contractId].name = name;
    }
    
    function deleteContract(uint256 contractId) public onlyOwner {
        delete _contractIdToInfo[contractId];
    }

    function _contractExists(uint256 contractId) internal view returns (bool) {
        return _contractIdToInfo[contractId].addr != address(0);
    }
    
    function _stringIsEmpty(string memory str) internal pure returns (bool) {
        return bytes(str).length == 0;
    }

    function countFunctions() external view returns(uint256) { return _allSelectors.length; }
    function countContracts() external view returns(uint256) { return _contractIdToInfo.length; }
    function getContractIdForFunction(bytes4 selector) external view returns(uint256) { return _selectorToContractId[selector]; }
    function getContractInfo(uint256 contractId) public view returns (address, string memory) {
        return (
            _contractIdToInfo[contractId].addr,
            _contractIdToInfo[contractId].name
        );
    }

}