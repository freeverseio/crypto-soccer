pragma solidity >= 0.6.3;
/**
 * @title Library of functions to serialize matchLogs
 */

contract EncodingMatchLogBase2  {

    uint256 private constant ONE256       = 1; 

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
    
    function getIsHomeStadium(uint256 log)  public pure returns (bool) {
        return ((log >> 250) & 1) == 1;
    }
    
    /// recall that 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
    function getHalfTimeSubs(uint256 log, uint8 pos)  public pure returns (uint8) {
        return uint8((log >> (185 + 5 * pos)) & 31);
    }

    function getNDefs(uint256 log, bool is2ndHalf)  public pure returns (uint8) {
        return uint8((log >> (200 + 4 * (is2ndHalf ? 1 : 0))) & 15);
    }

    function addNTot2ndHalf(uint256 log, uint8 nTot)  public pure returns (uint256) {
        return log | (uint256(nTot) << 208);
    }

    function getNTot2ndHalf(uint256 log)  public pure returns (uint8) {
        return uint8((log >> 208) & 15);
    }

    function getWinner(uint256 log) public pure returns (uint8) {
        return uint8((log >> 212) & 3);
    }
    
    function getTeamSumSkills(uint256 log) public pure returns (uint256) {
        return (log >> 214) & 16777215; /// 2^24 - 1
    }
    
    function addTrainingPoints(uint256 log, uint256 points)  public pure returns (uint256) {
        return log | (uint256(points) << 238);
    }

    function getTrainingPoints(uint256 log)  public pure returns (uint16) {
        return  uint16((log >> 238) & 4095); /// 2^12-1
    }
    
}
