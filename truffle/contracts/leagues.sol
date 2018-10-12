pragma solidity ^ 0.4.24;
import "./games.sol";

/*
    Main contract to manage Leagues
*/

contract League is GameEngine {

    /// @dev Creates a league
    function createLeague(uint[] _teamIdxs, uint _blockFirstGame, uint _blocksBetweenGames) 
        internal 
        {
        leagues.push(League({
            teamIdxs: _teamIdxs, 
            blockFirstGame: _blockFirstGame, 
            blocksBetweenGames: _blocksBetweenGames,
            resultsFirstHalf : 0,
            resultsSecondHalf : 0
            })
        );
    }

    function getTeamsIdxsInLeague(uint leagueIdx) internal view returns (uint[]) {
        return leagues[leagueIdx].teamIdxs;
    }

    function getNTeamsInLeague(uint leagueIdx) internal view returns (uint) {
        return leagues[leagueIdx].teamIdxs.length;
    }

/*
    function playRound(uint leagueIdx, uint8 round) internal {
        League memory thisLeague = leagues[leagueIdx];
        uint nTeams = getNTeamsInLeague(leagueIdx);
        uint16[2] result; 
        for (uint8 game=0; game < nTeams-1; game++) {
            result = playGame();
        }  
    }
    */
}