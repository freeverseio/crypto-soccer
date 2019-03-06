pragma solidity ^0.5.0;

import "./LeaguesScheduler.sol";

contract LeaguesScore is LeaguesScheduler {
    uint16 constant public DIVIDER = 0xffff;

    function encodeScore(uint8 home, uint8 visitor) public pure returns (uint16 score) {
        require(isValidScore(home, visitor), "invalid score");
        score |= home * 2 ** 8;
        score |= visitor;
    }

    function decodeScore(uint16 score) public pure returns (uint8 home, uint8 visitor) {
        require(isValidScore(score), "invalid score");
        home = uint8(score / 2 ** 8);
        visitor = uint8(score & 0x00ff);
    }

    function scoresCreate() public pure returns (uint16[] memory) {

    }

    function isValidScore(uint16 score) public pure returns (bool) {
        return score != DIVIDER;
    }

    function isValidScore(uint8 home, uint8 visitor) public pure returns (bool) {
        return !(home == 0xff && visitor == 0xff);
    }

    /// TODO: maybe addScoreToDay ?
    function scoresAppend(uint16[] memory scores, uint16 score) public pure returns (uint16[] memory) {
        require(score != DIVIDER, "invalid score");
        require(isValidDayScores(scores), "invalid day scores");
        uint16[] memory result = new uint16[](scores.length + 1);
        for (uint256 i = 0; i < scores.length ; i++)
            result[i] = scores[i];
        result[result.length-1] = score;
        return result;
    }

    function isValidDayScores(uint16[] memory dayScores) public pure returns (bool) {
        for (uint256 i = 0 ; i < dayScores.length ; i++) 
            if (dayScores[i] == DIVIDER)
                return false;
        return true;
    }

    function getDayScores(uint16[] memory leagueScores, uint256 day) public pure returns (uint16[] memory dayScores) {
        require(day < countDaysInTournamentScores(leagueScores), "out of range");
        uint256 current;
        uint256 i;
        for(i = 0 ; current < day; i++)
            if (leagueScores[i] == DIVIDER)
                current++;
        for(; i < leagueScores.length && leagueScores[i] != DIVIDER ; i++)
            dayScores = scoresAppend(dayScores, leagueScores[i]);
    }

    function addToTournamentScores(uint16[] memory tournamentScores, uint16[] memory dayScores) public pure returns (uint16[] memory) {
        require(isValidDayScores(dayScores), "invalid day scores");
        if (tournamentScores.length == 0)
            return dayScores;

        uint16[] memory result = new uint16[](tournamentScores.length + 1 + dayScores.length);
        for (uint256 i = 0 ; i < tournamentScores.length ; i++)
            result[i] = tournamentScores[i];
        result[tournamentScores.length] = DIVIDER;
        for (uint256 i = 0 ; i < dayScores.length ; i++)
            result[tournamentScores.length + 1 + i] = dayScores[i];
        return result;
    }

    /// @return number of scores days    
    function countDaysInTournamentScores(uint16[] memory scores) public pure returns (uint256) {
        if (scores.length == 0)
            return 0;

        uint256 count = 1;
        for (uint256 i = 0 ; i < scores.length ; i++) {
            if (scores[i] == DIVIDER)
                count++; 
        }
        return count;
    }
}