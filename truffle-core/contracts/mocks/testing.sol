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
    constructor(address teamFactory) public
    League(teamFactory)
    {
    }
    // WRAPPERS FOR GAMES

    // function test_playGame(uint teamIdx1, uint teamIdx2, uint seed)
    //     external
    //     returns (uint16[2] memory teamGoals) 
    // {
    //     return playGame(teamIdx1, teamIdx2, seed);
    // }


    // function test_playRound(uint leagueIdx, uint8 round, uint seed) 
    //     external 
    // {
    //     playRound(leagueIdx, round, seed);  
    // }

    // // WRAPPERS FOR LEAGUES

    function test_createLeague(uint[] _teamIdxs, uint _blockFirstGame, uint _blocksBetweenGames) 
        external
        returns (uint) 
    {
        return createLeague(_teamIdxs, _blockFirstGame, _blocksBetweenGames);
    }

    function test_getTeamsIdxsInLeague(uint leagueIdx) external view returns (uint[]) {
        return getTeamsIdxsInLeague(leagueIdx);
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
