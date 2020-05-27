pragma solidity >= 0.6.3;

/**
 @title Storage required by Proxy
 @author Freeverse.io, www.freeverse.io
*/

contract ProxyStorage {

    /// Proxy delegates to various contracts
    /// Contracts are first added, then activated/deactivated
    struct ContractInfo {
        address addr;
        bytes4[] selectors;
        bytes32 name;
        bool isActive;
    }
    ContractInfo[] internal _contractsInfo;
    mapping (bytes4 => address) internal _selectorToContractAddr;

    /// Roles
    address internal _company; 
    address internal _proposedCompany;
    address internal _superUser; 
    address internal _directory; 

    modifier onlyCompany() {
        require(msg.sender == _company, "Only company is authorized.");
        _;
    }
    
    modifier onlySuperUser() {
        require(msg.sender == _superUser, "Only superuser is authorized.");
        _;
    }
}