pragma solidity >=0.4.21 <0.6.0;
import "./EncodingMatchLogPart4.sol";

/**
 * @title Library of functions to serialize matchLogs
 */

        // uint8 nGoals, // 4b, offset 0
        // uint8 assistersIdx[12], 4b each, offset 4
        // uint8 shootersIdx[12], 4b each, offset 52
        // uint8 shooterFwdPos[12], 2b each, offset 100
        // bool[7] memory penalties, // 1b each, offset 128
        // uint8[2] memory outOfGames 4b each
        // uint8[4] memory yellowCards, 4b each, // first 2 for first half, other for half 2
        // uint8[2] memory outOfGameRounds,  
        // uint8[2] memory typesOutOfGames, 
        // bool[3] memory yellowCardedDidNotFinish1stHalf, 1b each,
        // bool isHomeStadium, // 1b each,
        // uint8[3] memory halfTimeSubstitutions, // 4b each, offset 201, the first 3 for half 1, the other for half 2
        // uint8[6] memory ingameSubs, // 2b each, offset 189
        //                             //  ...the first 3 for half 1, the other for half 2.
        //                             // ...0: no change required, 1: change happened, 2: change could not happen  
        // uint8[4] memory numDefTotWinner
        //                             // [4b, 4b, 4b, 2b], offset 213
        //                             // [nDefsHalf1, nDefsHalf2, nTotHalf2, winner]
        //                             // winner: 0 = home, 1 = away, 2 = draw
        // teamSumSkills: 24b // offset 227
contract EncodingMatchLogPart1 is EncodingMatchLogPart4{

    uint256 private constant ONE256       = 1; 
    uint256 private constant CHG_HAPPENED        = uint256(1); 
    uint256 private constant CHG_CANCELLED       = uint256(2); 

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

    function addOutOfGame(uint256 log, uint8 player, uint8 round, uint8 typeOfOutOfGame, bool is2ndHalf)  public pure returns (uint256) {
        uint8 offset = is2ndHalf ? 155 : 135;
        log |= (uint256(player) << offset);
        log |= (uint256(round) << (offset + 4));
        return log | (uint256(typeOfOutOfGame) << (offset + 8));
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
    
    // recall that 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
    function addHalfTimeSubs(uint256 log, uint8 player, uint8 pos)  public pure returns (uint256) {
        return log | (uint256(player) << (185 + 4 * pos));
    }

    function setInGameSubsHappened(uint256 log, uint8 happenedType, uint8 posInHalf, bool is2ndHalf) public pure returns (uint256) {
        uint8 offset = 173 + 2 * (posInHalf + (is2ndHalf ? 3 : 0));
        return (log & ~(uint256(3) << offset)) | (uint256(happenedType) << offset);
    }

    function addWinner(uint256 log, uint8 winner)  public pure returns (uint256) {
        return log | (uint256(winner) << 209);
    }

    function addTeamSumSkills(uint256 log, uint256 extraSumSkills)  public pure returns (uint256) {
        return log | (uint256(extraSumSkills) << 211);
    }

}