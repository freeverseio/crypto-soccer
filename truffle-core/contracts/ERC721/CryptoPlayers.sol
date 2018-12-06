pragma solidity ^0.4.24;

import "./CryptoPlayersMetadata.sol";

contract CryptoPlayers is CryptoPlayersMetadata {
    function addPlayer(string memory name, uint state, uint256 teamIdx, address owner) public {
        uint256 nextPlayerId = totalSupply() + 1;
        mintWithName(owner, nextPlayerId, name);
        _setState(nextPlayerId, state);
        _setTeam(nextPlayerId, teamIdx);
    }

    function getNCreatedPlayers() public view returns(uint) { 
        return totalSupply();
    }

    function playerExists(string name) public view returns (bool){
        uint256 id = getPlayerId(name);
        return id != 0;
    }

    function getTeamIndexByPlayer(string name) public view returns (uint256){
        uint256 id = getPlayerId(name);
        return getTeam(id);
    }
}
