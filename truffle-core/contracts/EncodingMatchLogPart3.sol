pragma solidity >=0.5.12 <0.6.2;
import "./EncodingMatchLogPart4.sol";
/**
 * @title Library of functions to serialize matchLogs
 */

contract EncodingMatchLogPart3 is EncodingMatchLogPart4 {

    function addAssister(uint256 log, uint8 player, uint8 pos)  public pure returns (uint256) {
        return log | (uint256(player) << (4 + 4 * pos));
    }
  
    function addShooter(uint256 log, uint8 player, uint8 pos)  public pure returns (uint256) {
        return log | (uint256(player) << (52 + 4 * pos));
    }
  
    function addForwardPos(uint256 log, uint8 player, uint8 pos)  public pure returns (uint256) {
        return log | (uint256(player) << (100 + 2 * pos));
    }
    
    function getNGoals(uint256 log) public pure returns (uint8) {
        return uint8(log & 15);
    }

    function addWinnerToBothLogs(uint256[2] memory logs, uint8 winner)  public pure returns (uint256[2] memory) {
        logs[0] |= (uint256(winner) << 209);
        logs[1] |= (uint256(winner) << 209);
        return logs;
    }

    function addNDefs(uint256 log, uint8 nDefs, bool is2ndHalf)  public pure returns (uint256) {
        return log | (uint256(nDefs) << (197 + 4 * (is2ndHalf ? 1 : 0)));
    }

    function setIsHomeStadium(uint256 log, bool val)  public pure returns (uint256) {
        return (log & ~(uint256(1) << 247)) | (uint256(val ? 1: 0) << 247);
    }
}
