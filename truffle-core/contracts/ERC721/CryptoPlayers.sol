pragma solidity ^0.4.24;

import "./CryptoPlayersMetadata.sol";

contract CryptoPlayers is CryptoPlayersMetadata {
    function addPlayer(string memory name, uint state, uint256 teamIdx) public {
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        players.push(Player({name: name, state: state}));
        playerToTeam[playerNameHash] = teamIdx;
    }

    function getPlayerState(uint playerIdx) public view returns(uint) {
        return players[playerIdx + 1].state;
    }

    function getNCreatedPlayers() public view returns(uint) { 
        return players.length - 1;
    }

    function getPlayerName(uint playerIdx) public view returns(string) {
        return players[playerIdx + 1].name;
    }

    function getTeamIndexByPlayer(string name) public view returns (uint256){
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        return playerToTeam[playerNameHash];
    }

    function playerExists(string name) public view returns (bool){
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 teamIdx = playerToTeam[playerNameHash];
        return teamIdx != 0;
    }
}
