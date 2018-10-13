pragma solidity ^ 0.4.24;
import "./games.sol";

/*
    Main contract to manage Leagues
*/

contract League is GameEngine {

    /// @dev Creates a league and returns the new league idx
    function createLeague(uint[] _teamIdxs, uint _blockFirstGame, uint _blocksBetweenGames) 
        internal 
        returns (uint)
    {
        leagues.push(League({
            teamIdxs: _teamIdxs, 
            blockFirstGame: _blockFirstGame, 
            blocksBetweenGames: _blocksBetweenGames,
            resultsFirstHalf : 0,
            resultsSecondHalf : 0
            })
        );
        return leagues.length-1;
    }

    /// @dev plays all games in a given round. 
    ///  For a league with nTeams, there are nTeams-1 games in a round.
    function playRound(uint leagueIdx, uint8 round, uint seed)
        internal
        view 
    {
        uint[] memory teamsIdxs = getTeamsIdxsInLeague(leagueIdx);
        uint8 nTeams = uint8(teamsIdxs.length);
        uint8 homeTeam;
        uint8 awayTeam;
        for (uint8 game = 0; game < nTeams-1; game++) {
            (homeTeam, awayTeam) = teamsInGame(round, game, nTeams);
            playGame(
                teamsIdxs[homeTeam],
                teamsIdxs[awayTeam],
                seed+game
            );
        }
    }

    function getTeamsIdxsInLeague(uint leagueIdx) internal view returns (uint[]) {
        return leagues[leagueIdx].teamIdxs;
    }

    function getNTeamsInLeague(uint leagueIdx) internal view returns (uint) {
        return leagues[leagueIdx].teamIdxs.length;
    }

    function getNLeaguesCreated() internal view returns(uint) {
        return leagues.length;
    }
}