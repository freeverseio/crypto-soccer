pragma solidity >=0.5.12 <=0.6.3;

/**
* @title Storage required by Proxy
*/
contract ProxyStorage {

    address internal _proxyOwner; 
    address internal _proposedProxyOwner;
    ContractInfo[] internal _contractsInfo;
    mapping (bytes4 => uint256) internal _selectorToContractId;

    struct ContractInfo {
        address addr;
        bytes4[] selectors;
        bytes32 name;
    }
}