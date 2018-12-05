pragma solidity ^0.4.24;

import "./CryptoPlayersMetadata.sol";

contract CryptoPlayers is CryptoPlayersMetadata {
    function addPlayer(string memory name, uint state, uint256 teamIdx, address owner) public {
        uint256 nextPlayerId = totalSupply() + 1;
        mintWithName(owner, nextPlayerId, name);
        _setState(nextPlayerId, state);
        _setTeam(nextPlayerId, teamIdx);
    }

    function getPlayerState(uint playerIdx) public view returns(uint) {
        return _getState(playerIdx);
    }

    function getNCreatedPlayers() public view returns(uint) { 
        return totalSupply();
    }

    function getPlayerName(uint playerIdx) public view returns(string) {
        return _getName(playerIdx);
    }

    function getTeamIndexByPlayer(string name) public view returns (uint256){
        return _getTeamIndexByPlayer(name);
    }

    function playerExists(string name) public view returns (bool){
        uint256 id = _getPlayer(name);
        return id != 0;
    }
}
