pragma solidity ^0.4.24;

import "../factories/teams.sol";

/**
 * @title TeamFactoryMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract TeamFactoryMock is TeamFactory {
    function createTeam_(string _teamName) public {
        super.createTeam(_teamName);
    }

    function getStatePlayerInTeam_(uint8 _playerIdx, uint _teamIdx)
        public
        view
        returns(uint)
    {
        return super.getStatePlayerInTeam(_playerIdx, _teamIdx);
    }

    function getNCreatedTeams_() public view returns(uint) {
        return super.getNCreatedTeams();
    }

    function getTeamName_(uint idx) public view returns(string) { 
        return super.getTeamName(idx);
    }
}