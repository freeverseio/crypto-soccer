pragma solidity >=0.5.12 <0.6.2;

import "./EncodingMatchLogPart1.sol";
import "./EncodingMatchLogPart2.sol";
import "./EncodingMatchLogPart3.sol";
// note: both part1 and part2 inherit from part4

/**
 * @title Library of functions to serialize matchLogs
 */

    // MAX_GOALS = ROUNDS_PER_MATCH = 12,  NO_OUT_OF_GAME_PLAYER = 14
    // 
    // uint8 nGoals, // 4b, offset 0 
    // uint8 assistersIdx[MAX_GOALS], 4b each, offset 4
    // uint8 shootersIdx[MAX_GOALS], 4b each, offset 52
    // uint8 shooterFwdPos[MAX_GOALS], 2b each, offset 100
    // bool[7] memory penaltyScored, // 1b each, offset 128
    // uint8[2] memory outOfGamePlayer 4b each, offset 155:135
    //      an int between [0, 13], 14 = NOTHING HAPPENED       
    // uint8[2] memory outOfGameRounds, 4b each, offset 155:135 +4 
    // uint8[2] memory outOfGameType, 2b each, offset 155:135 +8
    //      null:  0 (cannot return null if outOfGamePlayer != NO_OUT_OF_GAME_PLAYER)
    //      injuryHard:  1
    //      injuryLow:   2
    //      redCard:     3
    // uint8[4] memory yellowCards, 4b each, // first 2 for first half, other for half 2, offset 165:145
    // bool[2] memory yellowCardedDidNotFinish1stHalf, // 1b each, offset 163  
    // uint8[3] memory halfTimeSubstitutions, // 4b each, offset 185
    //      0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
    // uint8[6] memory ingameSubsHappened, // 2b each, offset 173
    //                             //  ...the first 3 for half 1, the other for half 2.
    //                             // ...0: no change required, 1: change happened, 2: change could not happen, cancelled.  
    // winner, 2b, winner: 0 = home, 1 = away, 2 = draw, offset 209 //     
    // nDefsHalf[2], 4b each, offset 197
    // NTot2ndHalf, 4b offset 205
    // teamSumSkills: 24b // offset 211
    // trainingPoints, 12b, offset 235
    // bool isHomeStadium, // 1b each, offset 247

contract EncodingMatchLog is EncodingMatchLogPart1, EncodingMatchLogPart2, EncodingMatchLogPart3 {

}
