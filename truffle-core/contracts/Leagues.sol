pragma solidity >=0.4.21 <0.6.0;

import "./Assets.sol";
import "./Engine.sol";
/**
 * @title Scheduling of leagues, and calls to Engine to resolve games.
 */

contract Leagues is Assets {
    
    uint8 constant public MATCHDAYS = 14;
    uint8 constant public MATCHES_PER_DAY = 4;
    Engine private _engine;

    function setEngineAdress(address addr) public {
        _engine = Engine(addr);
    }

    function getEngineAddress() public view returns (address) {
        return address(_engine);
    }

    function getTeamsInMatch(uint8 matchday, uint8 matchIdxInDay) public pure returns (uint8 homeIdx, uint8 visitorIdx) 
    {
        require(matchday < MATCHDAYS, "wrong match day");
        require(matchIdxInDay < MATCHES_PER_DAY, "wrong match");
        if (matchday < (TEAMS_PER_LEAGUE - 1))
            (homeIdx, visitorIdx) = _getTeamsInMatchFirstHalf(matchday, matchIdxInDay);
        else
            (visitorIdx, homeIdx) = _getTeamsInMatchFirstHalf(matchday - (TEAMS_PER_LEAGUE - 1), matchIdxInDay);
    }

    function _shiftBack(uint8 t) private pure returns (uint8)
    {
        if (t < TEAMS_PER_LEAGUE)
            return t;
        else
            return t-(TEAMS_PER_LEAGUE-1);
    }

    function _getTeamsInMatchFirstHalf(uint8 matchday, uint8 matchIdxInDay) private pure returns (uint8, uint8) 
    {
        uint8 team1 = 0;
        if (matchIdxInDay > 0)
            team1 = _shiftBack(TEAMS_PER_LEAGUE-matchIdxInDay+matchday);

        uint8 team2 = _shiftBack(matchIdxInDay+1+matchday);
        if ( (matchday % 2) == 0)
            return (team1, team2);
        else
            return (team2, team1);
    }

    // returns [scoreHome, scoreAway, scoreHome, scoreAway,...]
    // TODO: currentVerseSeed must be provided from getCurrentVerseSeed()
    // TODO: likewise, matchday should be computed outside
    function computeMatchday(
        uint8 matchday,
        uint256[PLAYERS_PER_TEAM_MAX][TEAMS_PER_LEAGUE] memory prevLeagueState,
        uint256[TEAMS_PER_LEAGUE] memory tacticsIds,
        uint256 currentVerseSeed
    )
        public
        view
        returns (uint8[2 * MATCHES_PER_DAY] memory scores)
    {
        uint8[2] memory score;
        uint8 homeTeamIdx;
        uint8 visitorTeamIdx;
        for (uint8 matchIdxInDay = 0; matchIdxInDay < MATCHES_PER_DAY ; matchIdxInDay++)
        {
            (homeTeamIdx, visitorTeamIdx) = getTeamsInMatch(matchday, matchIdxInDay);
            uint256 matchSeed = uint256(keccak256(abi.encode(currentVerseSeed, matchIdxInDay))); 
            uint256[2] memory tactics = [tacticsIds[homeTeamIdx], tacticsIds[visitorTeamIdx]];
            uint256[PLAYERS_PER_TEAM_MAX][2] memory states = [prevLeagueState[homeTeamIdx], prevLeagueState[visitorTeamIdx]];
            score = _engine.playMatch(
                matchSeed, 
                states,
                tactics,
                false,
                false
            );
            scores[matchIdxInDay * 2] = score[0];
            scores[matchIdxInDay * 2 +1 ] = score[1];
        }
    }    
}