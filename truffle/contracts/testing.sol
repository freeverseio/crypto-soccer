pragma solidity ^ 0.4.24;

/*
Contract to act as a wrapper for all functions that need to be tested directly.
This avoids declaring those as 'public' or 'external', when they will not be so when deploying

Inheritance structure:
    PlayerFactory is HelperFunctions
    TeamFactory is PlayerFactory
    GameEngine is TeamFactory
    Testing is GameEngine
*/

import "./games.sol";

contract Testing is GameEngine {

    // Wrappers for helpers:
    function test_encode(uint8 nElem, uint16[] nums, uint bits) public pure returns(uint result) {
        return encode(nElem, nums, bits);
    }

    function test_decode(uint8 nNumbers, uint longState, uint bits) public pure returns(uint16[] result) {
        return decode(nNumbers, longState, bits);
    }

    function test_getNumAtIndex(uint longState, uint8 index, uint bits) public pure returns(uint) {
        return getNumAtIndex(longState, index, bits);
    }

    function test_setNumAtIndex(uint value, uint longState, uint8 index, uint bits) public pure returns(uint) {
        return setNumAtIndex(value, longState, index, bits);
    }

    function test_computeKeccak256ForNumber(uint n) public pure returns(uint)
    {
        return computeKeccak256ForNumber(n);
    }

    function test_computeKeccak256(string s, uint n1, uint n2) public pure returns(uint) {
        return computeKeccak256(s, n1, n2);
    }

    function test_throwDice(uint weight1, uint weight2, uint rndNum, uint factor) public pure returns(uint8) {
        return throwDice(weight1, weight2, rndNum, factor);
    }

    function test_throwDiceArray(uint[] memory weights, uint rndNum, uint factor) public pure returns(uint8) {
        return throwDiceArray(weights, rndNum, factor);
    }


    // Wrappers for Players

    function test_createRandomPlayer(string _playerName, uint _teamIdx, uint16 _userChoice, uint8 _playerNumberInTeam, uint8 _playerRole) public {
        return createRandomPlayer(_playerName, _teamIdx, _userChoice, _playerNumberInTeam, _playerRole);
    }

    function test_createPlayer(
        string _playerName,
        uint _teamIdx,
        uint8 _playerNumberInTeam,
        uint _monthOfBirthAfterUnixEpoch,
        uint _defense,
        uint _speed,
        uint _pass,
        uint _shoot,
        uint _endurance,
        uint _role
        )
        public {
        return createPlayer(
            _playerName,
            _teamIdx, 
            _playerNumberInTeam,
            _monthOfBirthAfterUnixEpoch,
            _defense,
            _speed,
            _pass,
            _shoot,
            _endurance,
            _role
        );
    }

    function test_getRole(uint idx, uint8 first, uint8 second) public pure returns(uint8) {
        return  getRole(idx, first, second);
    }

    function test_getNCreatedPlayers() external view returns(uint) { return getNCreatedPlayers(); }
    function test_getPlayerState(uint playerIdx) external view returns(uint) { return getPlayerState(playerIdx); }


    // WRAPPERS FOR TEAMS

    function test_getNCreatedTeams() public view returns(uint) { return getNCreatedTeams(); }
    function test_getTeamName(uint idx) public view returns(string) { return getTeamName(idx); }
    function test_createTeam(string _teamName) public { return createTeam(_teamName); }



}
