pragma solidity >=0.4.21 <0.6.0;

import "../assets/Storage.sol";

contract StorageMock is Storage {
    constructor(address playerState)
    public
    Storage(playerState)
    {
    }
    
    function addTeam(string memory name) public returns (uint256) {
        return _addTeam(name);
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
}