pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize matchLogs
 */

contract EncodingMatchLogPart2  {

    uint256 private constant ONE256       = 1; 
    uint256 private constant CHG_HAPPENED        = uint256(1); 
    uint256 private constant CHG_CANCELLED       = uint256(2); 

    function getAssister(uint256 log, uint8 pos) public pure returns (uint8) {
        return uint8((log >> (4 + 4 * pos)) & 15);
    }

  
    function getShooter(uint256 log, uint8 pos) public pure returns (uint8) {
        return uint8((log >> (52 + 4 * pos)) & 15);
    }

    function getForwardPos(uint256 log, uint8 pos) public pure returns (uint8) {
        return uint8((log >> (100 + 2 * pos)) & 3);
    }
    
    function getPenalty(uint256 log, uint8 pos)  public pure returns (bool) {
        return ((log >> (128 + pos)) & 1) == 1;
    }
    
    function addIsHomeStadium(uint256 log)  public pure returns (uint256) {
        return log | (ONE256 << 211);
    }
    
    function getIsHomeStadium(uint256 log)  public pure returns (bool) {
        return ((log >> 211) & 1) == 1;
    }
    
    // recall that 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
    function getHalfTimeSubs(uint256 log, uint8 pos)  public pure returns (uint8) {
        return uint8((log >> (185 + 4 * pos)) & 15);
    }

    function getNDefs(uint256 log, bool is2ndHalf)  public pure returns (uint8) {
        return uint8((log >> (197 + 4 * (is2ndHalf ? 1 : 0))) & 15);
    }

    function addNTot2ndHalf(uint256 log, uint8 nTot)  public pure returns (uint256) {
        return log | (uint256(nTot) << 205);
    }

    function getNTot2ndHalf(uint256 log)  public pure returns (uint8) {
        return uint8((log >> 205) & 15);
    }

    function getWinner(uint256 log) public pure returns (uint8) {
        return uint8((log >> 209) & 3);
    }
    
    function getTeamSumSkills(uint256 log) public pure returns (uint256) {
        return (log >> 211) & 16777215; // 2^24 - 1
    }
    
    function addTrainingPoints(uint256 log, uint256 points)  public pure returns (uint256) {
        return log | (uint256(points) << 235);
    }

    function getTrainingPoints(uint256 log)  public pure returns (uint256) {
        return  (log >> 235) & 4095; // 2^12-1
    }
    
}