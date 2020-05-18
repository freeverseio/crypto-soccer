pragma solidity >=0.5.12 <=0.6.3;

/**
* @title Storage required by Proxy
*/
contract ProxyStorage {

    address internal _companyOwner; 
    address internal _proposedCompanyOwner;
    address internal _superUser; 
    
    ContractInfo[] internal _contractsInfo;
    mapping (bytes4 => address) internal _selectorToContractAddr;

    struct ContractInfo {
        address addr;
        bytes4[] selectors;
        bytes32 name;
        bool isActive;
    }
}