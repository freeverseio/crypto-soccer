pragma solidity ^ 0.4.24;
import "./games.sol";

/*
    Main contract to manage Leagues
*/

contract League is GameEngine {
    constructor(address teamFactory) public
    GameEngine(teamFactory)
    {
    }

    /// @dev The main League struct
    struct LeagueProps {
        uint[] teamIdxs;
        uint blockFirstGame;
        uint blocksBetweenGames;
        uint resultsFirstHalf;
        uint resultsSecondHalf;
    }
    /// @dev Array containing all leagues created so far
    LeagueProps[] leagues;

    /// @dev Creates a league and returns the new league idx
    function createLeague(uint[] _teamIdxs, uint _blockFirstGame, uint _blocksBetweenGames) 
        internal 
        returns (uint)
    {
        leagues.push(LeagueProps({
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
    {
        uint[] memory teamIdxs = leagues[leagueIdx].teamIdxs;
        uint8 nTeams = uint8(teamIdxs.length);
        uint8 homeTeam;
        uint8 awayTeam;
        for (uint8 game = 0; game < nTeams-1; game++) {
            (homeTeam, awayTeam) = _teamFactory.teamsInGame(round, game, nTeams);
            writeGameResult(
                leagueIdx,
                nTeams,
                round,
                game,
                playGame(
                    teamIdxs[homeTeam],
                    teamIdxs[awayTeam],
                    seed+game
                )
            );
        }
    }

    /// @dev writes the result of a game in the serialized uint256
    ///  since there are nTeams/2 games per round, and 2*(nTeams-1) rounds,
    ///  then pos(r,g) = r*nTeams/2 + g
    ///  The first half of these are written in one uint, the second, in another
    function  writeGameResult(
        uint leagueIdx, 
        uint8 nTeams, 
        uint8 round, 
        uint8 game,
        uint16[2] result
    ) 
        internal 
    {
        uint8 result2write = kTie;
        if (result[0] > result[1]) { result2write = kHomeWins; }
        if (result[0] < result[1]) { result2write = kAwayWins; }

        if (round < nTeams - 1) {
            leagues[leagueIdx].resultsFirstHalf = _teamFactory.setNumAtIndex(
                result2write, 
                leagues[leagueIdx].resultsFirstHalf, 
                round*nTeams+game, 
                kBitsPerGameResult
            );
        } else {
            leagues[leagueIdx].resultsSecondHalf = _teamFactory.setNumAtIndex(
                result2write, 
                leagues[leagueIdx].resultsSecondHalf, 
                (round-(nTeams-1))*nTeams+game, 
                kBitsPerGameResult
            );
        }
    }

    function getWrittenResult(uint leagueIdx, uint8 nTeams, uint8 round, uint8 game) 
        internal 
        view 
        returns(uint) 
    {
        if (round < nTeams - 1) {
            return _teamFactory.getNumAtIndex(
                leagues[leagueIdx].resultsFirstHalf, 
                round*nTeams+game, 
                kBitsPerGameResult
            );
        } else {
            return _teamFactory.getNumAtIndex(
                leagues[leagueIdx].resultsSecondHalf, 
                (round-(nTeams-1))*nTeams+game, 
                kBitsPerGameResult
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