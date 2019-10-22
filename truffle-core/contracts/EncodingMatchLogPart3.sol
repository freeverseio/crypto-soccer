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
contract EncodingMatchLogPart3 {

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


}