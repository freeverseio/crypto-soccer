pragma solidity ^0.5.0;

/**
 * teamId == 0 is invalid and represents the null team
 */
contract Storage {
    uint8 constant PLAYERS_PER_TEAM = 11;

    /// @dev The player skills in each team are obtained from hashing: name + userChoice
    /// @dev So userChoice allows the user to inspect lots of teams compatible with his chosen name
    /// @dev and select his favourite one.
    /// @dev playerIdx serializes each player idx, allowing 20 bit for each (>1M players possible)
    struct Team {
        string name;
    }

    /// @dev An array containing the Team struct for all teams in existence.
    /// @dev The ID of each team is actually his index in this array.
    Team[] private teams;

    constructor() public {
        teams.push(Team("_"));
    }

    function getPlayersPerTeam() external view returns (uint8) {
        return PLAYERS_PER_TEAM;
    }

    function countTeams() public view returns (uint256){
        return teams.length - 1;
    }

    function getTeamName(uint256 teamId) public view returns (string memory) {
        require(_teamExists(teamId), "invalid team id");
        return teams[teamId].name;
    }

    function _addTeam(string memory name) internal returns (uint256) {
        teams.push(Team(name));
        return teams.length - 1;
    }

    function _teamExists(uint256 teamId) internal view returns (bool) {
        return teamId != 0 && teamId < teams.length;
    }
}
