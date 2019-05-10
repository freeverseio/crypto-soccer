pragma solidity ^0.5.0;

import "./Storage.sol";

contract Players is Storage {
    /// this function uses the inverse of the following formula
    /// playerId = playersPerTeam * (teamId -1) + 1 + posInTeam;
    function getPlayerTeam(uint256 playerId) external view returns (uint256) {
        require(playerId != 0, "invalid player id");
        uint256 teamId = 1 + (playerId - 1) / PLAYERS_PER_TEAM; 
        require(teamId <= countTeams(), "playerId not created");
        return teamId;
    }
}