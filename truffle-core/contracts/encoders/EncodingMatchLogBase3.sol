pragma solidity >= 0.6.3;
import "./EncodingMatchLogBase4.sol";

/**
 @title Subset of Library of functions to serialize matchLogs
 @author Freeverse.io, www.freeverse.io
 @dev see EncodingMatchLog.sol for full spec
*/

contract EncodingMatchLogBase3 is EncodingMatchLogBase4 {

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
        logs[0] |= (uint256(winner) << 205);
        logs[1] |= (uint256(winner) << 205);
        return logs;
    }

    function addNDefs(uint256 log, uint8 nDefs, bool is2ndHalf)  public pure returns (uint256) {
        return log | (uint256(nDefs) << (193 + 4 * (is2ndHalf ? 1 : 0)));
    }

    function setIsHomeStadium(uint256 log, bool val)  public pure returns (uint256) {
        return (log & ~(uint256(1) << 243)) | (uint256(val ? 1: 0) << 243);
    }
}
