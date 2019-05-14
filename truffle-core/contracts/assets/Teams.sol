pragma solidity >=0.4.21 <0.6.0;

import "./Players.sol";

contract Teams is Players {
    event TeamCreation (string teamName, uint256 teamId);

    constructor(address playerState) public Players(playerState) {
    }

    function createTeam(string memory teamName, address owner) public {
        uint256 teamId = _addTeam(teamName, owner);
        emit TeamCreation(teamName, teamId);
    }

    function signToLeague(uint256 teamId, uint256 leagueId, uint8 posInLeague) public {
        require(_teamExists(teamId), "unexistent team");
        // TODO: looking to the usage I think _signToLeague fits more:
        // TODO: What happen inside that function stays inside that function
        _signToLeague(teamId, leagueId, posInLeague);
    }
}