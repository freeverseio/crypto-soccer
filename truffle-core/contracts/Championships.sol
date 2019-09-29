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
            groupIdx = 2 * (1 + (teamIdxInCup % 8));
        }
    }



    // function getTeamsInCupMatch(uint8 groupIdx, uint8 matchday, uint8 matchIdxInDay) public pure returns (uint8 homeIdx, uint8 visitorIdx) 
    // {
    //     require(matchday < MATCHDAYS, "wrong match day");
    //     require(matchIdxInDay < MATCHES_PER_DAY, "wrong match");
        
    //     if (matchday < (TEAMS_PER_LEAGUE - 1))
    //         (homeIdx, visitorIdx) = _getTeamsInMatchFirstHalf(matchday, matchIdxInDay);
    //     else
    //         (visitorIdx, homeIdx) = _getTeamsInMatchFirstHalf(matchday - (TEAMS_PER_LEAGUE - 1), matchIdxInDay);
    // }



    
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