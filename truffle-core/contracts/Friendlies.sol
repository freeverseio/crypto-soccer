pragma solidity >=0.4.21 <0.6.0;

import "./EncodingSkills.sol";
import "./EncodingIDs.sol";
import "./SortIdxsAnySize.sol";
/**
 * @title Creation of all game assets via creation of timezones, countries and divisions
 * @dev Timezones range from 1 to 24, with timeZone = 0 being null.
 */

contract Friendlies is SortIdxsAnySize {

    struct Friendly {
        bytes32[2] orgMapHash;
        bytes32[2] skillsHash;
        uint8 newestOrgMapIdx;
        uint8 newestSkillsIdx;
        bytes32 scoresRoot;
        uint8 updateCycleIdx;
        bytes32 actionsRoot;
    }        
    
    function getLeagueMatchDays(uint8 nTeams) private pure returns (uint8) { return 2 * (nTeams-1); }
    function getLeagueMatchesPerDay(uint8 nTeams) private pure returns (uint8) { return nTeams/2; }
    function getLeagueMatchesPerLeague(uint8 nTeams) private pure returns (uint8) { return nTeams * (nTeams-1); }
    
    function getTeamsInLeagueMatch(uint8 matchday, uint8 matchIdxInDay, uint8 nTeams) public pure returns (uint8 homeIdx, uint8 visitorIdx) 
    {
        require(matchday < getLeagueMatchDays(nTeams), "wrong match day");
        require(matchIdxInDay < getLeagueMatchesPerDay(nTeams), "wrong match");
        if (matchday < (nTeams - 1))
            (homeIdx, visitorIdx) = _getTeamsInMatchFirstHalf(matchday, matchIdxInDay, nTeams);
        else
            (visitorIdx, homeIdx) = _getTeamsInMatchFirstHalf(matchday - (nTeams - 1), matchIdxInDay, nTeams);
    }

    function _shiftBack(uint8 t, uint8 nTeams) private pure returns (uint8)
    {
        if (t < nTeams)
            return t;
        else
            return t-(nTeams-1);
    }

    function _getTeamsInMatchFirstHalf(uint8 matchday, uint8 matchIdxInDay, uint8 nTeams) private pure returns (uint8, uint8) 
    {
        uint8 team1 = 0;
        if (matchIdxInDay > 0)
            team1 = _shiftBack(nTeams - matchIdxInDay+matchday, nTeams);

        uint8 team2 = _shiftBack(matchIdxInDay+1+matchday, nTeams);
        if ( (matchday % 2) == 0)
            return (team1, team2);
        else
            return (team2, team1);
    }
    
    function getTeamsInCupMatch(uint8 matchIdxDay, uint8 nTeams, uint256 matchDaySeed) public pure returns (uint8, uint8) {
        if (matchDaySeed == 0) { return (2 * matchIdxDay, 2 * matchIdxDay - 1);}
        else {
            uint256[] memory randoms = new uint256[](nTeams);
            uint8[] memory order   = new uint8[](nTeams);
            for (uint8 i = 0; i < nTeams; i++) {
                order[i] = i;
                randoms[i] = uint256(keccak256(abi.encode(matchDaySeed, i)));
            }
            sortIdxs(randoms, order);
            return (order[2 * matchIdxDay], order[2 * matchIdxDay - 1]);
        }
    }
}
