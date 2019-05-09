pragma solidity ^0.5.0;

contract Storage {
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

    function countTeams() external view returns (uint256){
        return teams.length;
    }

    function _addTeam(string memory name) internal returns (uint256) {
        teams.push(Team(name));
        return teams.length;
    }

    // function _getTeamName();

    // function _getPlayersIds();
}
