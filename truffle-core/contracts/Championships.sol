pragma solidity >=0.4.21 <0.6.0;

import "./Engine.sol";
/**
 * @title Scheduling of leagues, and calls to Engine to resolve games.
 */

contract Championships {
    
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 constant public TEAMS_PER_LEAGUE = 8;
    uint8 constant public MATCHDAYS = 14;
    uint8 constant public MATCHES_PER_DAY = 4;
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

    // returns [scoreHome, scoreAway, scoreHome, scoreAway,...]
    // TODO: currentVerseSeed must be provided from getCurrentVerseSeed()
    // TODO: likewise, matchday should be computed outside
    // function computeMatchday(
    //     uint8 matchday,
    //     uint256[PLAYERS_PER_TEAM_MAX][TEAMS_PER_LEAGUE] memory prevLeagueState,
    //     uint256[TEAMS_PER_LEAGUE] memory tacticsIds,
    //     uint256 currentVerseSeed
    // )
    //     public
    //     view
    //     returns (uint8[2 * MATCHES_PER_DAY] memory scores)
    // {
    //     uint8[2] memory score;
    //     uint8 homeTeamIdx;
    //     uint8 visitorTeamIdx;
    //     for (uint8 matchIdxInDay = 0; matchIdxInDay < MATCHES_PER_DAY ; matchIdxInDay++)
    //     {
    //         (homeTeamIdx, visitorTeamIdx) = getTeamsInLeagueMatch(matchday, matchIdxInDay);
    //         uint256 matchSeed = uint256(keccak256(abi.encode(currentVerseSeed, matchIdxInDay))); 
    //         uint256[2] memory tactics = [tacticsIds[homeTeamIdx], tacticsIds[visitorTeamIdx]];
    //         uint256[PLAYERS_PER_TEAM_MAX][2] memory states = [prevLeagueState[homeTeamIdx], prevLeagueState[visitorTeamIdx]];
    //         uint8[4] memory events0;
    //         (score, events0, events0) = _engine.playMatch(
    //             matchSeed, 
    //             states,
    //             tactics,
    //             false,
    //             false
    //         );
    //         scores[matchIdxInDay * 2] = score[0];
    //         scores[matchIdxInDay * 2 +1 ] = score[1];
    //     }
    // }    
}