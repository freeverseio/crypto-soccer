pragma solidity >=0.4.21 <0.6.0;

import "../assets/Players.sol";

contract PlayersMock is Players {
    function addTeam(string memory name) public returns (uint256) {
        _addTeam(name);
    }
}