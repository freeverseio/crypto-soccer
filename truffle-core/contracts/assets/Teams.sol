pragma solidity ^0.5.0;

import "./Storage.sol";

contract Teams is Storage {
    event TeamCreation (string teamName, uint256 teamId);

    function createTeam(string memory teamName) public {
        uint256 teamId = _addTeam(teamName);
        emit TeamCreation(teamName, teamId);
    }

    function getPlayerIdFromTeamIdAndPos(uint256 teamId, uint8 posInTeam) external view returns (uint256) {
        require(_teamExists(teamId), "unexistent team");
    }
}