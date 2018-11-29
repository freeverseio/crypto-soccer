pragma solidity ^ 0.4.24;

import "./players.sol";

contract TeamFactory is PlayerFactory {
    constructor(address cryptoTeams, address cryptoPlayers) public 
    PlayerFactory(cryptoTeams, cryptoPlayers){

    }

    /// @dev Fired whenever a new team is created
    event TeamCreation(string teamName, uint nCreatedTeams, address owner);

    /// @dev An internal method that creates an empty new team and stores it. This
    /// @dev method doesn't do any checking and should only be called when the
    /// @dev input data is known to be valid.
    /// @param _teamName The name of the team, which happens to determine the team skills, via hash(_teamName, _userChoice)
    function createTeam(string _teamName) public {
        // TODO: require maxLen for _teamName

        /// @dev At this stage, playerIdx = 0.
        /// @dev A team is considered as 'created' if the owner has a non-null address.
        _cryptoTeams.addTeam(_teamName, msg.sender);

        // emit the team creation event
        emit TeamCreation(_teamName, _cryptoTeams.totalSupply(), msg.sender);
    }

    /// @dev Returns the entire state of the player (age, skills, etc.) given his idx in a given team
    function getStatePlayerInTeam(uint8 _playerIdx, uint _teamIdx)
        public
        view
        returns(uint)
    {
        uint playerIdx = getNumAtIndex(_cryptoTeams.getTeamPlayersIdx(_teamIdx), _playerIdx, kBitsPerPlayerIdx);
        return _cryptoPlayers.getPlayerState(playerIdx);
    }

/* 
    @dev Section with functions only for external/testing use.
*/    


    function getTeamState(uint256 team) public view returns(uint256[kMaxPlayersInTeam]){
        uint256[kMaxPlayersInTeam] memory teamState;
        for (uint8 p = 0; p < kMaxPlayersInTeam; p++) {
            teamState[p] = getStatePlayerInTeam(p, team);
        }
        return teamState;
    }
}
