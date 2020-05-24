pragma solidity >= 0.6.3;
import "./EncodingMatchLogBase4.sol";

/**
 * @title Library of functions to serialize matchLogs
 */

contract EncodingMatchLogBase1 is EncodingMatchLogBase4{

    uint256 private constant ONE256       = 1; 

    function addNGoals(uint256 log, uint8 goals) public pure returns (uint256) {
        return log + goals;
    }
    
    function addScoredPenalty(uint256 log, uint8 pos)  public pure returns (uint256) {
        return log | (ONE256 << (128 + pos));
    }

    function getOutOfGamePlayer(uint256 log, bool is2ndHalf)  public pure returns (uint256) {
        uint8 offset = is2ndHalf ? 155 : 135;
        return ((uint256(log) >> offset) & 15);
    }

    function getOutOfGameType(uint256 log, bool is2ndHalf)  public pure returns (uint256) {
        uint8 offset = is2ndHalf ? 163 : 143;
        return ((uint256(log) >> offset) & 3);
    }

    function getOutOfGameRound(uint256 log, bool is2ndHalf)  public pure returns (uint256) {
        uint8 offset = is2ndHalf ? 155 : 135;
        return ((uint256(log) >> offset +4 ) & 15);
    }

    function addOutOfGame(uint256 log, uint8 player, uint8 round, uint8 typeOfOutOfGame, bool is2ndHalf)  public pure returns (uint256) {
        uint8 offset = is2ndHalf ? 155 : 135;
        /// in total, we will write 4b + 4b + 2b = 10b
        log &= ~(uint256(1023) << offset); /// note: 2**10-1 = 1023
        log |= (uint256(player) << offset);
        log |= (uint256(round) << offset+4);
        return log | (uint256(typeOfOutOfGame) << offset+8);
    }

    function addYellowCard(uint256 log, uint8 player, uint8 posInHaf, bool is2ndHalf)  public pure returns (uint256) {
        uint8 offset = (is2ndHalf ? 165 : 145) + posInHaf * 4;
        return log | (uint256(player) << offset);
    }
    
    function getYellowCard(uint256 log, uint8 posInHaf, bool is2ndHalf)  public pure returns (uint8) {
        uint8 offset = (is2ndHalf ? 165 : 145) + posInHaf * 4;
        return uint8((log >> offset) & 15);
    }

    function setYellowedDidNotFinished1stHalf(uint256 log, uint8 posInHaf)  public pure returns (uint256) {
        return log | (ONE256 << (153 + posInHaf));
    }

    function getYellowedDidNotFinished1stHalf(uint256 log, uint8 posInHaf)  public pure returns (bool) {
        return ((log >> (153 + posInHaf)) & 1) == 1;
    }
    
    /// recall that 0 means no subs, and we store here lineUp[p]+1 (where lineUp[p] = player shirt in the 25 that was substituted)
    function addHalfTimeSubs(uint256 log, uint8 player, uint8 pos)  public pure returns (uint256) {
        return log | (uint256(player) << (185 + 5 * pos));
    }

    function setInGameSubsHappened(uint256 log, uint8 happenedType, uint8 posInHalf, bool is2ndHalf) public pure returns (uint256) {
        uint8 offset = 173 + 2 * (posInHalf + (is2ndHalf ? 3 : 0));
        return (log & ~(uint256(3) << offset)) | (uint256(happenedType) << offset);
    }

    function addWinner(uint256 log, uint8 winner)  public pure returns (uint256) {
        return log | (uint256(winner) << 212);
    }

    function addTeamSumSkills(uint256 log, uint256 extraSumSkills)  public pure returns (uint256) {
        return log | (uint256(extraSumSkills) << 214);
    }

}
