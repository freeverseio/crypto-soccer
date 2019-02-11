pragma solidity ^0.4.25;

import "./Leagues.sol";

contract LeaguesScheduler is Leagues {

    function countLeagueDays(uint256 id) public view returns (uint256) 
    {
        uint256 nTeams = countTeams(id);
        return 2*(nTeams - 1);
    }

    function getMatchPerDay(uint256 id) public view returns (uint256)
    {
        uint256 nTeams = countTeams(id);
        return nTeams / 2;
    }

    function shiftBack(uint256 t, uint256 nTeams) public pure returns (uint256)
    {
        if (t < nTeams)
            return t;
        else
            return t-(nTeams-1);
    }

    function getTeamsInMatchFirstHalf(uint256 matchday, uint256 matchIdx, uint256 nTeams) public pure returns (uint256, uint256) 
    {
        uint256 team1 = 0;
        if (matchIdx > 0)
            team1 = shiftBack(nTeams-matchIdx+matchday, nTeams);

        uint256 team2 = shiftBack(matchIdx+1+matchday, nTeams);
        if ( (matchday % 2) == 0)
            return (team1, team2);
        else
            return (team2, team1);
    }

    function getTeamsInMatch(
        uint256 id,
        uint256 matchday, 
        uint256 matchIdx
    ) 
        public 
        view 
        returns (uint256 team0Idx, uint256 team1Idx) 
    {
        require(matchday < countLeagueDays(id), "wrong match day");
        require(matchIdx < getMatchPerDay(id), "wrong match");
        uint256 nTeams = countTeams(id);
        if (matchday < (nTeams - 1))
            (team0Idx, team1Idx) = getTeamsInMatchFirstHalf(matchday, matchIdx, nTeams);
        else
            (team1Idx, team0Idx) = getTeamsInMatchFirstHalf(matchday - (nTeams - 1), matchIdx, nTeams);
    }
}