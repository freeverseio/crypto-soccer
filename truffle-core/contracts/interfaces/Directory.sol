pragma solidity >= 0.6.3;

import "../storage/UniverseInfo.sol";

/**
 @title Directory for contract addresses related to the project
 @author Freeverse.io, www.freeverse.io
 @dev keeps the deployed contracts in the arrays _names and _addresses
 @dev it keeps two copies of such arrays, and a pointer that points to the "active" copy
 @dev this is to allow the pattern: "first add the new data, which can cost a lot of gas
 @dev ...and, in one atomic transaction, just change the pointer".
 @dev We use this to let the Proxy upgrade this Directory as part of the proxy upgrade, atomically.
 @dev We keep track of deployBlockNum as a safety measure, to ensure that the contracts-to-become-active
 @dev were deployed later thant he currently active contracts. Otherwise, we need to call "revertActivation"
*/

contract Directory {

    event DeployedDirectory(bytes32[] names, address[] addrs, uint8 newActivePtr);
    event Activation(uint8 activePtr);

    UniverseInfo private _universeInfo;

    bytes32[][2] internal _names;
    address[][2] internal _addresses;
    uint256[2] deployBlockNum;
    address owner;
    uint8 activePtr; /// only allowed to be 0 or 1
    
    constructor (address proxyAddr) public {
        owner = proxyAddr;
        _universeInfo = UniverseInfo(proxyAddr);
    }
    
    /// We grant write permission either to the COO or to the Proxy contract itself
    /// so that it can activate the Directory atomically, as part of its upgrade.
    modifier onlyOwners() {
        require(msg.sender == owner || msg.sender == _universeInfo.COO(), "Non authorized attempt to write to Directory");
        _;
    }

    function deploy(bytes32[] calldata newNames, address[] calldata newAdresseses) external onlyOwners {
        uint256 nContr = newNames.length;
        require(nContr == newAdresseses.length, "non-matching number of names and addresses");
        uint8 newActivePtr = 1 - activePtr;
        delete _names[newActivePtr];
        delete _addresses[newActivePtr];
        for (uint256 contr = 0; contr < nContr; contr++) {
           _names[newActivePtr].push(newNames[contr]); 
           _addresses[newActivePtr].push(newAdresseses[contr]); 
        }
        deployBlockNum[newActivePtr] = block.number;
        emit DeployedDirectory(newNames, newAdresseses, newActivePtr);
    }

    function activateNewDeploy() external onlyOwners {
        uint8 currentPtr = activePtr;
        uint8 newPtr = 1 - currentPtr;
        require(deployBlockNum[newPtr] > deployBlockNum[currentPtr], "cannot activate a set of contracts that were deployed before the current contracts");
        activePtr = newPtr;
        emit Activation(newPtr);
    }

    function revertActivation() external onlyOwners  {
        uint8 currentPtr = activePtr;
        uint8 newPtr = 1 - currentPtr;
        require(deployBlockNum[newPtr] < deployBlockNum[currentPtr], "cannot revert to a set of contracts that were deployed after the current contracts");
        activePtr = newPtr;
        emit Activation(newPtr);
    }

    /// returns the active directory data
    function getDirectory() public view returns (bytes32[] memory, address[] memory) {
        uint8 ptr = activePtr;
        uint256 nContr = _names[activePtr].length;
        bytes32[] memory names = new bytes32[](nContr);
        address[] memory addresses = new address[](nContr);
        for (uint256 contr = 0; contr < nContr; contr++) {
            names[contr] = _names[ptr][contr];
            addresses[contr] = _addresses[ptr][contr];
        }
        return (names, addresses);
    }
}