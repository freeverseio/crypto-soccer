pragma solidity >=0.4.21 <0.6.0;

import "./Players.sol";

contract Teams is Players {
    event TeamCreation (string teamName, uint256 teamId);

    constructor(address playerState) public Players(playerState) {
    }

    function createTeam(string memory teamName) public {
        uint256 teamId = _addTeam(teamName);
        emit TeamCreation(teamName, teamId);
    }

    function getPlayerIdFromTeamIdAndPos(uint256 teamId, uint8 posInTeam) external view returns (uint256) {
        require(_teamExists(teamId), "unexistent team");
        require(posInTeam < PLAYERS_PER_TEAM, "invalid player pos");
        return PLAYERS_PER_TEAM * (teamId - 1) + 1 + posInTeam;
    }

    function signToLeague(uint256 teamId, uint256 leagueId, uint8 posInLeague) public {
        require(_teamExists(teamId), "unexistent team");
        _updateTeamCurrentHistory(teamId, leagueId, posInLeague);
    }
}