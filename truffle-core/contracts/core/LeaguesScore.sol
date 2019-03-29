pragma solidity ^0.5.0;

import "./LeagueUpdatable.sol";

contract LeaguesScore is LeagueUpdatable {
    function scoresCreate() public pure returns (uint16[] memory) {
    }

    function encodeScore(uint8 home, uint8 visitor) public pure returns (uint16 score) {
        score |= home * 2 ** 8;
        score |= visitor;
    }

    function decodeScore(uint16 score) public pure returns (uint8 home, uint8 visitor) {
        home = uint8(score / 2 ** 8);
        visitor = uint8(score & 0x00ff);
    }

    function scoresAppend(uint16[] memory scores, uint16 score) public pure returns (uint16[] memory) {
        uint16[] memory result = new uint16[](scores.length + 1);
        for (uint256 i = 0; i < scores.length ; i++)
            result[i] = scores[i];
        result[result.length-1] = score;
        return result;
    }

    function scoresConcat(uint16[] memory target, uint16[] memory scores) public pure returns (uint16[] memory) {
        uint16[] memory result = new uint16[](target.length + scores.length);
        for (uint256 i = 0 ; i < target.length ; i++) 
            result[i] = target[i];
        for (uint256 i = 0 ; i < scores.length ; i++)
            result[target.length + i] = scores[i];
        return result;
    }

    function scoresGetDay(uint256 id, uint256 day) public view returns (uint16[] memory dayScores) {
        uint256 matchPerDay = getMatchPerDay(id);
        uint16[] memory scores = getScores(id);
        uint256 first = day*matchPerDay;
        dayScores = new uint16[](matchPerDay);
        for (uint256 i = 0 ; i < matchPerDay ; i++) 
            dayScores[i] = scores[first+i];
    }
}