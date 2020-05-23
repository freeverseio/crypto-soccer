pragma solidity >= 0.6.3;

/**
* @title Storage required by Proxy
*/
contract ProxyStorage {

    struct ContractInfo {
        address addr;
        bytes4[] selectors;
        bytes32 name;
        bool isActive;
    }

    address internal _company; 
    address internal _proposedCompany;
    address internal _superUser; 
    ContractInfo[] internal _contractsInfo;
    mapping (bytes4 => address) internal _selectorToContractAddr;
  
    modifier onlyCompany() {
        require(msg.sender == _company, "Only company is authorized.");
        _;
    }
    
    modifier onlySuperUser() {
        require(msg.sender == _superUser, "Only superuser is authorized.");
        _;
    }
}