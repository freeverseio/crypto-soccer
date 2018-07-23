pragma solidity ^ 0.4.24;

import "./player_factory.sol";

contract TeamFactory is PlayerFactory {

    /// @dev Fired whenever a new team is created
    event TeamCreation(string teamName, uint nCreatedTeams, address owner);

    /// @dev An internal method that creates a new team and stores it. This
    ///  method doesn't do any checking and should only be called when the
    ///  input data is known to be valid.
    /// @param _teamName The name of the team, which happens to determine the team skills, via hash(_teamName, _userChoice)
    function createTeam(string _teamName) public {
        // TODO: require maxLen for _teamName
        bytes32 nameHash = keccak256(abi.encodePacked(_teamName));
        require(teamToOwnerAddr[nameHash]==0);

        Team memory emptyTeam;
        emptyTeam.name = _teamName;
        // At this stage, all playerIdx = 0.
        // We will signal that a team has been created by editing the first player's idx.
        // This enables to require that two players can't have the same name, via:
        //      require(playerToTeam[playerNameHash].playerIdx[0] == 0);
        emptyTeam.playersIdx = uint(0);
        emptyTeam.timeOfCreation = now;

        teams.push(emptyTeam);
        teamToOwnerAddr[nameHash] = msg.sender;

        // emit the team creation event
        emit TeamCreation(_teamName, teams.length, msg.sender);
    }

    function getNCreatedTeams() external view returns(uint) { return teams.length;}
    function getTeamName(uint idx) external view returns(string) { return teams[idx].name;}

    function getSkill(uint _teamIdx, uint8 _playerIdx)
        external
        view
        returns(uint)
    {
        require (_teamIdx < teams.length);
        require (_playerIdx < kMaxPlayersInTeam);
        Team storage t = teams[_teamIdx];
        uint playerIdx = getNumAtIndex(t.playersIdx, _playerIdx, 20);
        return (playerIdx != 0) ? players[playerIdx].state : getDefaultPlayerState(t, _playerIdx);

        // TODO (eaylon): the rest is not useful for any contract, move to the test that needs it:
        //for (uint8 p=0;p<kMaxPlayersInTeam;p++) {
        //    uint16[] memory thisSkills = readNumbersFromUint(7, state, 10000);
        //    for (uint8 sk=0;sk<7;sk++) {
        //        skills[p][sk] = thisSkills[sk];
        //    }
        //}
        //return skills;
        // skills2[0][0]=uint(4);
        // skills2[1][2]=uint(114);
        // return skills2;
    }
}
