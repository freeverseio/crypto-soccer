pragma solidity >=0.5.12 <=0.6.3;

import "./Storage.sol";

/**
* @title Manages the state variables of a DelegateProxy
*/
contract StorageProxy is Storage {

    event ContractAdded(uint256 contactId, bool requiresPermission, bytes32 name, bytes4[] selectors);
    event ContractsActivated(uint256[] contactIds);
    event ContractsDeleted(uint256[] contactIds);

    // TODO: is this future-proof? shall we have it re-settable?
    uint256 constant private FWD_GAS_LIMIT = 10000; 

    constructor() public {
        _storageOwner = msg.sender;
        // _contractsInfo[0] is the NULL contract:
        _contractsInfo.push(ContractInfo(address(0), false, new bytes4[](0), "")); 
    }
    
    modifier onlyOwner() 
    {
        require(msg.sender == _storageOwner, "Only owner is authorized.");
        _;
    }

    /**
    * @dev execute a delegate call via fallback function
    */
    fallback () external {
        ContractInfo memory info = _contractsInfo[_selectorToContractId[msg.sig]];
        require(info.selectors.length != 0, "function selector is not assigned to a valid contract");
        address contractAddress = info.addr;
        if (info.requiresPermission) {
            require(msg.sender == _storageOwner, "Only owner is authorized for this selector.");
        }
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

    function addContract(address addr, bool requiresPermission, bytes4[] memory selectors, bytes32 name) public onlyOwner {
            ContractInfo memory info;
            info.addr = addr;
            info.name = name;
            info.requiresPermission = requiresPermission;
            info.selectors = selectors;
            uint256 contractId = _contractsInfo.length;
            _contractsInfo.push(info);
            emit ContractAdded(contractId, requiresPermission, name, selectors);        
    }
    
    function deleteAndActivateContracts(uint256[] memory deactContractIds, uint256[] memory actContractIds) public onlyOwner {
        deleteContracts(deactContractIds);
        activateContracts(actContractIds);
    }
        
    function activateContracts(uint256[] memory contractIds) public onlyOwner {
        for (uint256 c = 0; c < contractIds.length; c++) {
            uint256 contractId = contractIds[c];
            bytes4[] memory selectors = _contractsInfo[contractId].selectors;
            for (uint256 s = 0; s < selectors.length; s++) {
                _selectorToContractId[selectors[s]] = contractId;
            }
        }
        emit ContractsActivated(contractIds);        
    }

    function deleteContracts(uint256[] memory contractIds) public onlyOwner {
        for (uint256 c = 0; c < contractIds.length; c++) {
            uint256 contractId = contractIds[c];
            bytes4[] memory selectors = _contractsInfo[contractId].selectors;
            for (uint256 s = 0; s < selectors.length; s++) {
                delete _selectorToContractId[selectors[s]];
            }
            delete _contractsInfo[contractId];
        }
        emit ContractsDeleted(contractIds);        
    }

    // function deletePreviousDeploy() private {
    //     // First delete all entries in the selectors mapping _selectorToContractAddr:
    //     for (uint256 contractId = 0; contractId < _contractsInfo.length; contractId++) {
    //         deleteAllSelectorsInContract(contractId);
    //     }
    //     // Finally delete the arrays:
    //     delete _contractsInfo;
    // }
    
    // function deleteAllSelectorsInContract(uint256 contractId) private {
    //     for (uint256 s = 0; s < _contractsInfo[contractId].selectors.length; s++) {
    //         delete _selectorToContractId[_contractsInfo[contractId].selectors[s]];
    //     }
    // }

    // function addSelectors(bytes4[] memory selectors, uint256 contractId) public onlyOwner {
    //     require(_contractExists(contractId), "selectors cannot point to a non-specified contract");
    //     for (uint256 s = 0; s < selectors.length; s++) {
    //         // If selectors has never been declared before, add to _allSelectorsInContract to keep track.
    //         if(_selectorToContractAddr[selectors[s]] == 0) { 
    //             _allSelectorsInContract[contractId].push(selectors[s]); 
    //         }
    //         _selectorToContractAddr[selectors[s]] = contractId;
           
    //     }
    // }


    // function deleteSelectors(bytes4[] memory selectors) public onlyOwner {
    //     for (uint256 s = 0; s < selectors.length; s++) {
    //         _selectorToContractAddr[selectors[s]] = 0;
    //     }
    // }
    
    // function setContract(uint256 contractId, address addr, bool isSetter, string memory name, bytes4[] memory selectors) public onlyOwner {
    //     require(contractId != 0, "contractId = 0 is reserved for NULL contract");
    //     require(!_stringIsEmpty(name), "cannot create a contract without name");
    //     require(addr != address(0), "cannot create a contract with null address");
    //     bool contractExists = contractId < _contractIdToInfo.length;
    //     require(contractExists || contractId == _contractIdToInfo.length, "contractId does not exist, and it is neither the next Id available");

    //     ContractInfo memory info;
    //     info.addr = addr;
    //     info.name = name;
    //     info.isSetter = isSetter;
    
    //     if (contractExists) { 
    //         deleteContract(contractId);
    //         _contractIdToInfo[contractId] = info;
    //     } else {
    //         _contractIdToInfo.push(info);
    //     }
    //     addSelectors(selectors, contractId);
    //     require(false,"--");
    //     emit ContractSet(_contractIdToInfo.length - 1, isSetter, name, selectors);
    // }

    // function changeContractAddr(uint256 contractId, address addr) public onlyOwner {
    //     require(addr != address(0), "cannot set a contract with null address");
    //     _contractIdToInfo[contractId].addr = addr;
    // }

    // function changeContractName(uint256 contractId, string memory name) public onlyOwner {
    //     require(!_stringIsEmpty(name), "cannot set a contract without name");
    //     _contractIdToInfo[contractId].name = name;
    // }

    // function changeContractIsSetter(uint256 contractId, bool isSetter) public onlyOwner {
    //     _contractIdToInfo[contractId].isSetter = isSetter;
    // }
    
    // function deleteContract(uint256 contractId) public onlyOwner {
    //     delete _contractIdToInfo[contractId];
    //     for (uint256 s = 0; s < _allSelectorsInContract[contractId].length; s++) {
    //         delete _selectorToContractAddr[_allSelectorsInContract[contractId][s]];
    //     }
    // }

    // function _contractExists(uint256 contractId) internal view returns (bool) {
    //     return _contractIdToInfo[contractId].addr != address(0);
    // }
    
    function _stringIsEmpty(string memory str) internal pure returns (bool) {
        return bytes(str).length == 0;
    }

    // function countFunctions() external view returns(uint256) { return _allSelectorsInContract.length; }
    function countContracts() external view returns(uint256) { return _contractsInfo.length; }
    function countAddressesInContract(uint256 contractId) external view returns(uint256) { return _contractsInfo[contractId].selectors.length; }
    function getContractAddressForSelector(bytes4 selector) external view returns(address) { 
        return _contractsInfo[_selectorToContractId[selector]].addr; 
    }
    function getContractInfo(uint256 contractId) public view returns (address, bool, bytes32, bytes4[] memory) {
        return (
            _contractsInfo[contractId].addr,
            _contractsInfo[contractId].requiresPermission,
            _contractsInfo[contractId].name,
            _contractsInfo[contractId].selectors
        );
    }

}