pragma solidity ^ 0.4.24;
import "./games.sol";

/*
    Main contract to manage Leagues
*/

contract League is GameEngine {

    /// @dev Creates a league
    function createLeague(uint[] _teamIdxs, uint _blockFirstGame, uint _blocksBetweenGames) internal {
        leagues.push(League({teamIdxs: _teamIdxs, blockFirstGame: _blockFirstGame, blocksBetweenGames: _blocksBetweenGames}));
    }

    function getTeamsInLeague(uint leagueIdx) internal view returns (uint[]) {
        return leagues[leagueIdx].teamIdxs;
    }
}