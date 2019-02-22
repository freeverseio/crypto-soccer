pragma solidity ^0.5.0;

contract Scores {
    uint16 constant public DIVIDER = 0xffff;

    function encodeScore(uint8 home, uint8 visitor) public pure returns (uint16 score) {
        require(isValidGoals(home), "invalid home goals");
        require(isValidGoals(visitor), "invalid visitor goals");
        score |= home * 2 ** 8;
        score |= visitor;
    }

    function isValidGoals(uint8 goals) public pure returns (bool) {
        return goals != 0xff;
    }

    function decodeScore(uint16 score) public pure returns (uint8 home, uint8 visitor) {
        require(isValidScore(score), "invalid score");
        home = uint8(score / 2 ** 8);
        visitor = uint8(score & 0x00ff);
    }

    function isValidScore(uint16 score) public pure returns (bool) {
        return score != DIVIDER;
    }

    /// TODO: maybe addScoreToDay ?
    function addToDayScores(uint16[] memory scores, uint16 score) public pure returns (uint16[] memory) {
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

    function scoresConcat(uint256[2][] memory left, uint256[2][] memory right) public pure returns (uint256[2][] memory) {
        if(left.length == 0)
            return right;
        if(right.length == 0)
            return left;

        uint256[2][] memory result = new uint256[2][](left.length + right.length + 1);
        uint256 i;
        for (i = 0; i < left.length ; i++){
            result[i][0] = left[i][0];
            result[i][1] = left[i][1];
        }
        result[left.length][0] = DIVIDER;
        result[left.length][1] = DIVIDER;
        for (i = 0 ; i < right.length ; i++){
            result[left.length + 1 + i][0] = right[i][0];
            result[left.length + 1 + i][1] = right[i][1];
        }

        return result;        
    }

    function isValid(uint256[] memory scores) public pure returns (bool)
    {
        if (scores.length == 0)
            return true;
        if (scores.length % 2 != 0)
            return false;
        if (scores[0] == DIVIDER)
            return false;
        if (scores[scores.length - 1] == DIVIDER)
            return false;
        for (uint256 i = 0 ; i < scores.length - 1 ; i++)
            if (scores[i] == DIVIDER && scores[i+1] == DIVIDER)
                return false;
        return true;
    }

    /// @return number of scores days    
    function scoresCountDays(uint256[] memory scores) public pure returns (uint256) {
        require(isValid(scores), "invalid scores");
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