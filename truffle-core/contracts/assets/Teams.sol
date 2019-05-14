pragma solidity >=0.4.21 <0.6.0;

import "./Players.sol";

contract Teams is Players {
    event TeamCreation (string teamName, uint256 teamId);

    constructor(address playerState) public Players(playerState) {
    }

    // TODO: name of the function carries information stored in the name of the params
    // TODO: getPlayerId(uint256 teamId, uint8 posInTeam) already gives all the info
    function getPlayerIdFromTeamIdAndPos(uint256 teamId, uint8 posInTeam) external view returns (uint256) {
        require(_teamExists(teamId), "unexistent team");
        require(posInTeam < PLAYERS_PER_TEAM, "invalid player pos");
        return PLAYERS_PER_TEAM * (teamId - 1) + 1 + posInTeam;
    }

    function createTeam(string memory teamName) public {
        uint256 teamId = _addTeam(teamName);
        emit TeamCreation(teamName, teamId);
    }

    function signToLeague(uint256 teamId, uint256 leagueId, uint8 posInLeague) public {
        require(_teamExists(teamId), "unexistent team");
        // TODO: looking to the usage I think _signToLeague fits more:
        // TODO: What happen inside that function stays inside that function
        _signToLeague(teamId, leagueId, posInLeague);
    }
}