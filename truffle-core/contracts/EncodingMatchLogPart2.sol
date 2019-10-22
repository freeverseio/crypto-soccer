pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize matchLogs
 */

        // uint8 nGoals, // 4b, offset 0
        // uint8 assistersIdx[14], 4b each, offset 4
        // uint8 shootersIdx[14], 4b each, offset 60
        // uint8 shooterFwdPos[14], 2b each, offset 116
        // bool[7] memory penalties, // 1b each, offset 144
        // uint8[2] memory outOfGames 4b each
        // uint8[6] memory yellowCards1, 4b each,
        // uint8[6] memory yellowCards2, 4b each
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
contract EncodingMatchLogPart2 {

    uint256 private constant ONE256       = 1; 
    uint256 private constant CHG_HAPPENED        = uint256(1); 
    uint256 private constant CHG_CANCELLED       = uint256(2); 

    function addNGoals(uint256 log, uint8 goals) public pure returns (uint256) {
        return log + goals;
    }

    function addAssister(uint256 log, uint8 player, uint8 pos)  public pure returns (uint256) {
        return log | (uint256(player) << (4 + 4 * pos));
    }
  
    function getAssister(uint256 log, uint8 pos) public pure returns (uint8) {
        return uint8((log >> (4 + 4 * pos)) & 15);
    }

    function addShooter(uint256 log, uint8 player, uint8 pos)  public pure returns (uint256) {
        return log | (uint256(player) << (60 + 4 * pos));
    }
  
    function getShooter(uint256 log, uint8 pos) public pure returns (uint8) {
        return uint8((log >> (60 + 4 * pos)) & 15);
    }

    function addForwardPos(uint256 log, uint8 player, uint8 pos)  public pure returns (uint256) {
        return log | (uint256(player) << (116 + 2 * pos));
    }
    
    function getForwardPos(uint256 log, uint8 pos) public pure returns (uint8) {
        return uint8((log >> (116 + 2 * pos)) & 3);
    }
    
    function addPenalty(uint256 log, bool penalty, uint8 pos)  public pure returns (uint256) {
        return log | (uint256(penalty ? 1 : 0) << (144 + pos));
    }

    function getPenalty(uint256 log, uint8 pos)  public pure returns (bool) {
        return ((log >> (144 + pos)) & 1) == 1;
    }
    
    function addOutOfGame(uint256 log, uint8 player, uint8 round, uint8 typeOfOutOfGame, bool is2ndHalf)  public pure returns (uint256) {
        uint8 offset = is2ndHalf ? 171 : 151;
        log |= (uint256(player) << offset);
        log |= (uint256(round) << (offset + 4));
        return log | (uint256(typeOfOutOfGame) << (offset + 8));
    }

    function getOutOfGame(uint256 log, bool is2ndHalf)  public pure 
    returns (uint8 player, uint8 round, uint8 typeOfOutOfGame) {
        uint8 offset = is2ndHalf ? 171 : 151;
        player = uint8((log >> offset) & 15);
        round  = uint8((log >> (offset + 4)) & 15);
        typeOfOutOfGame = uint8((log >> (offset + 8)) & 3);
    }

    function addYellowCard(uint256 log, uint8 player, uint8 posInHaf, bool is2ndHalf)  public pure returns (uint256) {
        uint8 offset = (is2ndHalf ? 181 : 161) + posInHaf * 4;
        return log | (uint256(player) << offset);
    }


    function setYellowedDidNotFinished1stHalf(uint256 log, uint8 posInHaf)  public pure returns (uint256) {
        return log | (ONE256 << (169 + posInHaf));
    }

    
    function setInGameSubsHappened(uint256 log, uint8 happenedType, uint8 posInHalf, bool is2ndHalf) public pure returns (uint256) {
        uint8 offset = 189 + 2 * (posInHalf + (is2ndHalf ? 3 : 0));
        return (log & ~(uint256(3) << offset)) | (uint256(happenedType) << offset);
    }

    
    function addIsHomeStadium(uint256 log)  public pure returns (uint256) {
        return log | (ONE256 << 227);
    }
    
    function getIsHomeStadium(uint256 log)  public pure returns (bool) {
        return ((log >> 227) & 1) == 1;
    }
    
    
    function getHalfTimeSubs(uint256 log, uint8 pos)  public pure returns (uint8) {
        return uint8((log >> (201 + 4 * pos)) & 15);
    }


    function getNDefs(uint256 log, bool is2ndHalf)  public pure returns (uint8) {
        return uint8((log >> (213 + 4 * (is2ndHalf ? 1 : 0))) & 15);
    }

    function addNTot2ndHalf(uint256 log, uint8 nTot)  public pure returns (uint256) {
        return log | (uint256(nTot) << 221);
    }

    function getNTot2ndHalf(uint256 log)  public pure returns (uint8) {
        return uint8((log >> 221) & 15);
    }

    function addWinner(uint256 log, uint8 winner)  public pure returns (uint256) {
        return log | (uint256(winner) << 225);
    }

    function getWinner(uint256 log) public pure returns (uint8) {
        return uint8((log >> 225) & 3);
    }
    
    function getTeamSumSkills(uint256 log) public pure returns (uint256) {
        return (log >> 227) & 16777215; // 2^24 - 1
    }
}