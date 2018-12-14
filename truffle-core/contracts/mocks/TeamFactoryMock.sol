pragma solidity ^0.4.24;

import "../factory/teams.sol";

/**
 * @title TeamFactoryMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract TeamFactoryMock is TeamFactory {
    // Wrappers for helpers:
    constructor(address cryptoTeams, address cryptoPlayers) public
    TeamFactory(cryptoTeams, cryptoPlayers){

    }

    function test_teamsInGame(uint8 round, uint8 game, uint8 nTeams) 
        public 
        pure 
        returns (uint8 team1, uint8 team2)
    {
        return teamsInGame(round, game, nTeams);
    }

    function test_serialize(uint8 nElem, uint16[] nums, uint bits) public pure returns(uint result) {
        return serialize(nElem, nums, bits);
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

    function test_computeKeccak256ForNumber(uint n) external pure returns(uint)
    {
        return computeKeccak256ForNumber(n);
    }

    function test_computeKeccak256(string s, uint n1, uint n2) external pure returns(uint) {
        return computeKeccak256(s, n1, n2);
    }

    function test_throwDice(uint weight1, uint weight2, uint rndNum, uint factor) external pure returns(uint8) {
        return throwDice(weight1, weight2, rndNum, factor);
    }

    function test_throwDiceArray(uint[] weights, uint rndNum, uint factor) external pure returns(uint8) {
        return throwDiceArray(weights, rndNum, factor);
    }

    function test_getRndNumArrays(uint seed, uint8 roundsPerGame, uint8 bitsPerRndNum) 
        external
        pure
        returns (uint16[] rndNumArray) 
    {
        return getRndNumArrays(seed, roundsPerGame, bitsPerRndNum);
    }

    function test_getGameId(uint teamIdx1, uint teamIdx2, uint seed) 
        external
        pure
        returns (uint gameId) 
    {
        return getGameId(teamIdx1, teamIdx2, seed);
    }


    // Wrappers for Players

    function test_createBalancedPlayer(
        string _playerName, 
        uint _teamIdx, 
        uint16 _userChoice, 
        uint8 _playerNumberInTeam, 
        uint8 _playerRole,
        address owner
    ) 
        external 
    {
        createBalancedPlayer(
            _playerName, 
            _teamIdx, 
            _userChoice,
            _playerNumberInTeam, 
            _playerRole,
            owner
        );
    }

    function test_createUnbalancedPlayer(
        string _playerName,
        uint _teamIdx,
        uint8 _playerNumberInTeam,
        uint _monthOfBirthAfterUnixEpoch,
        uint _defense,
        uint _speed,
        uint _pass,
        uint _shoot,
        uint _endurance,
        uint _role,
        address owner
    )
        external 
    {
        createUnbalancedPlayer(
            _playerName,
            _teamIdx, 
            _playerNumberInTeam,
            _monthOfBirthAfterUnixEpoch,
            _defense,
            _speed,
            _pass,
            _shoot,
            _endurance,
            _role,
            owner
        );
    }

    function test_getRole(uint idx, uint8 first, uint8 second) external pure returns(uint8) {
        return getRole(idx, first, second);
    }

    function test_getNCreatedPlayers() external view returns(uint) { return _cryptoPlayers.getNCreatedPlayers(); }
    function test_getPlayerState(uint playerIdx) external view returns(uint) { return _cryptoPlayers.getState(playerIdx); }
    function test_getPlayerName(uint playerIdx) external view returns(string) { return _cryptoPlayers.getName(playerIdx); }


    // WRAPPERS FOR TEAMS

    function totalSupply() external view returns(uint) { return _cryptoTeams.totalSupply(); }
    function test_getTeamName(uint idx) public view returns(string) { return _cryptoTeams.getName(idx); }
    function test_createTeam(string _teamName) external { return createTeam(_teamName); }
    function test_getStatePlayerInTeam(uint8 _playerIdx, uint _teamIdx) external view returns(uint) { return getStatePlayerInTeam(_playerIdx, _teamIdx); }
}