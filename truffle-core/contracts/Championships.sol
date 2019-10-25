pragma solidity >=0.4.21 <0.6.0;

import "./Engine.sol";
import "./SortIdxs.sol";
/**
 * @title Scheduling of leagues, and calls to Engine to resolve games.
 */

contract Championships is SortIdxs {
    
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 constant public TEAMS_PER_LEAGUE = 8;
    uint8 constant public MATCHDAYS = 14;
    uint8 constant public MATCHES_PER_DAY = 4;
    uint8 constant public MATCHES_PER_LEAGUE = 56; // = 4 * 14 = 7*8
    Engine private _engine;

    function setEngineAdress(address addr) public {
        _engine = Engine(addr);
    }

    function getEngineAddress() public view returns (address) {
        return address(_engine);
    }

    // groupIdx = 0,...,15
    // posInGroup = 0, ...7
    // teamIdx  = 0,...,128
    function getTeamIdxInCup(uint8 groupIdx, uint8 posInGroup) public pure returns(uint8) {
        if (groupIdx % 2 == 0) {
                return 8 * posInGroup + groupIdx / 2;
        } else {
                return 8 * posInGroup + (groupIdx - 1) / 2 + 64;
        }
    }

    function getGroupAndPosInGroup(uint8 teamIdxInCup) public pure returns(uint8 groupIdx, uint8 posInGroup) {
        if (teamIdxInCup < 64) {
            posInGroup = teamIdxInCup / 8;
            groupIdx = 2 * (teamIdxInCup % 8);
        } else {
            posInGroup = (teamIdxInCup-64) / 8;
            groupIdx = 1 + 2 * (teamIdxInCup % 8);
        }
    }

    // sortedTeamIdxInCup contains 64 teams, made up from the top 4 in each of the 16 leagues.
    // they are flattened by groupIdx, and then by final classifications in that group.
    // [groupIdx0_1st, ..., groupIdx0_4th; groupIdx1_1st, ...,]
    // so the index of the array is  groupIdx * 4 + classInGroup
    //      M(2m) 	= L(m mod M, 0) vs L(m+1 mod M, 3), 	m = 0,..., M-1,  M = 16
    //      M(2m+1) = L(m+2 mod M, 1) vs L(m+3 mod M, 2),	m = 0,..., M-1,  M = 16
    // returns indices in sortedTeamIdxInCup
    function getTeamsInCupPlayoffMatch(uint8 matchIdxInDay) public pure returns (uint8 team0, uint8 team1) {
        require(matchIdxInDay < 32, "there are only 32 mathches on day 9 of a cup");
        if (matchIdxInDay % 2 == 0) {
            team0 = (matchIdxInDay/2 % 16) * 4;
            team1 = ((matchIdxInDay/2 + 1) % 16) * 4 + 3;
        } else {
            team0 = (((matchIdxInDay-1)/2 + 2) % 16) * 4 + 1;
            team1 = (((matchIdxInDay-1)/2 + 3) % 16) * 4 + 2;
        }
    }
    
    // same as above, but returns teamIdxInCup as correctly provided by sortedTeamIdxInCup
    function getTeamsInCupPlayoffMatch(uint8 matchIdxInDay, uint8[64] memory sortedTeamIdxInCup) public pure returns (uint8 team0, uint8 team1) {
        (team0, team1) = getTeamsInCupPlayoffMatch(matchIdxInDay);
        return (sortedTeamIdxInCup[team0], sortedTeamIdxInCup[team1]);
    }



    function getTeamsInCupLeagueMatch(uint8 groupIdx, uint8 matchday, uint8 matchIdxInDay) public pure returns (uint8, uint8) 
    {
        require(matchday < MATCHDAYS/2, "wrong match day");
        (uint8 homeIdx, uint8 visitorIdx) = getTeamsInLeagueMatch(matchday, matchIdxInDay);
        return (getTeamIdxInCup(groupIdx, homeIdx), getTeamIdxInCup(groupIdx, visitorIdx)); 
    }
    
    function getTeamsInLeagueMatch(uint8 matchday, uint8 matchIdxInDay) public pure returns (uint8 homeIdx, uint8 visitorIdx) 
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

    // returns two sorted lists, [worst teamIdxInLeague, points], ....
    function computeLeagueLeaderBoard(uint8[2][MATCHES_PER_LEAGUE] memory results, uint8 matchDay) public pure returns (
        uint8[TEAMS_PER_LEAGUE] memory ranking, uint8[TEAMS_PER_LEAGUE] memory points
    ) {
        require(matchDay < MATCHDAYS, "wrong matchDay");
        uint8 team0;
        uint8 team1;
        for(uint8 m = 0; m < matchDay * 4; m++) {
            (team0, team1) = getTeamsInLeagueMatch(m / 4, m % 4); 
            if (results[m][0] == results[m][1]) {
                points[team0] += 1;
                points[team1] += 1;
            } else if (results[m][0] > results[m][1]) {
                points[team0] += 3;
            } else {
                points[team1] += 3;
            }
        }
        ranking = sortIdxs(points);
    }
}