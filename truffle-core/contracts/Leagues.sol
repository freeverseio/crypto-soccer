pragma solidity >=0.4.21 <0.6.0;

import "./Assets.sol";

contract Leagues is Assets {
    
    uint8 constant public MATCHDAYS = 14;
    uint8 constant public MATCHES_PER_DAY = 4;
    
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
}