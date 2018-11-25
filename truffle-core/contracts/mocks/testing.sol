pragma solidity ^ 0.4.24;

import "../leagues.sol";

/*
    Contract that acts as a wrapper for all functions that need to be tested directly.
    This avoids declaring those as 'public' or 'external', when they will not be so when deploying

    Inheritance structure:
        PlayerFactory is HelperFunctions
        TeamFactory is PlayerFactory
        GameEngine is TeamFactory
        League is GameEngine
        Testing is League
*/

contract Testing is League {

    // Wrappers for helpers:
    function test_serialize(uint8 nElem, uint16[] nums, uint bits) external returns(uint result) {
        return _teamFactory.serialize(nElem, nums, bits);
    }

    function test_decode(uint8 nNumbers, uint longState, uint bits) external returns(uint16[] result) {
        return _teamFactory.decode(nNumbers, longState, bits);
    }

    function test_getNumAtIndex(uint longState, uint8 index, uint bits) external returns(uint) {
        return _teamFactory.getNumAtIndex(longState, index, bits);
    }

    function test_setNumAtIndex(uint value, uint longState, uint8 index, uint bits) external returns(uint) {
        return _teamFactory.setNumAtIndex(value, longState, index, bits);
    }

    function test_computeKeccak256ForNumber(uint n) external returns(uint)
    {
        return _teamFactory.computeKeccak256ForNumber(n);
    }

    function test_computeKeccak256(string s, uint n1, uint n2) external returns(uint) {
        return _teamFactory.computeKeccak256(s, n1, n2);
    }

    function test_throwDice(uint weight1, uint weight2, uint rndNum, uint factor) external returns(uint8) {
        return _teamFactory.throwDice(weight1, weight2, rndNum, factor);
    }

    function test_throwDiceArray(uint[] weights, uint rndNum, uint factor) external returns(uint8) {
        return _teamFactory.throwDiceArray(weights, rndNum, factor);
    }

    function test_getRndNumArrays(uint seed, uint8 roundsPerGame, uint8 bitsPerRndNum) 
        external
        returns (uint16[] rndNumArray) 
    {
        return _teamFactory.getRndNumArrays(seed, roundsPerGame, bitsPerRndNum);
    }

    function test_getGameId(uint teamIdx1, uint teamIdx2, uint seed) 
        external
        returns (uint gameId) 
    {
        return _teamFactory.getGameId(teamIdx1, teamIdx2, seed);
    }


    // Wrappers for Players

    function test_createBalancedPlayer(
        string _playerName, 
        uint _teamIdx, 
        uint16 _userChoice, 
        uint8 _playerNumberInTeam, 
        uint8 _playerRole
    ) 
        external 
    {
        _teamFactory.createBalancedPlayer(
            _playerName, 
            _teamIdx, 
            _userChoice,
            _playerNumberInTeam, 
            _playerRole
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
        uint _role
    )
        external 
    {
        _teamFactory.createUnbalancedPlayer(
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

    function test_getRole(uint idx, uint8 first, uint8 second) external returns(uint8) {
        return _teamFactory.getRole(idx, first, second);
    }

    function test_getNCreatedPlayers() external view returns(uint) { return _teamFactory.getNCreatedPlayers(); }
    function test_getPlayerState(uint playerIdx) external view returns(uint) { return _teamFactory.getPlayerState(playerIdx); }
    function test_getPlayerName(uint playerIdx) external view returns(string) { return _teamFactory.getPlayerName(playerIdx); }


    // WRAPPERS FOR TEAMS

    function test_getNCreatedTeams() external view returns(uint) { return _teamFactory.getNCreatedTeams(); }
    function test_getTeamName(uint idx) external view returns(string) { return _teamFactory.getTeamName(idx); }
    function test_createTeam(string _teamName) external { return _teamFactory.createTeam(_teamName); }
    function test_getStatePlayerInTeam(uint8 _playerIdx, uint _teamIdx) external view returns(uint) { return _teamFactory.getStatePlayerInTeam(_playerIdx, _teamIdx); }


    // WRAPPERS FOR GAMES

    function test_playGame(uint teamIdx1, uint teamIdx2, uint seed)
        external
        returns (uint16[2] memory teamGoals) 
    {
        return playGame(teamIdx1, teamIdx2, seed);
    }

    // WRAPPERS FOR LEAGUES

    function test_createLeague(uint[] _teamIdxs, uint _blockFirstGame, uint _blocksBetweenGames) 
        external
        returns (uint) 
    {
        return createLeague(_teamIdxs, _blockFirstGame, _blocksBetweenGames);
    }

    function test_getTeamsIdxsInLeague(uint leagueIdx) external view returns (uint[]) {
        return getTeamsIdxsInLeague(leagueIdx);
    }

    function test_teamsInGame(uint8 round, uint8 game, uint8 nTeams) 
        external 
        returns (uint8 team1, uint8 team2)
    {
        return _teamFactory.teamsInGame(round, game, nTeams);
    }

    function test_playRound(uint leagueIdx, uint8 round, uint seed) 
        external 
    {
        playRound(leagueIdx, round, seed);  
    }

    function test_getNLeaguesCreated() external view returns(uint) {
        return getNLeaguesCreated();
    }

    function test_getNTeamsInLeague(uint leagueIdx) external view returns(uint) {
        return getNTeamsInLeague(leagueIdx);
    }    

    function test_getWrittenResult(uint leagueIdx, uint8 nTeams, uint8 round, uint8 game)
        external
        view 
        returns(uint)
    {
        return getWrittenResult(leagueIdx, nTeams, round, game);
    }

}
