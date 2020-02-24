pragma solidity >=0.5.12 <=0.6.3;

import "./Constants.sol";

/**
* @title Storage required by Proxy
*/
contract ProxyStorage {

    address internal _storageOwner; // TODO: move to a "proposed new owner" + "accept" instead of stright "set net owner"
    ContractInfo[] internal _contractsInfo;
    mapping (bytes4 => uint256) internal _selectorToContractId;

    struct ContractInfo {
        address addr;
        bytes4[] selectors;
        bytes32 name;
    }
}