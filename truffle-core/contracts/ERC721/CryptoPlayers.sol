pragma solidity ^0.4.24;

import "./CryptoPlayersMetadata.sol";

contract CryptoPlayers is CryptoPlayersMetadata {
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
