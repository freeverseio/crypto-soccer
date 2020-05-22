pragma solidity >=0.5.12 <=0.6.3;

import "./ProxyStorage.sol";

/**
* @title Manages the state variables of a DelegateProxy
*/
contract Proxy is ProxyStorage {

    event ContractAdded(uint256 contractId, bytes32 name, bytes4[] selectors);
    event ContractsActivated(uint256[] contractIds, uint256 time);
    event ContractsDeactivated(uint256[] contractIds, uint256 time);

    address constant private NULL_ADDR  = address(0);
    address constant private PROXY_DUMMY_ADDR = address(1);

    // TODO: is this future-proof? shall we have it re-settable?
    uint256 constant private FWD_GAS_LIMIT = 10000; 

    /**
    * @dev Sets CompanyOwner and SuperUser
    * @dev Stores proxy selectors in _contractsInfo[0], pointing to PROXY_DUMMY_ADDR
    */
    constructor(address companyOwner, address superUser, bytes4[] memory proxySelectors) public {
        _superUser = msg.sender;
        _contractsInfo.push(ContractInfo(PROXY_DUMMY_ADDR, proxySelectors, "Proxy", false));
        activateContracts(new uint256[](1)); 
        _company = companyOwner;
        _superUser = superUser;
    }
    
    /**
    * @dev execute a delegate call via fallback function
    */
    fallback () external {
        address contractAddr = _selectorToContractAddr[msg.sig];
        require(contractAddr != NULL_ADDR, "function selector is not assigned to a valid contract");
        delegate(contractAddr, msg.data);
    } 
    
    /**
    * @dev Delegates call. It returns the entire context execution
    * @dev NOTE: does not check if the implementation (code) address is a contract,
    *      so having an incorrect implementation could lead to unexpected results
    * @param _target Target address to perform the delegatecall
    * @param _calldata Calldata for the delegatecall
    */
    function delegate(address _target, bytes memory _calldata) internal {
        uint256 fwdGasLimit = FWD_GAS_LIMIT;
        assembly {
            let result := delegatecall(sub(gas(), fwdGasLimit), _target, add(_calldata, 0x20), mload(_calldata), 0, 0)
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
    * @dev Proposes a new owners, who need to later accept it
    */
    function proposeCompany(address addr) public onlyCompany {
        _proposedCompany = addr;
    }

    /**
    * @dev The proposed owners can call these functions to become the owners
    */
    function acceptCompany() public  {
        require(msg.sender == _proposedCompany, "only proposed owner can become owner");
        _company = _proposedCompany;
        _proposedCompany = address(0);
    }

    // SuperUser manages the proxy contract. No need to propose/accept, since it can be changed by company.
    function setSuperUser(address addr) public onlyCompany {
        _superUser = addr;
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
    function addContract(uint256 contractId, address addr, bytes4[] memory selectors, bytes32 name) public onlySuperUser {
        // we require that the contract gets assigned an Id that is as specified from outside, 
        // to make deployment more predictable, and avoid having to parse the emitted event to get contractId:
        require(contractId == _contractsInfo.length, "trying to add a new contract to a contractId that is non-consecutive");
        assertPointsToContract(addr);
        ContractInfo memory info;
        info.addr = addr;
        info.name = name;
        info.isActive = false;
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
    function deactivateAndActivateContracts(uint256[] memory deactContractIds, uint256[] memory actContractIds) public onlySuperUser {
        deactivateContracts(deactContractIds);
        activateContracts(actContractIds);
    }
        
    /**
    * @dev  Activates a set of contracts, by adding an entry in the 
    *       _selectorToContractAddr mapping for each selector of the contract. 
    * @param contractIds The ids of the contracts to be activated
    */
    function activateContracts(uint256[] memory contractIds) public onlySuperUser {
        for (uint256 c = 0; c < contractIds.length; c++) {
            uint256 contractId = contractIds[c];
            require(!_contractsInfo[contractId].isActive, "cannot activate a contract that is already Active");
            bytes4[] memory selectors = _contractsInfo[contractId].selectors;
            address addr = _contractsInfo[contractId].addr;
            for (uint256 s = 0; s < selectors.length; s++) {
                require(_selectorToContractAddr[selectors[s]] != PROXY_DUMMY_ADDR, "Found a collision with a function in the Proxy contract");
                _selectorToContractAddr[selectors[s]] = addr;
            }
            _contractsInfo[contractId].isActive = true;
        }
        emit ContractsActivated(contractIds, now);        
    }

    /**
    * @dev  De-activates a set of contracts, by adding an entry in the 
    *       _selectorToContractAddr mapping for each selector of the contract. 
    * @param contractIds The ids of the contracts to be activated
    */
    function deactivateContracts(uint256[] memory contractIds) public onlySuperUser {
        for (uint256 c = 0; c < contractIds.length; c++) {
            uint256 contractId = contractIds[c];
            require(contractId != 0, "cannot deactivate the proxy contract, with id = 0");
            require(_contractsInfo[contractId].isActive, "cannot deactivate a contract that is Active");
            bytes4[] memory selectors = _contractsInfo[contractId].selectors;
            for (uint256 s = 0; s < selectors.length; s++) {
                delete _selectorToContractAddr[selectors[s]];
            }
            _contractsInfo[contractId].isActive = false;
        }
        emit ContractsDeactivated(contractIds, now);        
    }


   /**
    * @dev Reverts unless contractAddress points to a legit contract.
    *      Makes sure that the hash of the external code is neither 0x0 (not-yet created),
    *       nor an account without code: keccak256('') = 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470
    *      See EIP-1052 for more info
    *      This check is important to avoid delegateCall returning OK when delegating to nowhere
    */
    function assertPointsToContract(address contractAddress) internal view {
        bytes32 emptyContractHash = 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470;
        bytes32 codeHashAtContractAddress;
        assembly { codeHashAtContractAddress := extcodehash(contractAddress) }
        require(codeHashAtContractAddress != emptyContractHash && codeHashAtContractAddress != 0x0, "pointer to a non Contract found!");
    }


    /**
    * @dev  Standard getters
    */
    function countContracts() external view returns(uint256) { return _contractsInfo.length; }
    function countAddressesInContract(uint256 contractId) external view returns(uint256) { return _contractsInfo[contractId].selectors.length; }
    function getContractAddressForSelector(bytes4 selector) public view returns(address) { 
        return _selectorToContractAddr[selector]; 
    }
    function getContractInfo(uint256 contractId) external view returns (address, bytes32, bytes4[] memory, bool) {
        return (
            _contractsInfo[contractId].addr,
            _contractsInfo[contractId].name,
            _contractsInfo[contractId].selectors,
            _contractsInfo[contractId].isActive
        );
    }

}