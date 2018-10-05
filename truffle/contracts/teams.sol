pragma solidity ^ 0.4.24;

import "./players.sol";

contract TeamFactory is PlayerFactory {

    /// @dev Fired whenever a new team is created
    event TeamCreation(string teamName, uint nCreatedTeams, address owner);

    /// @dev An internal method that creates an empty new team and stores it. This
    /// @dev method doesn't do any checking and should only be called when the
    /// @dev input data is known to be valid.
    /// @param _teamName The name of the team, which happens to determine the team skills, via hash(_teamName, _userChoice)
    function createTeam(string _teamName) internal {
        // TODO: require maxLen for _teamName

        /// @dev Make sure team name did not exist before.
        bytes32 nameHash = keccak256(abi.encodePacked(_teamName));
        require(teamToOwnerAddr[nameHash]==0);

        /// @dev Create empty team and store assigned name.
        Team memory emptyTeam;
        emptyTeam.name = _teamName;

        /// @dev At this stage, playerIdx = 0.
        /// @dev A team is considered as 'created' if the owner has a non-null address.
        teams.push(emptyTeam);
        teamToOwnerAddr[nameHash] = msg.sender;

        // emit the team creation event
        emit TeamCreation(_teamName, teams.length, msg.sender);
    }

    /// @dev Returns the entire state of the player (age, skills, etc.) given his idx in a given team
    function getSkill(uint _teamIdx, uint8 _playerIdx)
        internal
        view
        returns(uint)
    {
        uint playerIdx = getNumAtIndex(teams[_teamIdx].playersIdx, _playerIdx, bitsPerPlayerIdx());
        return players[playerIdx].state;
    }

/* 
    @dev Section with functions only for external/testing use.
*/    
    function getNCreatedTeams() internal view returns(uint) { return teams.length;}
    function getTeamName(uint idx) internal view returns(string) { return teams[idx].name;}


}
