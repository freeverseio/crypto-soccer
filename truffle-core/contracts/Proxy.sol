pragma solidity >=0.5.12 <=0.6.3;

import "./ProxyStorage.sol";

/**
* @title Manages the state variables of a DelegateProxy
*/
contract Proxy is ProxyStorage {

    event ContractAdded(uint256 contactId, bytes32 name, bytes4[] selectors);
    event ContractsActivated(uint256[] contactIds);
    event ContractsDeleted(uint256[] contactIds);

    address constant private NULL_ADDR = address(0);

    // TODO: is this future-proof? shall we have it re-settable?
    uint256 constant private FWD_GAS_LIMIT = 10000; 

    /**
    * @dev Sets owner of proxy to whoever deployed it
    * @dev And sets _contractsInfo[0] as the NULL contract
    */
    constructor() public {
        _proxyOwner = msg.sender;
        _contractsInfo.push(ContractInfo(NULL_ADDR, new bytes4[](0), "")); 
    }
    
    modifier onlyOwner() 
    {
        require(msg.sender == _proxyOwner, "Only owner is authorized.");
        _;
    }

    /**
    * @dev execute a delegate call via fallback function
    */
    fallback () external {
        address contractAddr = getContractAddressForSelector(msg.sig);
        require(contractAddr != NULL_ADDR, "function selector is not assigned to a valid contract");
        delegate(contractAddr, msg.data);
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
    
    /**
    * @dev Proposes a new proxy owner, who needs to later accept it
    */
    function proposeProxyOwner(address proposedOwner) public onlyOwner {
        _proposedProxyOwner = proposedOwner;
    }

    /**
    * @dev The proposed owner can call this function to become the owner
    */
    function acceptProxyOwner() public  {
        require(msg.sender == _proposedProxyOwner, "only proposed owner can become owner");
        _proxyOwner = _proposedProxyOwner;
        _proposedProxyOwner = address(0);
    }

    /**
    * @dev Stores the info about a contract to be later called via delegate call,
    * @dev by pushing to the _contractsInfo array, and emits an event with all the info.
    * @dev NOTE: it does not activate it until "activateContracts" is invoked
    * @param contractId The index in the array _contractsInfo where this contract should be placed
    *   It must be equal to the next available idx in the array. Although not strictly necessary, 
    *   it allows the external caller to ensure that the idx is as expected without parsing the event.
    * @param addr Address of the contract that will be used in the delegate call
    * @param selectors An array of all selectors needed inside the contract
    * @param name The name of the added contract, only for reference
    */
    function addContract(uint256 contractId, address addr, bytes4[] memory selectors, bytes32 name) public onlyOwner {
        // we require that the contract gets assigned an Id that is as specified from outside, 
        // to make deployment more predictable, and avoid having to parse the emitted event to get contractId:
        require(contractId == _contractsInfo.length, "trying to add a new contract to a contractId that is non-consecutive");
        ContractInfo memory info;
        info.addr = addr;
        info.name = name;
        info.selectors = selectors;
        _contractsInfo.push(info);
        emit ContractAdded(contractId, name, selectors);        
    }
    
    /**
    * @dev  Deactivates a set of contracts, and then activates another set,
    *       in one single atomic transaction. 
    *       Note: it only removes the mapped selectors, not the contract info. 
    * @param deactContractIds The ids of the contracts to be de-activated
    * @param actContractIds The ids of the contracts to be activated
    */
    function deactivateAndActivateContracts(uint256[] memory deactContractIds, uint256[] memory actContractIds) public onlyOwner {
        deactivateContracts(deactContractIds);
        activateContracts(actContractIds);
    }
        
    /**
    * @dev  Activates a set of contracts, by adding an entry in the 
    *       _selectorToContractId mapping for each selector of the contract. 
    * @param contractIds The ids of the contracts to be activated
    */
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

    /**
    * @dev  De-activates a set of contracts, by adding an entry in the 
    *       _selectorToContractId mapping for each selector of the contract. 
    * @param contractIds The ids of the contracts to be activated
    */
    function deactivateContracts(uint256[] memory contractIds) public onlyOwner {
        for (uint256 c = 0; c < contractIds.length; c++) {
            uint256 contractId = contractIds[c];
            require(contractId != 0, "cannot deactivate the null contract, with id = 0");
            bytes4[] memory selectors = _contractsInfo[contractId].selectors;
            for (uint256 s = 0; s < selectors.length; s++) {
                delete _selectorToContractId[selectors[s]];
            }
        }
        emit ContractsDeleted(contractIds);        
    }


    /**
    * @dev  Standard getters
    */
    function countContracts() external view returns(uint256) { return _contractsInfo.length; }
    function countAddressesInContract(uint256 contractId) external view returns(uint256) { return _contractsInfo[contractId].selectors.length; }
    function getContractAddressForSelector(bytes4 selector) public view returns(address) { 
        return _contractsInfo[_selectorToContractId[selector]].addr; 
    }
    function getContractInfo(uint256 contractId) external view returns (address, bytes32, bytes4[] memory) {
        return (
            _contractsInfo[contractId].addr,
            _contractsInfo[contractId].name,
            _contractsInfo[contractId].selectors
        );
    }

}