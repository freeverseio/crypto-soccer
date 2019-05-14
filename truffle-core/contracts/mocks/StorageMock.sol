pragma solidity >=0.4.21 <0.6.0;

import "../assets/Storage.sol";

contract StorageMock is Storage {
    constructor(address playerState)
    public
    Storage(playerState)
    {
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

    function addTeam(string memory name, address owner) public returns (uint256) {
        return _addTeam(name, owner);
    }

    function playerExists(uint256 playerId) public view returns (bool) {
        return _playerExists(playerId);
    }

    function isVirtual(uint256 playerId) public view returns (bool) {
        return _isVirtual(playerId);
    }

    function setPlayerState(uint256 state) public {
        _setPlayerState(state);
    }

    function teamExists(uint256 teamId) public view returns (bool){
        return _teamExists(teamId);
    }
}