pragma solidity >=0.4.21 <0.6.0;

import "../assets/Players.sol";

contract PlayersMock is Players {
    constructor(address playerState)
    public
    Players(playerState)
    {
    }

    function addTeam(string memory name, address owner) public returns (uint256) {
        return _addTeam(name, owner);
    }

    
}