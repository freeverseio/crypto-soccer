pragma solidity >= 0.6.3;

/**
* @title Directory for contract addresses related to the project
* @dev it keeps the current deployed contracts in the array bytes32[] _names
* @dev      ...and their addresses in the mapping _name2addr
* @dev eve ry time we invoke the deploy function, it deletes the _names array
* @dev      ...and writes the new mapping entries (some will be over-written, some will remain)
*/
contract Directory {

    event DeployedDirectory(bytes32[] names, address[] adresseses);

    mapping (bytes32 => address) internal _name2addr;
    bytes32[] internal _names;
    address _owner;
    
    constructor () public {
        _owner = msg.sender;
    }
    
    function setOwner(address newOwner) public { 
        require(_owner == msg.sender, "Only the owner can set a new owner");
        _owner = newOwner;        
    }

    function deploy(bytes32[] memory newNames, address[] memory newAdresseses) public {
        require(_owner == msg.sender, "Only the owner can deploy new directory");
        uint256 nContr = newNames.length;
        require(nContr == newAdresseses.length, "non-matching number of names and addresses");
        delete _names;
        for (uint256 contr = 0; contr < nContr; contr++) {
           _names.push(newNames[contr]); 
           _name2addr[newNames[contr]] = newAdresseses[contr]; 
        }
        emit DeployedDirectory(newNames, newAdresseses);
    }

    function getAddress(bytes32 name) public view returns (address) {
        return _name2addr[name];
    }
    
    function getDirectory() public view returns (bytes32[] memory, address[] memory) {
        uint256 nContr = _names.length;
        bytes32[] memory names = new bytes32[](nContr);
        address[] memory addresses = new address[](nContr);
        for (uint256 contr = 0; contr < nContr; contr++) {
            names[contr] = _names[contr];
            addresses[contr] = _name2addr[_names[contr]];
        }
        return (names, addresses);
    }

}