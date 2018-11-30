pragma solidity ^0.4.24;

import "./CryptoPlayersMetadata.sol";

contract CryptoPlayers is CryptoPlayersMetadata {
    function addPlayer(string memory name, uint state, uint256 teamIdx, address owner) public {
        _addPlayer(name, state, teamIdx, owner);
    }

    function getPlayerState(uint playerIdx) public view returns(uint) {
        return _getPlayerState(playerIdx);
    }

    function getNCreatedPlayers() public view returns(uint) { 
        return _getNCreatedPlayers();
    }

    function getPlayerName(uint playerIdx) public view returns(string) {
        return _getPlayerName(playerIdx);
    }

    function getTeamIndexByPlayer(string name) public view returns (uint256){
        return _getTeamIndexByPlayer(name);
    }

    function playerExists(string name) public view returns (bool){
        return _playerExists(name);
    }
}
