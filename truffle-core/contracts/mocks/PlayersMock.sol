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

    function computeSkills(uint256 rnd) public pure returns (uint16[5] memory) {
        return _computeSkills(rnd);
    }

    function intHash(string memory arg) public pure returns (uint256) {
        return _intHash(arg);
    }

    function computeBirth(uint256 rnd) public pure returns (uint16) {
        uint256 currentTime = 1557495456;
        return _computeBirth(rnd, currentTime);
    }

    function setPlayerState(uint256 state) public {
        _setPlayerState(state);
    }
}