pragma solidity >=0.4.21 <0.6.0;

import "../assets/Players.sol";

contract PlayersMock is Players {
    function addTeam(string memory name) public returns (uint256) {
        return _addTeam(name);
    }

    function computeSkills(uint256 seed) public pure returns (uint16[5] memory) {
        return _computeSkills(seed);
    }

    function intHash(string memory arg) public pure returns (uint256) { 
        return _intHash(arg);
    }
}