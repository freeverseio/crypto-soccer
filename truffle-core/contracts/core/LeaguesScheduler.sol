pragma solidity ^0.5.0;

import "./LeaguesBase.sol";

contract LeaguesScheduler is LeaguesBase {
    function countLeagueDays(uint256 id) public view returns (uint256) 
    {
        uint256 nTeams = countTeams(id);
        return 2*(nTeams - 1);
    }

    function getEndBlock(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        uint256 nMatchDays = countLeagueDays(id);
        uint256 initBlock = getInitBlock(id);
        uint256 step = getStep(id);
        return initBlock + (nMatchDays - 1) * step;
    }

    function getMatchPerDay(uint256 id) public view returns (uint256)
    {
        uint256 nTeams = countTeams(id);
        return nTeams / 2;
    }

    function getMatchDayBlockHash(uint256 id, uint256 day) public view returns (bytes32)
    {
        uint256 initBlock = getInitBlock(id);
        uint256 step = getStep(id);
        bytes32 blockHash = blockhash(initBlock + step * day);
        require(blockHash != 0, "unable to retrive block hash");
        return blockHash;
    }

    function getTeamsInMatch(
        uint256 id,
        uint256 matchday, 
        uint256 matchIdx
    ) 
        public 
        view 
        returns (uint256 homeIdx, uint256 visitorIdx) 
    {
        require(matchday < countLeagueDays(id), "wrong match day");
        require(matchIdx < getMatchPerDay(id), "wrong match");
        uint256 nTeams = countTeams(id);
        if (matchday < (nTeams - 1))
            (homeIdx, visitorIdx) = _getTeamsInMatchFirstHalf(matchday, matchIdx, nTeams);
        else
            (visitorIdx, homeIdx) = _getTeamsInMatchFirstHalf(matchday - (nTeams - 1), matchIdx, nTeams);
    }

    function _shiftBack(uint256 t, uint256 nTeams) private pure returns (uint256)
    {
        if (t < nTeams)
            return t;
        else
            return t-(nTeams-1);
    }

    function _getTeamsInMatchFirstHalf(uint256 matchday, uint256 matchIdx, uint256 nTeams) private pure returns (uint256, uint256) 
    {
        uint256 team1 = 0;
        if (matchIdx > 0)
            team1 = _shiftBack(nTeams-matchIdx+matchday, nTeams);

        uint256 team2 = _shiftBack(matchIdx+1+matchday, nTeams);
        if ( (matchday % 2) == 0)
            return (team1, team2);
        else
            return (team2, team1);
    }
}